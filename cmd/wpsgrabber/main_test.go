package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(t *testing.T) {

	currentWorkingDir, _ := os.Getwd()

	configFile := filepath.Join(currentWorkingDir, "testdata", "config.yaml")

	yamlFile, _ := ioutil.ReadFile(configFile)
	yaml.Unmarshal(yamlFile, &configuration)

	timeout := time.After(3 * time.Second)
	done := make(chan bool)
	go func() {
		err := New(configFile)
		if err != nil {
			t.Errorf("err: %v", err)
		}
		done <- true
	}()

	// Testing fsnotify
	time.Sleep(1 * time.Second)

	sourcePath := filepath.Join(configuration.RootDir, "example_succeeded.xml")
	destDir := filepath.Join(configuration.RootDir, "new")
	destPath := filepath.Join(destDir, "0.xml")
	os.Mkdir(destDir, 0755)

	time.Sleep(1 * time.Second)
	sourceFile, _ := ioutil.ReadFile(sourcePath)
	ioutil.WriteFile(destPath, sourceFile, 0644)

	<-timeout

	// Clean-up
	OutputDir, _ := os.Open(configuration.OutputDir)
	defer OutputDir.Close()

	OutputDirFiles, _ := OutputDir.Readdirnames(-1)

	for _, OutputDirFile := range OutputDirFiles {
		os.RemoveAll(filepath.Join(configuration.OutputDir, OutputDirFile))
	}

	os.RemoveAll(destDir)
}
