package main

import "net/http"

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("hello, worldy"))
	return
}
