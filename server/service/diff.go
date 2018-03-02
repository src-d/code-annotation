package service

import "github.com/pmezard/go-difflib/difflib"

// Diff service generates diff for files
type Diff struct {
	context int
}

// NewDiff creates Diff service
func NewDiff() *Diff {
	return &Diff{context: 6} // keep it hard coded for now
}

// Generate return unified diff string for 2 files
func (d *Diff) Generate(nameA, nameB, contentA, contentB string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(contentA),
		B:        difflib.SplitLines(contentB),
		FromFile: nameA,
		ToFile:   nameB,
		Context:  d.context,
	}

	return difflib.GetUnifiedDiffString(diff)
}
