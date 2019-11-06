package wpsgrabber

import (
	"encoding/csv"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func createCSV(executeResponse *ExecuteResponse) error {

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

	filename := filepath.Join(configuration.CSVOutputDir, uuid.New().String()+".csv")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	w := csv.NewWriter(file)

	// Write any buffered data to the underlying writer
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Print("error writing record to csv:", err)
			return err
		}
	}

	log.Println("CSV created:", filename)

	return nil
}
