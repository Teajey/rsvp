package fixtures

type RssProps struct {
	XMLName  string `xml:"rss"`
	Version  string `xml:"version,attr"`
	Atom     string `xml:"xmlns:atom,attr"`
	Channels []RssChannel
}

type RssAtomLink struct {
	XMLName string `xml:"atom:link"`
	Href    string `xml:"href,attr"`
	Rel     string `xml:"rel,attr"`
	Type    string `xml:"type,attr"`
}

type RssChannel struct {
	XMLName             string `xml:"channel"`
	Title               string `xml:"title"`
	Link                string `xml:"link"`
	Description         string `xml:"description"`
	Language            string `xml:"language"`
	LastBuildDateRFC822 string `xml:"lastBuildDate"`
	AtomLink            RssAtomLink
	Items               []RssItem
}

type RssItem struct {
	XMLName       string   `xml:"item"`
	Title         string   `xml:"title"`
	Link          string   `xml:"link,omitempty"`
	Description   string   `xml:"description"`
	PubDateRFC822 string   `xml:"pubDate"`
	Guid          string   `xml:"guid"`
	Category      []string `xml:"category"`
}
