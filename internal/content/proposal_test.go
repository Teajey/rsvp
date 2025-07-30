package content_test

import (
	"testing"

	"github.com/Teajey/rsvp/internal/assert"
	"github.com/Teajey/rsvp/internal/content"
)

func TestParseProposalHtml(t *testing.T) {
	proposal, err := content.ParseProposal("text/html")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/html")
	assert.Eq(t, "Weight value", proposal.Weight, 1.)
}

func TestParseProposalHtmlWeighted(t *testing.T) {
	proposal, err := content.ParseProposal("text/html;q=0.8")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/html")
	assert.Eq(t, "Weight value", proposal.Weight, 0.8)
}

func TestParseProposalHtmlWeirdSpacing(t *testing.T) {
	proposal, err := content.ParseProposal(" text / html ; q = 0.8 ")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/html")
	assert.Eq(t, "Weight value", proposal.Weight, 0.8)
}

func TestParseProposalEmpty(t *testing.T) {
	_, err := content.ParseProposal("")
	assert.FatalErrIs(t, "Parsing proposal", err, content.ErrorProposalEmpty)
}

func TestParseProposalEmptySuper(t *testing.T) {
	_, err := content.ParseProposal("/html")
	assert.FatalErrIs(t, "Parsing proposal", err, content.ErrorProposalEmptySuper)
}

func TestParseProposalEmptySub(t *testing.T) {
	_, err := content.ParseProposal("text/")
	assert.FatalErrIs(t, "Parsing proposal", err, content.ErrorProposalEmptySub)
}

func TestParseProposalBadWeightPrefix(t *testing.T) {
	proposal, err := content.ParseProposal("text/html;w=0.8")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/html")
	assert.Eq(t, "Weight value", proposal.Weight, 1.)
}

func TestParseProposalBadWeightFloat(t *testing.T) {
	_, err := content.ParseProposal("text/html;q=*")
	var badWeightFloat content.ErrorProposalBadWeightFloat
	assert.FatalErrAs(t, "Parsing proposal", err, &badWeightFloat)
}

func TestParseProposalDouble(t *testing.T) {
	_, err := content.ParseProposal("text/html;q=0.8,text/plain;q=0.7")
	var badWeightFloat content.ErrorProposalBadWeightFloat
	assert.FatalErrAs(t, "Parsing proposal", err, &badWeightFloat)
}

func TestParseProposalWild(t *testing.T) {
	proposal, err := content.ParseProposal("*/*;q=0.8")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "*/*")
	assert.Eq(t, "Weight value", proposal.Weight, 0.8)
}

func TestParseProposalWildSuper(t *testing.T) {
	_, err := content.ParseProposal("*/html;q=0.8")
	assert.FatalErrIs(t, "Parsing proposal", err, content.ErrorProposalWildSuper)
}

func TestParseProposalWildSub(t *testing.T) {
	proposal, err := content.ParseProposal("text/*;q=0.8")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/*")
	assert.Eq(t, "Weight value", proposal.Weight, 0.8)
}

func TestParseProposalOverWeighted(t *testing.T) {
	proposal, err := content.ParseProposal("text/html;q=1.1")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/html")
	assert.Eq(t, "Weight value", proposal.Weight, 1)
}

func TestParseProposalUnderWeighted(t *testing.T) {
	proposal, err := content.ParseProposal("text/html;q=-0.1")
	assert.FatalErr(t, "Parsing proposal", err)
	assert.Eq(t, "MediaType value", proposal.MediaType(), "text/html")
	assert.Eq(t, "Weight value", proposal.Weight, 0)
}
