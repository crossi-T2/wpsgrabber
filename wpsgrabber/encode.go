package wpsgrabber

import (
	"encoding/csv"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func createCSV(executeResponse *ExecuteResponse) {

	// TODO: if we add the request.xml in "not used" we would need to encode commas

	records := [][]string{
		{configuration.ProcessorCode, configuration.ProcessorVersion, "wpsid", "not used", "start", "stop", strconv.Itoa(executeResponse.Status.ProcessStatus)},
	}

	file, err := os.Create(filepath.Join(configuration.CSVOutputDir, uuid.New().String()+".csv"))
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
}
