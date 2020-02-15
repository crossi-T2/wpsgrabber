package main

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateCSV(t *testing.T) {

	var status Status
	var process Process
	var response ExecuteResponse

	status.ProcessStatus = 0
	status.CreationTime, _ = time.Parse(time.RFC3339, "2020-04-26T11:16:58.912+02:00")
	status.EndTime, _ = time.Parse(time.RFC3339, "2020-04-30T11:16:58.912+02:00")
	status.ProcessStatus = 0

	process.Identifier = "com.terradue.wps_oozie.process.OozieAbstractAlgorithm"
	process.Version = "1.0.0"
	process.Title = "SRTM Digital Elevation Model"
	process.WorkflowIdentifier = ""

	response.Process = process
	response.Status = status

	CSVfilename := filepath.Join(configuration.OutputDir, "_run.csv")

	err := createCSV(CSVfilename, &response)

	if err != nil {
		t.Errorf("createCSV failed")
	}

	CSVfile, err := os.Open(CSVfilename)

	if err != nil {
		t.Errorf("can't open file: %s", CSVfilename)
	}

	defer CSVfile.Close()

	var expectedCSV = "com.terradue.wps_oozie.process.OozieAbstractAlgorithm,1.0.0,,_request.xml,2020-04-26T09:16:58Z,2020-04-30T09:16:58Z,0"

	scanner := bufio.NewScanner(CSVfile)
	scanner.Scan()
	actualCSV := scanner.Text()

	if actualCSV != expectedCSV {
		t.Errorf("createCSV returned %s instead of %s", actualCSV, expectedCSV)
	}

	os.Remove(CSVfilename)
}
