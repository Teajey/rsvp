package rsvp_test

import (
	"encoding/csv"
	html "html/template"
	"slices"
	"testing"
	text "text/template"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

func TestNilResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data(nil)
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestBlankResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Blank()
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestEmptyStringResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data("")
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStructResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data(struct{}{})
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestSliceResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data([]struct{}{})
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestMapResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data(map[string]any{})
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStringResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data("Hello")
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestBytesResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data([]byte("Hello"))
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeBytes,
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

type myCsv struct{}

func (myCsv) MarshalCsv(w *csv.Writer) error {
	panic("not implemented")
}

func TestCsvResponse(t *testing.T) {
	cfg := rsvp.Config{}
	resp := rsvp.Data(myCsv{})
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeCsv,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStructResponseWithHtmlTemplate(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
	}
	resp := rsvp.Body{Data: struct{}{}, TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStructWithoutTemplateNameResponseWithHtmlTemplate(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
	}
	resp := rsvp.Body{Data: struct{}{}, TemplateName: ""}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestCsvResponseWithHtmlTemplate(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
	}
	resp := rsvp.Body{Data: myCsv{}, TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeCsv,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestCsvResponseWithTextTemplate(t *testing.T) {
	cfg := rsvp.Config{
		TextTemplate: text.New(""),
	}
	resp := rsvp.Body{Data: myCsv{}, TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeCsv,
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStructResponseWithTextTemplate(t *testing.T) {
	cfg := rsvp.Config{
		TextTemplate: text.New(""),
	}
	resp := rsvp.Body{Data: struct{}{}, TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStructWithoutTemplateNameResponseWithTextTemplate(t *testing.T) {
	cfg := rsvp.Config{
		TextTemplate: text.New(""),
	}
	resp := rsvp.Body{Data: struct{}{}, TemplateName: ""}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStructResponseWithTextAndHtmlTemplates(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
		TextTemplate: text.New(""),
	}
	resp := rsvp.Body{Data: struct{}{}, TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStringResponseWithTextAndHtmlTemplates(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
		TextTemplate: text.New(""),
	}
	resp := rsvp.Body{Data: "Hello", TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestCsvResponseWithTextAndHtmlTemplates(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
		TextTemplate: text.New(""),
	}
	resp := rsvp.Body{Data: myCsv{}, TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeCsv,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestBytesResponseWithHtmlTemplate(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
	}
	resp := rsvp.Body{Data: []byte("Hello"), TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypeBytes,
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}

func TestStringResponseWithHtmlTemplate(t *testing.T) {
	cfg := rsvp.Config{
		HtmlTemplate: html.New(""),
	}
	resp := rsvp.Body{Data: "Hello", TemplateName: "tm"}
	actual := slices.Collect(resp.MediaTypes(cfg))

	expected := []string{
		rsvp.SupportedMediaTypePlaintext,
		rsvp.SupportedMediaTypeJson,
		rsvp.SupportedMediaTypeXml,
		rsvp.SupportedMediaTypeHtml,
		rsvp.SupportedMediaTypeGob,
	}

	assert.SlicesEq(t, "media types", expected, actual)
}
