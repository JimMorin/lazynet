package actions

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/sureshkrishnan-v/lazynet/internal/models"
)

// ExportToCSV saves the provided list of connections to a CSV file.
// It returns the filename or an error.
func ExportToCSV(connections []models.Connection) (string, error) {
	filename := "lazynet_export.csv"
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	headers := []string{"ProcessName", "PID", "Status", "Family", "Type", "LocalAddr", "RemoteAddr"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// Write data rows
	for _, c := range connections {
		row := []string{
			c.ProcessName,
			strconv.Itoa(int(c.Pid)),
			c.Status,
			c.Family,
			c.Type,
			c.Laddr,
			c.Raddr,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return filename, nil
}
