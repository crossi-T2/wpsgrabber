package wpsgrabber

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func createCSV(filename string, executeResponse *ExecuteResponse) error {

	loc, _ := time.LoadLocation("UTC")

	processIdentifier := executeResponse.Process.Identifier
	if configuration.ProcessIdentifier != "" {
		processIdentifier = configuration.ProcessIdentifier
	}

	processVersion := executeResponse.Process.Version
	if configuration.ProcessVersion != "" {
		processVersion = configuration.ProcessVersion
	}

	records := [][]string{
		{processIdentifier,
			processVersion,
			executeResponse.Process.WorkflowIdentifier,
			"not used",
			executeResponse.Status.CreationTime.In(loc).Format(time.RFC3339),
			executeResponse.Status.EndTime.In(loc).Format(time.RFC3339),
			strconv.Itoa(executeResponse.Status.ProcessStatus)},
	}

	file, err := os.Create(filename)
	if err != nil {
		err = fmt.Errorf("can't create file %s : %v ", filename, err)
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)

	// Write any buffered data to the underlying writer
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			err = fmt.Errorf("can't write record to CSV %s: %v", filename, err)
			return err
		}
	}

	return nil
}
