package server

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

func processFile(filePath string) ([]string, error) {

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

func parseTextFile(filePath string) ([]string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error parsing your txt file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	datas := []string{}
	for scanner.Scan() == true {
		datas = append(datas, scanner.Text())

	}
	return datas, nil
}

func parseCSVFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error parsing your txt file: %s", err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error parsing your csv file: %s", err)
	}

	datas := []string{}
	for _, data := range records {
		datas = append(datas, data[0])

	}

	return datas, nil

}

func parseJSONFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error parsing your json file: %s", err)
	}
	defer file.Close()

	datas := []string{}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&datas)
	if err != nil {
		return nil, fmt.Errorf("Error decoding json: %s", err)
	}

	return datas, nil

}

func parsePDFFile(filePath string) ([]string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error parsing pdf: %s", err)
	}

	defer f.Close()
	datas := []string{}
	for i := 1; i <= r.NumPage(); i++ {
		record := r.Page(i)
		text, err := record.GetPlainText(nil)
		if err != nil {

			return nil, fmt.Errorf("Error reading page %d: %s", i, err)
		}
		datas = append(datas, text)
	}

	return datas, nil

}
