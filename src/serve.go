package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

type Sequence struct {
	Sequence string
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var sequence Sequence
	err := dec.Decode(&sequence)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := []byte(sequence.Sequence)

	if len(NSplit(key, a.tokenSize)) < 4 {
		w.Write([]byte("Sequence too short"))
		w.WriteHeader(400)
		return
	}

	// seqBytes := []bytes(sequence)

	// lock the sequence while a PUT or DELETE is in progress
	if r.Method == "PUT" || r.Method == "POST" || r.Method == "DELETE" || r.Method == "UNLINK" {
		if !a.LockSequence(key) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		defer a.UnlockSequence(key)
	}

	switch r.Method {
	case "POST":
		rec, id := a.UpdateSequence(key)
		var remote string
		if rec.deleted == SOFT || rec.deleted == HARD {
			if a.fallback == "" {
				w.Header().Set("Content-Length", "0")
				w.WriteHeader(404)
				return
			}
			// fall through to fallback
			remote = fmt.Sprintf("http://%s%s", a.fallback, key)
		} else {
			kvolumes := key2volume(key, a.volumes, a.replicas, a.subvolumes)
			if needs_rebalance(rec.rvolumes, kvolumes) {
				w.Header().Set("Key-Balance", "unbalanced")
				fmt.Println("on wrong volumes, needs rebalance")
			} else {
				w.Header().Set("Key-Balance", "balanced")
			}
			w.Header().Set("Key-Volumes", strings.Join(rec.rvolumes, ","))

		}
		good := false
		for _, vn := range rand.Perm(len(rec.rvolumes)) {
			remote = fmt.Sprintf("http://%s%s", rec.rvolumes[vn], key2path(key))
			found, _ := remote_head(remote, a.voltimeout)
			if found {
				good = true
				break
			}
		}
		// if not found on any volume servers, fail before the redirect
		if !good {
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(404)
			return
		}
		// note: this can race and fail, but in that case the client will handle the retry
		w.Header().Set("Location", remote)
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(302)
	case "PUT":
		// no empty values
		if r.ContentLength == 0 {
			w.WriteHeader(411)
			return
		}

		// check if we already have the key, and it's not deleted
		rec := a.GetSequence(key)
		if rec.deleted == NO {
			// Forbidden to overwrite with PUT
			w.WriteHeader(403)
			return
		}

		// we don't have the key, compute the remote URL
		kvolumes := key2volume(key, a.volumes, a.replicas, a.subvolumes)

		// push to leveldb initially as deleted, and without a hash since we don't have it yet
		if !a.PutSequence(key, Node{kvolumes, SOFT, ""}) {
			w.WriteHeader(500)
			return
		}

		// write to each replica
		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		bodylen := r.ContentLength
		for i := 0; i < len(kvolumes); i++ {
			if i != 0 {
				// if we have already read the contents into the TeeReader
				body = bytes.NewReader(buf.Bytes())
			}
			remote := fmt.Sprintf("http://%s%s", kvolumes[i], key2path(key))
			if remote_put(remote, bodylen, body) != nil {
				// we assume the remote wrote nothing if it failed
				fmt.Printf("replica %d write failed: %s\n", i, remote)
				w.WriteHeader(500)
				return
			}
		}

		var hash = ""
		if a.md5sum {
			// compute the hash of the value
			hash = fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
		}

		// push to leveldb as existing
		// note that the key is locked, so nobody wrote to the leveldb
		if !a.Sequence(key, Sequence{kvolumes, NO, hash}) {
			w.WriteHeader(500)
			return
		}

		// 201, all good
		w.WriteHeader(201)
	}