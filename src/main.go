package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	rand.Seed(time.Now().Unix())
	port := flag.Int("port", 8000, "Port for the master server. (default 8000)")
	dbPath := flag.String("db", "./db", "Path to leveldb (default ./db)")
	fallback := flag.String("fallback", "", "Fallback server for missing keys")
	replicas := flag.Int("replicas", 4, "Number of replicas to make of the data")
	volumesStr := flag.String("volumes", "localhost:8001", "Volumes to use (for storage servers), comma separated (default only one - localhost:8001)")
	subVolumes := flag.Int("subvolumes", 10, "Number of subvolumes.")
	protect := flag.Bool("protect", false, "UNLINK before DELETE")
	md5sum := flag.Bool("md5sum", true, "Calculate and store MD5 checksum of values")
	voltimeout := flag.Duration("voltimeout", 1*time.Second, "Volume servers must respond to GET/HEAD requests in this amount of time or they are considered down, as duration")

	flag.Parse()
	volumes := strings.Split(*volumesStr, ",")

	cmd := flag.Arg(0)

	if cmd != "serve" && cmd != "rebuild" && cmd != "rebalance" {
		fmt.Println("Usage: ./steres <serve, rebuild, rebalance>")
		flag.PrintDefaults()
		return
	}
	// flag checks
	if *dbPath == "" {
		panic("Invalid path for the database")
	}

	if len(volumes) < *replicas {
		panic("Number of volumes must be at least as many as replicas")
	}

	db, err := leveldb.OpenFile(*dbPath, nil)
	if err != nil {
		panic(fmt.Sprintf("Error: Failed to open LevelDB %s", err))
	}
	defer db.Close()
	app := App{
		db:            db,
		lock:          make(map[string]struct{}),
		volumes:       volumes,
		subvolumes:    *subVolumes,
		fallback:      *fallback,
		replicas:      *replicas,
		protect:       *protect,
		md5sum:        *md5sum,
		volumeTimeout: *voltimeout,
	}
	if cmd == "serve" {
		http.ListenAndServe(fmt.Sprintf(":%d", *port), &app)
	} else {
		fmt.Println("Other options not working at the moment.")
	}

}
