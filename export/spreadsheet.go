package export

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/nahidhasan98/crawling/helper"
	"github.com/nahidhasan98/crawling/model"
	"github.com/xuri/excelize/v2"
)

// createNewFileFromTemplate copies the content of the source file to a new destination file.
// It takes the source and destination file paths as parameters and returns an error if the operation fails.
func createNewFileFromTemplate(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// writeReviewDetails writes the review details of a product to an Excel sheet.
// It takes a slice of ReviewDetails and a serial number as parameters.
// The function returns the top-left and bottom-right cell references of the written data.
func writeReviewDetails(reviewDetails []model.ReviewDetails, serial int) (string, string) {
	filePath := "product.xlsx"
	f, err := excelize.OpenFile(filePath)
	helper.ErrorCheck(err)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	reviewSheet := "Review"
	rows, err := f.GetRows(reviewSheet)
	helper.ErrorCheck(err)

	nextRow := len(rows) + 1
	startRow := nextRow

	f.SetCellValue(reviewSheet, "A"+strconv.Itoa(nextRow), serial+1)

	for _, v := range reviewDetails {
		f.SetCellValue(reviewSheet, "B"+strconv.Itoa(nextRow), v.Date)
		f.SetCellValue(reviewSheet, "C"+strconv.Itoa(nextRow), v.Rating)
		f.SetCellValue(reviewSheet, "D"+strconv.Itoa(nextRow), v.Title)
		f.SetCellValue(reviewSheet, "E"+strconv.Itoa(nextRow), v.Description)
		f.SetCellValue(reviewSheet, "F"+strconv.Itoa(nextRow), v.ReviewerID)
		nextRow++
	}

	topLeft := fmt.Sprintf("A%d", startRow)
	bottomRight := fmt.Sprintf("A%d", nextRow-1)

	if startRow > 0 && startRow < nextRow-1 {
		err := f.MergeCell(reviewSheet, topLeft, bottomRight)
		helper.ErrorCheck(err)

		style, err := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})

		helper.ErrorCheck(err)

		f.SetCellStyle(reviewSheet, topLeft, bottomRight, style)
	}

	err = f.Save()
	helper.ErrorCheck(err)

	bottomRight = fmt.Sprintf("F%d", nextRow-1)
	return topLeft, bottomRight
}

// writeTaleOfSize writes the TaleOfSize details of a product to an Excel sheet.
// It takes a SizeTale struct and a serial number as parameters.
// The function returns the top-left and bottom-right cell references of the written data.
func writeTaleOfSize(taleOfSize model.SizeTale, serial int) (string, string) {
	filePath := "product.xlsx"
	f, err := excelize.OpenFile(filePath)
	helper.ErrorCheck(err)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sizeSheet := "TaleOfSize"

	rows, err := f.GetRows(sizeSheet)
	helper.ErrorCheck(err)

	nextRow := len(rows) + 1
	startRow := nextRow

	f.SetCellValue(sizeSheet, "A"+strconv.Itoa(nextRow), serial+1)

	header := taleOfSize.SizeChart["0"].Header["0"]

	keys := make([]string, 0, len(header))
	for key := range header {
		keys = append(keys, key)
	}
	// Sort the keys
	sort.Strings(keys)
	// Iterate over the sorted keys and print the key-value pairs
	for _, key := range keys {
		// fmt.Printf("%s: %s\n", key, header[key])
		f.SetCellValue(sizeSheet, "B"+strconv.Itoa(nextRow), header[key].Value)
		nextRow++
	}

	body := taleOfSize.SizeChart["0"].Body
	col := "C"
	keys2 := make([]string, 0, len(body))
	for key2 := range body {
		keys2 = append(keys2, key2)
	}
	// Sort the keys
	sort.Strings(keys2)
	// Iterate over the sorted keys and print the key-value pairs
	for _, key2 := range keys2 {
		// fmt.Printf("%s: %s\n", key, body[key])
		nextRow = startRow

		keys3 := make([]string, 0, len(body[key2]))
		for key3 := range body[key2] {
			keys3 = append(keys3, key3)
		}
		// Sort the keys
		sort.Strings(keys3)
		// Iterate over the sorted keys and print the key-value pairs
		for _, key3 := range keys3 {
			f.SetCellValue(sizeSheet, fmt.Sprintf("%s%d", col, nextRow), body[key2][key3].Value)
			nextRow++
		}
		colRune := []rune(col)[0]
		colRune++
		col = string(colRune)
	}

	if startRow > 0 && startRow < nextRow-1 {
		topLeft := fmt.Sprintf("A%d", startRow)
		bottomRight := fmt.Sprintf("A%d", nextRow-1)
		err := f.MergeCell(sizeSheet, topLeft, bottomRight)
		helper.ErrorCheck(err)

		style, err := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})

		helper.ErrorCheck(err)

		f.SetCellStyle(sizeSheet, topLeft, bottomRight, style)
	}

	err = f.Save()
	helper.ErrorCheck(err)

	colRune := []rune(col)[0]
	colRune--
	col = string(colRune)

	topLeft := fmt.Sprintf("A%d", startRow)
	bottomRight := fmt.Sprintf("%s%d", col, nextRow-1)
	return topLeft, bottomRight
}

// prepareImageURL formats a slice of image URLs into a numbered list as a string.
func prepareImageURL(imageURLs []string) string {
	res := ""

	for i := 0; i < len(imageURLs); i++ {
		res += fmt.Sprintf("%d. %s\n", i+1, imageURLs[i])
	}

	return res
}

// prepareAvailableSize formats a slice of available sizes into a comma-separated string.
func prepareAvailableSize(availableSize []string) string {
	res := ""

	for i := 0; i < len(availableSize); i++ {
		res += availableSize[i]

		if i < len(availableSize)-1 {
			res += ", "
		}
	}

	return res
}

// prepareKWs formats a slice of keywords into a comma-separated string.
func prepareKWs(kws []string) string {
	res := ""

	for i := 0; i < len(kws); i++ {
		res += kws[i]

		if i < len(kws)-1 {
			res += ", "
		}
	}

	return res
}

// Spreadsheet creates an Excel file containing product details using data from a slice of Product structs.
// It uses a template Excel file, writes product data to the file, and saves it as "product.xlsx".
// The function returns an error if any operation fails during file creation or data writing.
func Spreadsheet(products []model.Product) error {
	src := "./template/template.xlsx"
	dst := "./product.xlsx"
	filePath := "product.xlsx"

	err := createNewFileFromTemplate(src, dst)
	if err != nil {
		return err
	}

	for i := 0; i < len(products); i++ {
		topLeft, bottomRight := writeTaleOfSize(products[i].TaleOfSize, i)
		topLeft2, bottomRight2 := writeReviewDetails(products[i].Review.Details, i)

		f, err := excelize.OpenFile(filePath)
		if err != nil {
			return err
		}

		basicSheet := "Basic"
		rows, err := f.GetRows(basicSheet)
		if err != nil {
			return err
		}
		nextRow := len(rows) + 1

		f.SetCellValue(basicSheet, "A"+strconv.Itoa(nextRow), i+1)
		f.SetCellValue(basicSheet, "B"+strconv.Itoa(nextRow), products[i].URL)
		f.SetCellValue(basicSheet, "C"+strconv.Itoa(nextRow), products[i].Breadcrumb)
		f.SetCellValue(basicSheet, "D"+strconv.Itoa(nextRow), products[i].Category)
		f.SetCellValue(basicSheet, "E"+strconv.Itoa(nextRow), products[i].Name)
		f.SetCellValue(basicSheet, "F"+strconv.Itoa(nextRow), fmt.Sprintf("%s %s", products[i].Currency, products[i].Price))

		imageURL := prepareImageURL(products[i].ImageURL)
		f.SetCellValue(basicSheet, "G"+strconv.Itoa(nextRow), imageURL)

		style, err := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Vertical: "top",
				WrapText: true,
			},
		})
		if err != nil {
			return err
		}
		f.SetCellStyle(basicSheet, "G"+strconv.Itoa(nextRow), "G"+strconv.Itoa(nextRow), style)

		availableSize := prepareAvailableSize(products[i].AvailableSize)
		f.SetCellValue(basicSheet, "H"+strconv.Itoa(nextRow), availableSize)
		f.SetCellValue(basicSheet, "I"+strconv.Itoa(nextRow), products[i].SenseOfSize)

		f.SetCellValue(basicSheet, "J"+strconv.Itoa(nextRow), products[i].Description.Title)
		f.SetCellValue(basicSheet, "K"+strconv.Itoa(nextRow), products[i].Description.General)
		f.SetCellValue(basicSheet, "L"+strconv.Itoa(nextRow), products[i].Description.Itemization)

		link := fmt.Sprintf("TaleOfSize!%s:%s", topLeft, bottomRight)
		display, tooltip := "View Tale of Size", "Click to see Tale Of Size"
		f.SetCellHyperLink(basicSheet, "M"+strconv.Itoa(nextRow), link, "Location", excelize.HyperlinkOpts{
			Display: &display,
			Tooltip: &tooltip,
		})

		f.SetCellValue(basicSheet, "N"+strconv.Itoa(nextRow), products[i].SpecialFunction)
		f.SetCellValue(basicSheet, "O"+strconv.Itoa(nextRow), products[i].Review.Rating)
		f.SetCellValue(basicSheet, "P"+strconv.Itoa(nextRow), products[i].Review.NumberOfReviews)
		f.SetCellValue(basicSheet, "Q"+strconv.Itoa(nextRow), products[i].Review.RecommendedRate)
		f.SetCellValue(basicSheet, "R"+strconv.Itoa(nextRow), products[i].Review.SenseOfFitting)
		f.SetCellValue(basicSheet, "S"+strconv.Itoa(nextRow), products[i].Review.AppropriationOfLength)
		f.SetCellValue(basicSheet, "T"+strconv.Itoa(nextRow), products[i].Review.QualityOfMaterial)
		f.SetCellValue(basicSheet, "U"+strconv.Itoa(nextRow), products[i].Review.Comfort)

		link2 := fmt.Sprintf("Review!%s:%s", topLeft2, bottomRight2)
		display2, tooltip2 := "View Review Details", "Click to see Review Details"
		f.SetCellHyperLink(basicSheet, "V"+strconv.Itoa(nextRow), link2, "Location", excelize.HyperlinkOpts{
			Display: &display2,
			Tooltip: &tooltip2,
		})

		kws := prepareKWs(products[i].KWs)
		f.SetCellValue(basicSheet, "W"+strconv.Itoa(nextRow), kws)

		err = f.Save()
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	fmt.Println("Data exported to", filePath, "successfully.")
	return nil
}
