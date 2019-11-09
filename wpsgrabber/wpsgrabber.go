package wpsgrabber

import (
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
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
	ScanFrom          time.Time `yaml:"ScanFrom"`
	CSVOutputDir      string    `yaml:"CSVOutputDir"`
	ProcessIdentifier string    `yaml:"ProcessIdentifier"`
	ProcessVersion    string    `yaml:"ProcessVersion"`
}

var configuration Configuration = Configuration{}

func New(configFile string) error {

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &configuration)

	// If ScanFrom was configured, it would scan the RootDir for WPS Execute response reports
	if !configuration.ScanFrom.IsZero() {

		// Walks the RootDir for reports
		err := filepath.Walk(configuration.RootDir,
			func(path string, file os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !file.IsDir() {
					if file.ModTime().After(configuration.ScanFrom) {

						response := parseExecuteResponse(path)

						if response.Status.ProcessStatus == 0 ||
							response.Status.ProcessStatus == 1 {

							log.Println("Found:", path)
							err = createCSV(response)

							if err != nil {
								errors.Wrap(err, "error creating CSV from response")
								return err
							}
						}
					}
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}

	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errors.Wrap(err, "error creating a fsnotify watcher")
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
						errors.Wrap(err, "error getting information from "+event.Name)
						return err
					}

					if file.Mode().IsDir() {
						log.Println("New directory:", event.Name)
						err = watcher.Add(event.Name)
						if err != nil {
							log.Print(err)
							errors.Wrap(err, "error getting information from "+event.Name)
							return err
						}
						log.Println("Watching:", event.Name)
					} else {
						log.Println("New file:", event.Name)
						// We do expect updates in XML files in the form 0.xml 1.xml 2.xml etc.
						matched, _ := regexp.MatchString(`^[0-9]+.xml$`, path.Base(event.Name))
						if matched {
							response := parseExecuteResponse(event.Name)

							if response.Status.ProcessStatus == 0 ||
								response.Status.ProcessStatus == 1 {

								err = createCSV(response)

								if err != nil {
									errors.Wrap(err, "error creating CSV from response")
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
							log.Println("Filename does not match regex ^[0-9]+.xml$")
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
		errors.Wrap(err, "error watching root dir "+configuration.RootDir)
		return err
	}

	log.Println("Watching:", configuration.RootDir)
	<-done

	return nil

}
