package model

type ProductIDs struct {
	List []string `json:"articles_sort_list"`
}

type data struct {
	Value string `json:"value"`
}

type sizeData struct {
	Body   map[string]map[string]data `json:"body"`
	Header map[string]map[string]data `json:"header"`
}

type SizeTale struct {
	SizeChart map[string]sizeData `json:"size_chart"`
}

type ReviewDetails struct {
	Date        string
	Rating      string
	Title       string
	Description string
	ReviewerID  string
}

type Review struct {
	Rating                string
	NumberOfReviews       string
	RecommendedRate       string
	SenseOfFitting        string
	AppropriationOfLength string
	QualityOfMaterial     string
	Comfort               string
	Details               []ReviewDetails
}

type DescriptionDetails struct {
	Title       string
	General     string
	Itemization string
}

type Product struct {
	ID              string
	Model           string
	URL             string
	Breadcrumb      string
	ImageURL        []string
	Category        string
	Name            string
	Price           string
	Currency        string
	AvailableSize   []string
	SenseOfSize     string
	Description     DescriptionDetails
	TaleOfSize      SizeTale
	SpecialFunction string
	Review          Review
	KWs             []string
}
