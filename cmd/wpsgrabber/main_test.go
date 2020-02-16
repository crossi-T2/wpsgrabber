package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// mutex to avoid a race condition during tests
type TestConfiguration struct {
	value Configuration
	m     sync.Mutex
}

func (conf *TestConfiguration) Get() Configuration {
	conf.m.Lock()
	defer conf.m.Unlock()
	return conf.value
}

func (conf *TestConfiguration) Set(configFile string) {
	conf.m.Lock()
	defer conf.m.Unlock()

	yamlFile, _ := ioutil.ReadFile(configFile)
	yaml.Unmarshal(yamlFile, &conf.value)
}

func TestMain(t *testing.T) {

	timeoutProducer := time.After(1 * time.Second)
	timeoutGeneral := time.After(3 * time.Second)

	currentWorkingDir, _ := os.Getwd()
	configFile := filepath.Join(currentWorkingDir, "testdata", "config.yaml")

	configuration := &TestConfiguration{}
	configuration.Set(configFile)

	go func() {
		err := New(configFile)
		if err != nil {
			t.Errorf("err: %v", err)
		}
	}()

	go func() {

		<-timeoutProducer

		// Testing fsnotify
		NewRequestsDir := filepath.Join(configuration.Get().RootDir, "new")
		os.Mkdir(NewRequestsDir, 0755)

		time.Sleep(100 * time.Millisecond)
		sourcePath := filepath.Join(configuration.Get().RootDir, "example_succeeded.xml")
		destPath := filepath.Join(NewRequestsDir, "0.xml")

		sourceFile, _ := ioutil.ReadFile(sourcePath)
		ioutil.WriteFile(destPath, sourceFile, 0644)
	}()

	<-timeoutGeneral

	// Clean-up
	OutputDir, _ := os.Open(configuration.Get().OutputDir)
	defer OutputDir.Close()

	OutputDirFiles, _ := OutputDir.Readdirnames(-1)

	for _, OutputDirFile := range OutputDirFiles {
		os.RemoveAll(filepath.Join(configuration.Get().OutputDir, OutputDirFile))
	}

	os.RemoveAll(filepath.Join(configuration.Get().RootDir, "new"))
}
