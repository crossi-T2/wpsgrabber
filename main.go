package main

import (
	"flag"
	wpsgrabber "github.com/crossi-T2/wpsgrabber/wpsgrabber"
	"log"
	"time"
)

func main() {

	loc, _ := time.LoadLocation("UTC")

	configFile := flag.String("config", "/etc/wpsgrabber/config.json", "Configuration file path")
	watchFrom := flag.String("watch-from", time.Now().In(loc).Format(time.RFC3339), "Start time ")

	flag.Parse()

	log.Println(*watchFrom)

	err := wpsgrabber.New(*configFile)

	if err != nil {
		log.Fatal(err)
	}
}
