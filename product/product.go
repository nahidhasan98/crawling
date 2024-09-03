package product

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nahidhasan98/crawling/helper"
	"github.com/nahidhasan98/crawling/model"
)

func GatherIDs(limit int) []string {
	var productIDs []string

	apiURL := "https://shop.adidas.jp/f/v1/pub/product/list"
	page := 1

	for {
		URL := fmt.Sprintf("%s?gender=mens&limit=120&page=%d", apiURL, page)

		response := helper.GETRequest(URL)
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		helper.ErrorCheck(err)

		var tempList model.ProductIDs
		err = json.Unmarshal(responseBody, &tempList)
		helper.ErrorCheck(err)

		productIDs = append(productIDs, tempList.List...)

		if len(productIDs) >= limit {
			break
		}
		page = page + 1

	}
	if len(productIDs) > limit {
		productIDs = productIDs[:limit]
	}

	return productIDs
}

func getBreadcrumb(doc *goquery.Document, product *model.Product) string {
	productModel := ""
	breadcrumb := ""

	doc.Find(".breadcrumbListItem a").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			if len(breadcrumb) > 0 {
				breadcrumb += " / "
			}
			breadcrumb += strings.TrimSpace(s.Text())

			model, exist := s.Attr("href")
			if exist {
				if strings.HasPrefix(strings.TrimSpace(model), "/model/") {
					model = strings.TrimPrefix(model, "/model/")
					model = strings.TrimSuffix(model, "/")
				}
			}
			productModel = model
		}
	})
	product.Model = productModel

	return breadcrumb
}

func getImageURL(doc *goquery.Document, host string) []string {
	imageURL := []string{}

	body := doc.Find("script#__NEXT_DATA__").Text()

	var bodyInterfacer map[string]interface{}
	err := json.Unmarshal([]byte(body), &bodyInterfacer)
	helper.ErrorCheck(err)

	props := bodyInterfacer["props"].(map[string]interface{})
	pageProps := props["pageProps"].(map[string]interface{})
	apis := pageProps["apis"].(map[string]interface{})
	pdpInitialProps := apis["pdpInitialProps"].(map[string]interface{})
	detailApi := pdpInitialProps["detailApi"].(map[string]interface{})
	product := detailApi["product"].(map[string]interface{})
	article := product["article"].(map[string]interface{})
	image := article["image"].(map[string]interface{})
	details := image["details"].([]interface{})

	for _, v := range details {
		detail := v.(map[string]interface{})
		imageUrl := detail["imageUrl"].(map[string]interface{})
		large := imageUrl["large"].(string)

		imageURL = append(imageURL, fmt.Sprintf("%s%s", host, large))
	}

	return imageURL
}

func getCategory(doc *goquery.Document) string {
	category := strings.TrimSpace(doc.Find(".groupName").Text())
	return category
}

func getName(doc *goquery.Document) string {
	name := strings.TrimSpace(doc.Find(".itemTitle").Text())
	return name
}

func getPrice(doc *goquery.Document, product *model.Product) string {
	price := strings.TrimSpace(doc.Find(".price-value ").Text())
	price = strings.ReplaceAll(price, ",", "")
	product.Currency = "¥"

	return price
}

func getAvailableSize(doc *goquery.Document) []string {
	sizes := []string{}
	doc.Find(".sizeSelectorListItem button").Each(func(i int, s *goquery.Selection) {
		size := strings.TrimSpace(s.Text())
		sizes = append(sizes, size)
	})
	return sizes
}

func getSenseOfSize(doc *goquery.Document, responseBody []byte) string {
	percentage := ""
	classes, exist := doc.Find(".bar .marker").Attr("class")
	if exist {
		classList := strings.Split(classes, " ")

		for _, class := range classList {
			if strings.HasPrefix(class, "mod-marker_") {
				pattern := fmt.Sprintf(`\.bar \.marker\.%s\{[^}]*left:(\d+(\.\d+)?)%%;`, regexp.QuoteMeta(class))
				regex := regexp.MustCompile(pattern)
				match := regex.FindStringSubmatch(string(responseBody))

				if len(match) > 0 {
					percentage = strings.TrimSpace(match[1])
				}

				break
			}
		}
	}

	if len(percentage) > 0 {
		percentage = fmt.Sprintf("Appropriate: %s%%", percentage)
	}

	return percentage
}

func getDescription(doc *goquery.Document) model.DescriptionDetails {
	title := strings.TrimSpace(doc.Find(".itemFeature").Text())
	general := strings.TrimSpace(doc.Find(".description_part.details").Text())
	itemization := ""

	doc.Find(".articleFeaturesItem").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			itemization += "\r\n"
		}
		itemization += "• " + strings.TrimSpace(s.Text())
	})

	description := model.DescriptionDetails{
		Title:       title,
		General:     general,
		Itemization: itemization,
	}

	return description
}

func getTaleOfSize(productModel string) model.SizeTale {
	apiURL := "https://shop.adidas.jp/f/v1/pub/size_chart"
	URL := fmt.Sprintf("%s/%s", apiURL, productModel)

	response := helper.GETRequest(URL)
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	helper.ErrorCheck(err)

	var sizeTale model.SizeTale
	err = json.Unmarshal(responseBody, &sizeTale)
	helper.ErrorCheck(err)

	return sizeTale
}

func getSpecialFunction(doc *goquery.Document) string {
	specialFunction := ""

	title := strings.TrimSpace(doc.Find(".item_part.details a").Text())
	description := strings.TrimSpace(strings.TrimPrefix(doc.Find(".item_part.details").Text(), title))

	if len(title) > 0 && len(description) > 0 {
		specialFunction = fmt.Sprintf("[ %s ] %s", title, description)
	} else if len(title) > 0 {
		specialFunction = title
	} else if len(description) > 0 {
		specialFunction = description
	}

	return specialFunction
}

func getReview(productID, productModel string) model.Review {
	apiURL := "https://adidasjp.ugc.bazaarvoice.com/7896-ja_jp"
	URL := fmt.Sprintf("%s/%s/reviews.djs?format=embeddedhtml&productattribute_itemKcod=%s", apiURL, productModel, productID)

	response := helper.GETRequest(URL)
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	helper.ErrorCheck(err)

	regex := regexp.MustCompile(`materials\s*=\s*\{\s*"BVRRRatingSummarySourceID":\s*"(.*?)"\s*\}`)
	match := regex.FindStringSubmatch(string(responseBody))

	var bodyReader *strings.Reader
	if len(match) > 1 {
		htmlContent := match[1]

		replacer := strings.NewReplacer(
			`\n`, "", // Remove newline characters
			`\r`, "", // Remove carriage return characters
			`\"`, `"`, // Unescape double quotes
			`\/`, "/", // Unescape forward slash
		)

		cleanedHTML := replacer.Replace(htmlContent)
		bodyReader = strings.NewReader(cleanedHTML)
	}

	rating, numberOfReviews, recommendedRate, senseOfFitting, appropriationOfLength, qualityOfMaterial, comfort := "", "", "", "", "", "", ""
	reviewDetails := []model.ReviewDetails{}
	review := model.Review{
		Details: reviewDetails,
	}

	if bodyReader == nil {
		return review
	}

	document, err := goquery.NewDocumentFromReader(bodyReader)
	helper.ErrorCheck(err)

	rating = strings.TrimSpace(document.Find("#BVRRWidgetID #BVRRRatingOverall_ .BVRRRatingNumber").Text())
	numberOfReviews = strings.TrimSpace(document.Find("#BVRRWidgetID  .BVRRBuyAgainTotal").Text())
	recommendedRate = strings.TrimSpace(document.Find("#BVRRWidgetID  .BVRRBuyAgainPercentage").Text())
	fitting, exist := document.Find(".BVRRSecondaryRatingsContainer .BVRRRatingFit .BVRRRatingRadioImage img").Attr("title")
	if exist {
		senseOfFitting = strings.TrimSpace(fitting)
	}

	length, exist := document.Find(".BVRRSecondaryRatingsContainer .BVRRRatingLength .BVRRRatingRadioImage img").Attr("title")
	if exist {
		appropriationOfLength = strings.TrimSpace(length)
	}

	quality, exist := document.Find(".BVRRSecondaryRatingsContainer .BVRRRatingQuality .BVRRRatingRadioImage img").Attr("title")
	if exist {
		qualityOfMaterial = strings.TrimSpace(quality)
	}

	comfortness, exist := document.Find(".BVRRSecondaryRatingsContainer .BVRRRatingComfort .BVRRRatingRadioImage img").Attr("title")
	if exist {
		comfort = strings.TrimSpace(comfortness)
	}

	document.Find("#BVRRWidgetID #BVRRDisplayContentBodyID .BVRRContentReview").Each(func(i int, s *goquery.Selection) {
		rating := fmt.Sprintf("%s %s %s", strings.TrimSpace(s.Find(".BVRRRatingNormalOutOf .BVRRRatingNumber").Text()), strings.TrimSpace(s.Find(".BVRRRatingNormalOutOf .BVRRSeparatorText").Text()), strings.TrimSpace(s.Find(".BVRRRatingNormalOutOf .BVRRRatingRangeNumber").Text()))

		title := strings.TrimSpace(s.Find(".BVRRReviewTitle").Text())
		date := strings.TrimSpace(s.Find(".BVRRReviewDate").Text())
		description := strings.TrimSpace(s.Find(".BVRRReviewText").Text())
		id := strings.TrimSpace(s.Find(".BVRRNickname").Text())

		temp := model.ReviewDetails{
			Date:        date,
			Rating:      rating,
			Title:       title,
			Description: description,
			ReviewerID:  id,
		}
		reviewDetails = append(reviewDetails, temp)
	})

	review = model.Review{
		Rating:                rating,
		NumberOfReviews:       numberOfReviews,
		RecommendedRate:       recommendedRate,
		SenseOfFitting:        senseOfFitting,
		AppropriationOfLength: appropriationOfLength,
		QualityOfMaterial:     qualityOfMaterial,
		Comfort:               comfort,
		Details:               reviewDetails,
	}

	return review
}

func getKWs(doc *goquery.Document) []string {
	var kws []string
	doc.Find(".itemTagsPosition .inner a").Each(func(i int, s *goquery.Selection) {
		kw := strings.TrimSpace(s.Text())
		kws = append(kws, kw)
	})

	return kws
}

func GetDetails(productID string) *model.Product {
	host := "https://shop.adidas.jp"
	URL := fmt.Sprintf("%s/products/%s/", host, productID)

	response := helper.GETRequest(URL)
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	helper.ErrorCheck(err)

	bodyReader := strings.NewReader(string(responseBody))

	document, err := goquery.NewDocumentFromReader(bodyReader)
	helper.ErrorCheck(err)

	product := model.Product{
		ID:  productID,
		URL: URL,
	}

	product.Breadcrumb = getBreadcrumb(document, &product)
	product.ImageURL = getImageURL(document, host)
	product.Category = getCategory(document)
	product.Name = getName(document)
	product.Price = getPrice(document, &product)
	product.AvailableSize = getAvailableSize(document)
	product.SenseOfSize = getSenseOfSize(document, responseBody)
	product.Description = getDescription(document)
	product.TaleOfSize = getTaleOfSize(product.Model)
	product.SpecialFunction = getSpecialFunction(document)
	product.Review = getReview(product.ID, product.Model)
	product.KWs = getKWs(document)

	return &product
}
