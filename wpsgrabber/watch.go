package wpsgrabber

import (
	"github.com/fsnotify/fsnotify"
	"github.com/tkanos/gonfig"
	"log"
	"os"
	"path"
	"regexp"
)

type Configuration struct {
	RootDir          string
	CSVOutputDir     string
	ProcessorCode    string
	ProcessorVersion string
}

var configuration Configuration = Configuration{}

func Watch(configfile string) {

	err := gonfig.GetConf(configfile, &configuration)

	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// We are only interested in newly created files
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("Created:", event.Name)
					file, err := os.Stat(event.Name)
					if err != nil {
						log.Println(err)
						return
					}

					if file.Mode().IsDir() {
						err = watcher.Add(event.Name)
						if err != nil {
							log.Fatal(err)
						}
						log.Println("Watching:", event.Name)
					} else {
						// We do expect updates in XML files in the form 0.xml 1.xml 2.xml etc.
						matched, _ := regexp.MatchString(`^[0-9]+.xml$`, path.Base(event.Name))
						if matched {
							response := parseExecuteResponse(event.Name)

							if response.Status.ProcessStatus == 0 ||
								response.Status.ProcessStatus == 1 {

								createCSV(response)

								// TODO: After this step we should release the file descriptor
								// of the parent directory, by removing it from the watcher

							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(configuration.RootDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching:", configuration.RootDir)
	<-done

}
