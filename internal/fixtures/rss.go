package fixtures

type RssProps struct {
	Title               string
	Link                string
	Description         string
	AtomLink            string
	LastBuildDateRFC822 string
	Items               []RssItem
}

type RssItem struct {
	Title         string
	Link          string
	Description   string
	PubDateRFC822 string
	Guid          string
	Category      []string
}
