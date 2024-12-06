package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"zenprogramming.it/ado-uat-generator/models"
)

// Reads UAT resuls with expected schema from provided file
func ReadUatResults(file string) (
	[]models.UatResult, error) {

	// Get current working directory and compose full file path
	currentDir, _ := os.Getwd()
	fullFilePath := path.Join(currentDir, file)

	// Read file content and return error if file cannot be read
	content, err := os.ReadFile(fullFilePath)
	if err != nil {
		return nil, fmt.Errorf("File %v not found", fullFilePath)
	}

	// Deserialize results file
	var data = []models.UatResult{}
	jsonErr := json.Unmarshal([]byte(content), &data)

	if jsonErr != nil {
		return nil, fmt.Errorf(`File '%v' do not contains a valid JSON: %v`, fullFilePath, data)
	}

	return data, nil
}
