package main

import wpsgrabber "github.com/crossi-T2/wpsgrabber/wpsgrabber"

func main() {

	configfile := "/Users/crossi/go/src/github.com/crossi-T2/wpsgrabber/wpsgrabber/config.json"

	wpsgrabber.Watch(configfile)
}
