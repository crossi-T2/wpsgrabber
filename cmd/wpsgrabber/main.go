package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

type Configuration struct {
	RootDir           string    `yaml:"RootDir"`
	LogFile           string    `yaml:"LogFile"`
	ScanFrom          time.Time `yaml:"ScanFrom"`
	OutputDir         string    `yaml:"OutputDir"`
	ProcessIdentifier string    `yaml:"ProcessIdentifier"`
	ProcessVersion    string    `yaml:"ProcessVersion"`
}

var configuration Configuration = Configuration{}

func main() {

	configFile := flag.String("config", "/etc/wpsgrabber/config.yaml", "Configuration file path")
	flag.Parse()

	err := New(*configFile)

	if err != nil {
		log.Fatal(err)
	}
}

func New(configFile string) error {

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("can't read config file %s: %v ", configFile, err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &configuration)

	if configuration.LogFile != "" {
		logFile, err := os.OpenFile(configuration.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

		if err != nil {
			err = fmt.Errorf("can't open log file %s: %v ", configuration.LogFile, err)
			return err
		}

		log.SetOutput(logFile)
	}

	// If ScanFrom was configured, it would scan the RootDir for WPS Execute response reports
	if !configuration.ScanFrom.IsZero() {

		// Walks the RootDir for reports
		err := filepath.Walk(configuration.RootDir,
			func(path string, file os.FileInfo, err error) error {
				if err != nil {
					err = fmt.Errorf("can't walk %s: %v ", path, err)
					return err
				}
				if !file.IsDir() {
					if file.ModTime().After(configuration.ScanFrom) {

						// We do expect updates in XML files in the form 0.xml 1.xml 2.xml etc.
						matched, _ := regexp.MatchString(`^[0-9]+.xml$`, filepath.Base(path))
						if matched {
							response, err := parseExecuteResponse(path)

							if err != nil {
								err = fmt.Errorf("can't parse %s: %v ", path, err)
							}

							if response.Status.ProcessStatus == 0 ||
								response.Status.ProcessStatus == 1 {

								log.Println("found:", path)

								err = EncodeResponse(response, path)

								if err != nil {
									err = fmt.Errorf("can't encode %s: %v ", path, err)
									return err
								}

							}
						}
					}
				}
				return nil
			})
		if err != nil {
			err = fmt.Errorf("can't walk %s: %v ", configuration.RootDir, err)
			return err
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		err = fmt.Errorf("can't create fsnotify watcher: %v", err)
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() error {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return errors.New("fsnotify error")
				}
				// We are only interested in newly created files
				if event.Op&fsnotify.Create == fsnotify.Create {
					file, err := os.Stat(event.Name)
					if err != nil {
						err = fmt.Errorf("can't get information from %s: %v ", event.Name, err)
						return err
					}

					if file.Mode().IsDir() {
						log.Println("new directory:", event.Name)
						err = watcher.Add(event.Name)
						if err != nil {
							err = fmt.Errorf("can't watch %s: %v ", event.Name, err)
							return err
						}
						log.Println("watching:", event.Name)
					} else {
						log.Println("new file:", event.Name)
						// We do expect updates in XML files in the form 0.xml 1.xml 2.xml etc.
						matched, _ := regexp.MatchString(`^[0-9]+.xml$`, path.Base(event.Name))
						if matched {
							response, err := parseExecuteResponse(event.Name)
							if err != nil {
								err = fmt.Errorf("can't parse %s: %v ", event.Name, err)
								return err
							}

							if response.Status.ProcessStatus == 0 ||
								response.Status.ProcessStatus == 1 {

								err = EncodeResponse(response, event.Name)

								if err != nil {
									err = fmt.Errorf("can't encode %s: %v ", event.Name, err)
									return err
								}

								// At this stage, there is no need to continue watching the parent
								// folder, since the processing execution information has been managed.
								// Removing then parent directory from the watcher list to release
								// the related file descriptor
								parentDir := filepath.Dir(event.Name)

								if parentDir != configuration.RootDir {
									watcher.Remove(parentDir)
								}
							}
						} else {
							log.Printf("filename does not match regex '^[0-9]+.xml$', skipping %s", event.Name)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return err
				}
			}
		}
	}()

	err = watcher.Add(configuration.RootDir)
	if err != nil {
		err = fmt.Errorf("can't watch %s: %v", configuration.RootDir, err)
		return err
	}

	log.Println("watching:", configuration.RootDir)
	<-done

	return nil

}

func EncodeResponse(response *ExecuteResponse, sourcePath string) error {

	CSVFilename := filepath.Join(configuration.OutputDir, response.Process.WorkflowIdentifier+"_run.csv")
	XMLFilename := filepath.Join(configuration.OutputDir, response.Process.WorkflowIdentifier+"_request.xml")

	err := createCSV(CSVFilename, response)

	if err != nil {
		err = fmt.Errorf("failed CSV encoding for %s: %v", sourcePath, err)
		return err
	}

	XMLRequestPath := filepath.Join(filepath.Dir(sourcePath), "request.xml")
	XMLRequestFile, err := ioutil.ReadFile(XMLRequestPath)

	if err != nil {
		err = fmt.Errorf("failed reading %s: %v", XMLRequestPath, err)
		return err
	}

	ioutil.WriteFile(XMLFilename, XMLRequestFile, 0644)

	log.Println("CSV file created:", CSVFilename)
	log.Println("XML file created:", XMLFilename)

	return nil
}
