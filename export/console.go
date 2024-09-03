package export

import (
	"fmt"

	"github.com/nahidhasan98/crawling/model"
)

// PrintToConsole prints the details of a product to the console.
// It takes a pointer to a Product struct as its parameter.
func PrintToConsole(products *model.Product) {
	fmt.Println(products)
	fmt.Println()
}
