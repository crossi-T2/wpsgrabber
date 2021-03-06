package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type ExecuteResponse struct {
	XMLName xml.Name `xml:"ExecuteResponse"`
	Process Process  `xml:"Process"`
	Status  Status   `xml:"Status"`
}

type Process struct {
	XMLName            xml.Name `xml:"Process"`
	Identifier         string   `xml:"Identifier"`
	Version            string   `xml:"processVersion,attr"`
	Title              string   `xml:"Title"`
	WorkflowIdentifier string
}

type Status struct {
	XMLName          xml.Name  `xml:"Status"`
	CreationTime     time.Time `xml:"creationTime,attr"`
	EndTime          time.Time
	ProcessFailed    xml.Name `xml:"ProcessFailed,omitempty"`
	ProcessSucceeded xml.Name `xml:"ProcessSucceeded,omitempty"`
	ProcessStatus    int
}

func parseExecuteResponse(responseFile string) (*ExecuteResponse, error) {

	file, err := os.Open(responseFile)
	if err != nil {
		err = fmt.Errorf("can't open file: %v ", err)
		return nil, err
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var executeResponse ExecuteResponse

	// Default value for ProcessStatus, otherwise it will set
	// the zero-value for it, which means Succeeded in this context
	executeResponse.Status.ProcessStatus = 999

	if err := xml.Unmarshal(byteValue, &executeResponse); err != nil {
		err = fmt.Errorf("can't unmarshal file: %v ", err)
		return nil, err
	}

	if executeResponse.Status.ProcessSucceeded.Local != "" ||
		executeResponse.Status.ProcessFailed.Local != "" {

		// Determine the EndTime of the executeResponse by inspecting
		// the modification time of the responseFile
		stat, _ := os.Stat(responseFile)
		executeResponse.Status.EndTime = stat.ModTime()

		// Determine the WorkflowIdentifier, based on the name of the parent
		// directory
		executeResponse.Process.WorkflowIdentifier = filepath.Base(filepath.Dir(responseFile))

		// Determine the value for executeResponse.Status.ProcessStatus
		if executeResponse.Status.ProcessSucceeded.Local != "" {
			executeResponse.Status.ProcessStatus = 0
		} else if executeResponse.Status.ProcessFailed.Local != "" {
			executeResponse.Status.ProcessStatus = 1
		}
	}

	return &executeResponse, nil
}
