package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/locmai/assistant/actions"
	"log"
	"os"
)

var (
	addr string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&addr, "http", os.Getenv("ADDR"), "HTTP listen address")
	flag.Parse()

	fs := NewServer()
	fs.Addr = addr
	fs.DisableBasicAuth = true

	fs.Actions.Set("create_cluster", actions.CreateClusterHandler)

	if err := fs.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
