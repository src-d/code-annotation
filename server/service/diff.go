package service

import (
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

// Diff service generates diff for files
type Diff struct {
	context int
}

// NewDiff creates Diff service
func NewDiff() *Diff {
	return &Diff{context: 6} // keep it hard coded for now
}

// DiffPreprocessorFunc type is function signature to preprocess diffs
type DiffPreprocessorFunc func(string) string

// Generate return unified diff string for 2 files
func (d *Diff) Generate(nameA, nameB, contentA, contentB string, preprocessors ...DiffPreprocessorFunc) (string, error) {
	for _, p := range preprocessors {
		contentA = p(contentA)
		contentB = p(contentB)
	}

	return d.generate(nameA, nameB, contentA, contentB)
}

func (d *Diff) generate(nameA, nameB, contentA, contentB string) (string, error) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(contentA),
		B:        difflib.SplitLines(contentB),
		FromFile: nameA,
		ToFile:   nameB,
		Context:  d.context,
	}

	return difflib.GetUnifiedDiffString(diff)
}

// ReplaceInvisible preprocessor function that replace invisible character with visible onces
func ReplaceInvisible(content string) string {
	content = strings.Replace(content, " ", "·", -1)
	content = strings.Replace(content, "\t", "→", -1)
	content = strings.Replace(content, "\r", "^M", -1)
	return strings.Replace(content, "\n", "↵\n", -1)
}
