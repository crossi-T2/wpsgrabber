package wpsgrabber

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
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

func parseExecuteResponse(responseFile string) *ExecuteResponse {

	file, err := os.Open(responseFile)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var executeResponse ExecuteResponse

	if err := xml.Unmarshal(byteValue, &executeResponse); err != nil {
		log.Fatal(err)
	}

	if executeResponse.Status.ProcessSucceeded.Local != "" ||
		executeResponse.Status.ProcessFailed.Local != "" {

		// Determine the EndTime of the executeResponse by inspecting
		// the creation time of the responseFile
		stat, _ := os.Stat(responseFile)
		executeResponse.Status.EndTime = stat.ModTime()

		// Determine the WorkflowIdentifier, based on the name of the parent
		// directory
		executeResponse.Process.WorkflowIdentifier = filepath.Base(filepath.Dir(responseFile))

		// Determine the value for executeResponse.Status.ProcessStatus
		if executeResponse.Status.ProcessSucceeded.Local != "" {
			executeResponse.Status.ProcessStatus = 0
		} else {
			executeResponse.Status.ProcessStatus = 1
		}
	}

	return &executeResponse
}
