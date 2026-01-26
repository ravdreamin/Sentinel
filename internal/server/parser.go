package server

import (
	"fmt"
	"strings"
)

func processFile(filepath string) ([]string, error) {

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt":
		return parseTextFile(filePath)
	case ".csv":
		return parseCSVFile(filePath)
	case ".json":
		return parseJSONFile(filePath)
	case ".pdf":
		return parsePDFFile(filePath)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}
