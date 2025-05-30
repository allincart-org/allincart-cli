package admintwiglinter

import (
	"github.com/shyim/go-version"

	"github.com/allincart-org/allincart-cli/internal/html"
)

type SkeletonBarFixer struct{}

func init() {
	AddFixer(SkeletonBarFixer{})
}

func (s SkeletonBarFixer) Check(nodes []html.Node) []CheckError {
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-skeleton-bar" {
			errors = append(errors, CheckError{
				Message:    "sw-skeleton-bar is removed, use mt-skeleton-bar instead.",
				Severity:   "error",
				Identifier: "sw-skeleton-bar",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (s SkeletonBarFixer) Supports(v *version.Version) bool {
	return allincart67Constraint.Check(v)
}

func (s SkeletonBarFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-skeleton-bar" {
			node.Tag = "mt-skeleton-bar"
		}
	})
	return nil
}
