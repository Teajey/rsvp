package rsvp_test

import (
	html "html/template"
	"net/http/httptest"
	"testing"
	text "text/template"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
	"github.com/Teajey/rsvp/internal/fixtures"
)

func TestStringBody(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	err := res.Write(rec, req, rsvp.DefaultConfig())
	if err != nil {
		t.Fatalf("Write failed: %s", err)
	}
	statusCode := rec.Result().StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	s := rec.Body.String()
	assert.Eq(t, "body contents", body, s)
}

func TestStringBodyAcceptApp(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/*")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "\"Hello,\\nWorld!\""+"\n", s)
}

func TestStringBodyJsonExt(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/message.json", nil)

	// Even if Accept is set, the file extension takes precedence
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()

	// Note Config.ExtensionToProposalMap must be set for this to work
	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "\"Hello,\\nWorld!\""+"\n", s)
}

func TestUnsupportedExtIgnored(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}

	// If the extension isn't supported, it's ignored
	req := httptest.NewRequest("GET", "/message.blah", nil)

	req.Header.Set("Accept", "application/*")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "\"Hello,\\nWorld!\""+"\n", s)
}

func TestListBody(t *testing.T) {
	body := []string{"hello", "world", "123"}
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	statusCode := rec.Result().StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	s := rec.Body.String()
	assert.Eq(t, "body contents", `["hello","world","123"]`+"\n", s)
}

func TestBytesBody(t *testing.T) {
	body := []byte("Hello, World!")
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/*")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "application/octet-stream", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "Hello, World!", s)
}

func TestHtmlTemplate(t *testing.T) {
	body := "Hello <input> World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.HtmlTemplate = html.New("")
	cfg.HtmlTemplate = html.Must(cfg.HtmlTemplate.Parse(`{{define "tm"}}<div>{{if .}}{{.}}{{else}}Nothin!{{end}}</div>{{end}}`))
	cfg.TextTemplate = text.New("")
	cfg.TextTemplate = text.Must(cfg.TextTemplate.Parse(`{{define "tm"}}{{if .}}Message: {{.}}{{else}}Nothin!{{end}}{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "text/html; charset=utf-8", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "<div>Hello &lt;input&gt; World!</div>", s)
}

func TestTextTemplate(t *testing.T) {
	body := "Hello, World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/plain")
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.HtmlTemplate = html.New("")
	cfg.HtmlTemplate = html.Must(cfg.HtmlTemplate.Parse(`{{define "tm"}}<div>{{if .}}{{.}}{{else}}Nothin!{{end}}</div>{{end}}`))
	cfg.TextTemplate = text.New("")
	cfg.TextTemplate = text.Must(cfg.TextTemplate.Parse(`{{define "tm"}}{{if .}}Message: {{.}}{{else}}Nothin!{{end}}{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "Message: Hello, World!", s)
}

func TestRss(t *testing.T) {
	body := fixtures.RssProps{
		Title:               "My stuff",
		Description:         "It's cool",
		LastBuildDateRFC822: "April 1st",
		Items: []fixtures.RssItem{
			{
				Title:         "New post",
				Description:   "I got a pet",
				PubDateRFC822: "April 1st",
				Guid:          "123",
			},
		},
	}

	res := rsvp.Response{Body: body, TemplateName: "rss.gotmpl"}
	req := httptest.NewRequest("GET", "/posts.xml", nil)
	rec := httptest.NewRecorder()
	// Override what rsvp will set
	rec.Header().Set("Content-Type", "application/rss+xml")

	cfg := rsvp.DefaultConfig()
	// This will make sure rsvp uses the text/html rendering scheme which I use as a hack to avoid rendering non-XML. I should probably just use encoding/xml
	cfg.ExtToProposalMap["xml"] = "text/html"

	htmlTmpl, err := html.ParseFiles("./internal/fixtures/rss.gotmpl")
	assert.FatalErr(t, "loading rss template", err)
	cfg.HtmlTemplate = htmlTmpl

	cfg.TextTemplate = text.New("")
	cfg.TextTemplate = text.Must(cfg.TextTemplate.Parse(`{{define "rss.gotmpl"}}{{if .}}Message: {{.}}{{else}}Nothin!{{end}}{{end}}`))

	err = res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "application/rss+xml", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
   <channel>
      <title>My stuff</title>
      <link></link>
      <description>It&#39;s cool</description>
      <language>en</language>
      <lastBuildDate>April 1st</lastBuildDate>
      <atom:link href="" rel="self" type="application/rss+xml" />
      <item>
         <title>New post</title>
         <description>I got a pet</description>
         <pubDate>April 1st</pubDate>
         <guid>123</guid>
      </item>
   </channel>
</rss>
`, s)
}
