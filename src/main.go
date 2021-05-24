package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type App struct {
	db *leveldb.DB
}

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("hello, worldy"))
	return
}

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	rand.Seed(time.Now().Unix())
	port := flag.Int("port", 8000, "Port for the master server.")
	dbPath := flag.String("db", ".", "Path to leveldb")
	flag.Parse()

	cmd := flag.Arg(0)
	if cmd != "serve" && cmd != "rebuild" && cmd != "rebalance" {
		fmt.Println("Usage: ./steres <serve, rebuild, rebalance>")
		flag.PrintDefaults()
		return
	}
	db, err := leveldb.OpenFile(*dbPath, nil)
	if err != nil {
		panic(fmt.Sprintf("LevelDB open failed: %s", err))
	}
	defer db.Close()
	app := App{
		db: db,
	}
	if cmd == "serve" {
		http.ListenAndServe(fmt.Sprintf(":%d", *port), &app)
	} else {
		fmt.Println("Other options not working at the moment.")
	}

}
