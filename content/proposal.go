package content

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
)

type Proposal struct {
	superType string
	subType   string
	Weight    *float32
}

var ErrorProposalEmpty = errors.New("Empty media-type")
var ErrorProposalEmptySuper = errors.New("Empty media supertype, i.e. <supertype>/<subtype>")
var ErrorProposalEmptySub = errors.New("Empty media subtype, i.e. <supertype>/<subtype>")

type ErrorProposalBadWeightPrefix struct {
	Weighting string
}

func (e ErrorProposalBadWeightPrefix) Error() string {
	return fmt.Sprintf("Invalid weighting '%s'. Expected 'q=' followed by a number ranging from 0 to 1", e.Weighting)
}

type ErrorProposalBadWeightFloat struct {
	Err error
}

func (e ErrorProposalBadWeightFloat) Error() string {
	return fmt.Sprintf("Couldn't parse weight as float32: %s", e.Err)
}

func (e ErrorProposalBadWeightFloat) Unwrap() error {
	return e.Err
}

func ParseProposal(s string) (Proposal, error) {
	nextElem, stopElem := iter.Pull(splitAndTrimSpace(s, ";"))
	defer stopElem()
	mediaType, _ := nextElem()
	if mediaType == "" {
		return Proposal{}, ErrorProposalEmpty
	}

	nextType, stopType := iter.Pull(splitAndTrimSpace(mediaType, "/"))
	defer stopType()
	superType, _ := nextType()
	if superType == "" {
		return Proposal{}, ErrorProposalEmptySuper
	}
	subType, ok := nextType()
	if !ok || subType == "" {
		return Proposal{}, ErrorProposalEmptySub
	}

	proposal := Proposal{
		superType: superType,
		subType:   subType,
		Weight:    nil,
	}

	weighting, ok := nextElem()
	if !ok || weighting == "" {
		return proposal, nil
	}

	nextPair, stopPair := iter.Pull(splitAndTrimSpace(weighting, "="))
	defer stopPair()

	key, ok := nextPair()
	if !ok {
		return proposal, nil
	}
	if key != "q" {
		return proposal, ErrorProposalBadWeightPrefix{weighting}
	}
	weighting, ok = nextPair()
	if !ok {
		return proposal, nil
	}

	if len(weighting) > 4 {
		return proposal, ErrorProposalBadWeightFloat{errors.New("Weight is unlikely to be valid because it contains more than 4 characters")}
	}

	val, err := strconv.ParseFloat(weighting, 32)
	if err != nil {
		return proposal, ErrorProposalBadWeightFloat{err}
	}

	switch {
	case val > 1:
		f := float32(1.)
		proposal.Weight = &f
	case val < 0:
		f := float32(0.)
		proposal.Weight = &f
	default:
		f := float32(val)
		proposal.Weight = &f
	}
	return proposal, nil
}

func (p *Proposal) MediaType() string {
	return fmt.Sprintf("%s/%s", p.superType, p.subType)
}
