package main

import (
	"fmt"

	"github.com/nahidhasan98/crawling/export"
	"github.com/nahidhasan98/crawling/model"
	"github.com/nahidhasan98/crawling/product"
)

func main() {
	fmt.Println("Programming is running...")

	productIDs := product.GatherIDs(300)

	var products []model.Product
	for i := 0; i < len(productIDs); i++ {
		fmt.Println("Getting product", i+1, ":")

		tempProduct := product.GetDetails(productIDs[i])
		products = append(products, *tempProduct)
	}

	fmt.Println("Writting data to file...")
	err1 := export.WriteToFile(products)
	if err1 != nil {
		fmt.Println("Error exporting to Google Sheet:", err1)
	}

	fmt.Println("Exporting data to Spreadsheet...")
	err := export.Spreadsheet(products)
	if err != nil {
		fmt.Println("Error exporting to Spreadsheet:", err)
	}
}
