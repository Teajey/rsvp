package content

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
)

type proposal struct {
	superType string
	subType   string
	Weight    float32
	params    map[string][]string
}

func (p *proposal) MediaType() string {
	return fmt.Sprintf("%s/%s", p.superType, p.subType)
}

func proposalEqual(a, b proposal) bool {
	return a.subType == b.subType && a.superType == b.superType && a.Weight == b.Weight
}

func proposalCmp(a, b proposal) int {
	if proposalEqual(a, b) {
		return 0
	}

	if a.Weight < b.Weight {
		return 1
	}

	if a.superType == "*" && b.superType != "*" {
		return 1
	}

	if b.superType == "*" && a.superType != "*" {
		return -1
	}

	if a.subType == "*" && b.subType != "*" {
		return 1
	}

	if b.subType == "*" && a.subType != "*" {
		return -1
	}

	return -1
}

var ErrorProposalEmpty = errors.New("Empty media-type")
var ErrorProposalEmptySuper = errors.New("Empty media supertype, i.e. <supertype>/<subtype>")
var ErrorProposalEmptySub = errors.New("Empty media subtype, i.e. <supertype>/<subtype>")
var ErrorProposalWildSuper = errors.New("Supertype may not be wild on it's own, i.e. */<subtype>")

type ErrorProposalBadWeightFloat struct {
	Err error
}

func (e ErrorProposalBadWeightFloat) Error() string {
	return fmt.Sprintf("Couldn't parse weight as float32: %s", e.Err)
}

func (e ErrorProposalBadWeightFloat) Unwrap() error {
	return e.Err
}

func ParseProposal(s string) (proposal, error) {
	nextElem, stopElem := iter.Pull(splitAndTrimSpace(s, ";"))
	defer stopElem()
	mediaType, _ := nextElem()
	if mediaType == "" {
		return proposal{}, ErrorProposalEmpty
	}

	nextType, stopType := iter.Pull(splitAndTrimSpace(mediaType, "/"))
	defer stopType()
	superType, _ := nextType()
	if superType == "" {
		return proposal{}, ErrorProposalEmptySuper
	}
	subType, ok := nextType()
	if !ok || subType == "" {
		return proposal{}, ErrorProposalEmptySub
	}
	if superType == "*" && subType != "*" {
		return proposal{}, ErrorProposalWildSuper
	}

	proposal := proposal{
		superType: superType,
		subType:   subType,
		Weight:    1.,
	}

	paramsStr, ok := nextElem()
	if !ok || paramsStr == "" {
		return proposal, nil
	}

	proposal.params = collectAppend(parseParameters(paramsStr))
	weightings, ok := proposal.params["q"]
	if !ok {
		return proposal, nil
	}
	if len(weightings) < 1 {
		return proposal, nil
	}
	weighting := weightings[0]

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
		proposal.Weight = f
	case val < 0:
		f := float32(0.)
		proposal.Weight = f
	default:
		f := float32(val)
		proposal.Weight = f
	}
	return proposal, nil
}
