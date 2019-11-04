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

func createCSV(executeResponse *ExecuteResponse) {

	loc, _ := time.LoadLocation("UTC")

	records := [][]string{
		{configuration.ProcessorCode,
			configuration.ProcessorVersion,
			executeResponse.Process.CurrentIdentifier,
			"not used",
			executeResponse.Status.CreationTime.In(loc).Format(time.RFC3339),
			executeResponse.Status.EndTime.In(loc).Format(time.RFC3339),
			strconv.Itoa(executeResponse.Status.ProcessStatus)},
	}

	filename := filepath.Join(configuration.CSVOutputDir, uuid.New().String()+".csv")
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)

	// Write any buffered data to the underlying writer
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	log.Println("CSV created:", filename)
}
