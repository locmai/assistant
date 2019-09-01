package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/dialogflow/fulfillment"
	"github.com/locmai/assistant/actions"
	"log"
	"os"
)

var (
	addr           string
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.StringVar(&addr, "http", os.Getenv("ADDR"), "HTTP listen address")
	flag.Parse()


	fs := fulfillment.NewServer()
	fs.Addr = addr
	fs.DisableBasicAuth = true
	fs.Actions.Set("create_cluster", actions.CreateCluster)


	if err := fs.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
