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

}
