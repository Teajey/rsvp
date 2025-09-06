package rsvp_test

import (
	"bytes"
	html "html/template"
	"net/http"
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
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
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

func TestUnsupportedExtHasBlank404(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}

	req := httptest.NewRequest("GET", "/message.blah", nil)

	req.Header.Set("Accept", "application/*")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 404, statusCode)

	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestListBody(t *testing.T) {
	body := []string{"hello", "world", "123"}
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Status code", 200, statusCode)
	s := rec.Body.String()
	assert.Eq(t, "body contents", `["hello","world","123"]`+"\n", s)
}

func TestBytesBody(t *testing.T) {
	body := []byte{0x29, 0x46, 0x4c, 0xff, 0x2f, 0x0e, 0x62, 0x41, 0xb5, 0xe3, 0xbb, 0xff, 0x06, 0x89, 0xa2, 0xef, 0xf0, 0xe2, 0x90, 0x4b, 0x62, 0x93, 0xa2, 0x6c, 0xc9, 0xcf, 0x08, 0xae, 0x18, 0xb0, 0xc2, 0xfc}
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

	s := rec.Body.Bytes()
	assert.SlicesEq(t, "body contents", body, s)
}

func TestHtmlTemplate(t *testing.T) {
	body := "Hello <input> World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")
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

func TestTextTemplateWithName(t *testing.T) {
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

func TestTextTemplateWithoutName(t *testing.T) {
	body := "Hello, World!"
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/message.txt", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.TextTemplate = text.New("")
	cfg.TextTemplate = text.Must(cfg.TextTemplate.Parse(`{{define "tm"}}{{if .}}Message: {{.}}{{else}}Nothin!{{end}}{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "Hello, World!", s)
}

func TestAttemptToRenderNonTextAsText(t *testing.T) {
	body := map[string]string{"I'm": "a map"}
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/message.txt", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 404, statusCode)

	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestRss(t *testing.T) {
	body := fixtures.RssProps{
		Version: "2.0",
		Atom:    "http://www.w3.org/2005/Atom",
		Channels: []fixtures.RssChannel{
			{
				Title:               "My stuff",
				Description:         "It's cool",
				Language:            "en",
				LastBuildDateRFC822: "April 1st",
				AtomLink: fixtures.RssAtomLink{
					Href: "",
					Rel:  "self",
					Type: "application/rss+xml",
				},
				Items: []fixtures.RssItem{
					{
						Title:         "New post",
						Description:   "I got a pet",
						PubDateRFC822: "April 1st",
						Guid:          "123",
					},
				},
			},
		},
	}

	res := rsvp.Response{Body: body, TemplateName: "rss.gotmpl"}
	// NOTE: Ideally, some middleware would be added to map .rss -> .xml, instead of using .rss.xml
	req := httptest.NewRequest("GET", "/posts.rss.xml", nil)
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/rss+xml")

	cfg := rsvp.DefaultConfig()

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)

	assert.Eq(t, "Content type", "application/rss+xml", resp.Header.Get("Content-Type"))

	s := rec.Body.String() + "\n"
	assert.TextSnapshot(t, "rss.xml", s)
}

func TestNotFound(t *testing.T) {
	res := rsvp.Response{Body: "404 Not Found", Status: http.StatusNotFound}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "*/*")
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.HtmlTemplate = html.Must(html.New("").Parse(`{{define "tm"}}<div>{{if .}}{{.}}{{else}}Nothin!{{end}}</div>{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 404, statusCode)
	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "404 Not Found", s)
}

func TestBlankOk(t *testing.T) {
	res := rsvp.Ok()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestEmptyStringBody(t *testing.T) {
	res := rsvp.Response{Body: ""}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	// TODO: Is it right that Content-Type is set here? In this scenario, by default, I feel it's reasonable not to set it
	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestNilBody(t *testing.T) {
	res := rsvp.Response{Body: nil}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", `null`+"\n", s)
}

func TestSeeOtherCanRender(t *testing.T) {
	res := rsvp.SeeOther("/", "POST successful")
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusSeeOther, statusCode)
	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", res.Body.(string), s)
}

func TestSeeOtherDoesNotRenderHtml(t *testing.T) {
	res := rsvp.SeeOther("/", nil)
	res.Html("<div></div>")
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusSeeOther, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestFoundCanRender(t *testing.T) {
	res := rsvp.Found("/", "POST successful")
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusFound, statusCode)
	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", res.Body.(string), s)
}

func TestFoundDoesNotRenderHtml(t *testing.T) {
	res := rsvp.Found("/", nil)
	res.Html("<div></div>")
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusFound, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestPermanentRedirectDoesNotRender(t *testing.T) {
	res := rsvp.PermanentRedirect("/")
	res.Body = "POST successful"
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusPermanentRedirect, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestMovedPermanentlyDoesNotRender(t *testing.T) {
	res := rsvp.MovedPermanently("/")
	res.Body = "POST successful"
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusMovedPermanently, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestNotFoundJson(t *testing.T) {
	res := rsvp.Response{Body: "404 Not Found", Status: http.StatusNotFound}
	req := httptest.NewRequest("GET", "/post.json", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.HtmlTemplate = html.Must(html.New("").Parse(`{{define "tm"}}<div>{{if .}}{{.}}{{else}}Nothin!{{end}}</div>{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 404, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", `"404 Not Found"`+"\n", s)
}

func TestHtmlFromString(t *testing.T) {
	res := rsvp.Ok()
	res.Html("<div></div>")
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "text/html; charset=utf-8", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", res.Body.(string), s)
}

func TestExplicitTextRequestWithoutTextTemplate(t *testing.T) {
	body := "Hello <input> World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/home.txt", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.HtmlTemplate = html.Must(html.New("").Parse(`{{define "tm"}}<div>{{if .}}{{.}}{{else}}Nothin!{{end}}</div>{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusNotFound, statusCode)

	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestExplicitHtmlRequestWithoutHtmlTemplate(t *testing.T) {
	body := "Hello <input> World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/home.html", nil)
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	cfg.TextTemplate = text.Must(text.New("").Parse(`{{define "tm"}}{{if .}}Message: {{.}}{{else}}Nothin!{{end}}{{end}}`))

	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusNotFound, statusCode)

	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestNestedFile(t *testing.T) {
	body := `Hello,
World!`
	res := rsvp.Response{Body: body}
	req := httptest.NewRequest("GET", "/files/file.txt", nil)
	rec := httptest.NewRecorder()
	cfg := rsvp.DefaultConfig()
	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", body, s)
}

func TestPutWithOkResponse(t *testing.T) {
	res := rsvp.Ok()
	req := httptest.NewRequest("PUT", "/files/file.md", bytes.NewBufferString("Some submission"))
	rec := httptest.NewRecorder()

	cfg := rsvp.DefaultConfig()
	err := res.Write(rec, req, cfg)
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestRequestUnknownFormat(t *testing.T) {
	body := "Hello <input> World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "unknown/unknown")
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
	assert.Eq(t, "Status code", 406, statusCode)

	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", "", s)
}

func TestComplexDataStructuresAreJsonByDefault(t *testing.T) {
	body := []string{"I", "am", "livid"}
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

	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))

	s := rec.Body.String()
	assert.Eq(t, "body contents", `["I","am","livid"]`+"\n", s)
}

func TestFirefoxAcceptHeader(t *testing.T) {
	body := "Hello <input> World!"
	res := rsvp.Response{Body: body, TemplateName: "tm"}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
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

func TestSeeOtherBlank(t *testing.T) {
	res := rsvp.SeeOther("/", nil)
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusSeeOther, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "", body)
}

func TestFoundBlank(t *testing.T) {
	res := rsvp.Found("/", nil)
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", http.StatusFound, statusCode)
	assert.Eq(t, "Content type", "", resp.Header.Get("Content-Type"))
	assert.Eq(t, "Location", "/", resp.Header.Get("Location"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "", body)
}

func TestRequestJsonEmptyString(t *testing.T) {
	res := rsvp.Response{Body: ""}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", `""`+"\n", body)
}

func TestRequestJsonNull(t *testing.T) {
	res := rsvp.Response{Body: nil}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "null\n", body)
}

func TestRespondJsonEmptyString(t *testing.T) {
	res := rsvp.Response{Body: ""}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", `""`+"\n", body)
}

func TestRespondJsonNull(t *testing.T) {
	res := rsvp.Response{Body: nil}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "null\n", body)
}

func TestRequestXmlEmptyString(t *testing.T) {
	res := rsvp.Response{Body: ""}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/xml")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/xml", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "<string></string>", body)
}

func TestRequestXmlNull(t *testing.T) {
	res := rsvp.Response{Body: nil}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/xml")
	rec := httptest.NewRecorder()

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/xml", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "", body)
}

func TestRespondXmlEmptyString(t *testing.T) {
	res := rsvp.Response{Body: ""}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/xml")

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/xml", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "<string></string>", body)
}

func TestRespondXmlNull(t *testing.T) {
	res := rsvp.Response{Body: nil}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/xml")

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 200, statusCode)
	assert.Eq(t, "Content type", "application/xml", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "", body)
}

func TestRequestForXmlButServingJson(t *testing.T) {
	res := rsvp.Response{Body: nil}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/xml")
	rec := httptest.NewRecorder()
	rec.Header().Set("Content-Type", "application/json")

	err := res.Write(rec, req, rsvp.DefaultConfig())
	assert.FatalErr(t, "Write response", err)

	resp := rec.Result()
	statusCode := resp.StatusCode
	assert.Eq(t, "Status code", 406, statusCode)
	assert.Eq(t, "Content type", "application/json", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Eq(t, "body contents", "", body)
}
