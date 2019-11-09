package main

import (
	"flag"
	wpsgrabber "github.com/crossi-T2/wpsgrabber/wpsgrabber"
	"log"
)

func main() {

	configFile := flag.String("config", "/etc/wpsgrabber/config.json", "Configuration file path")
	flag.Parse()

	err := wpsgrabber.New(*configFile)

	if err != nil {
		log.Fatal(err)
	}
}
