package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

// hash functions
// we wanna hash Nodes not the node
// func keyToPath(key []byte)

func key2path(key []byte) string {
	//

	//create request
}

type sortvol struct {
	score  []byte
	volume string
}
type byScore []sortvol

func (s byScore) Len() int      { return len(s) }
func (s byScore) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byScore) Less(i, j int) bool {
	return bytes.Compare(s[i].score, s[j].score) == 1
}

func key2volume(key []byte, volumes []string, count int, svcount int) []string {
	// this is an intelligent way to pick the volume server for a file
	// stable in the volume server name (not position!)
	// and if more are added the correct portion will move (yay md5!)
	var sortvols []sortvol
	for _, v := range volumes {
		hash := md5.New()
		hash.Write(key)
		hash.Write([]byte(v))
		score := hash.Sum(nil)
		sortvols = append(sortvols, sortvol{score, v})
	}
	sort.Stable(byScore(sortvols))
	// go should have a map function
	// this adds the subvolumes
	var ret []string
	for i := 0; i < count; i++ {
		sv := sortvols[i]
		var volume string
		if svcount == 1 {
			// if it's one, don't use the path structure for it
			volume = sv.volume
		} else {
			// use the least significant compare dword for the subvolume
			// using only a byte would cause potential imbalance
			svhash := uint(sv.score[12])<<24 + uint(sv.score[13])<<16 +
				uint(sv.score[14])<<8 + uint(sv.score[15])
			volume = fmt.Sprintf("%s/sv%02X", sv.volume, svhash%uint(svcount))
		}
		ret = append(ret, volume)
	}
	//fmt.Println(string(key), ret[0])
	return ret
}

func needs_rebalance(volumes []string, nvolumes []string) bool {
	if len(volumes) != len(nvolumes) {
		return true
	}
	for i := 0; i < len(volumes); i++ {
		if volumes[i] != nvolumes[i] {
			return true
		}
	}
	return false
}

func remote_delete(remote string) error {
	req, err := http.NewRequest("DELETE", remote, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 && resp.StatusCode != 404 {
		return fmt.Errorf("remote_delete: wrong status code %d", resp.StatusCode)
	}
	return nil
}

func remote_put(remote string, length int64, body io.Reader) error {
	req, err := http.NewRequest("PUT", remote, body)
	if err != nil {
		return err
	}
	req.ContentLength = length
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 && resp.StatusCode != 204 {
		return fmt.Errorf("remote_put: wrong status code %d", resp.StatusCode)
	}
	return nil
}

func remote_get(remote string) (string, error) {
	resp, err := http.Get(remote)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("remote_get: wrong status code %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func remote_head(remote string, timeout time.Duration) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "HEAD", remote, nil)
	if err != nil {
		return false, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200, nil
}
