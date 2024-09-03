package export

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nahidhasan98/crawling/model"
)

// WriteToFile serializes a slice of Product structs to JSON format and writes it to a file.
// The function returns an error if any file operation or JSON marshaling fails.
func WriteToFile(product []model.Product) error {
	filename := "product.txt"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serialize the struct to JSON
	jsonData, err := json.MarshalIndent(product, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("Data written to", filename)
	return nil
}
