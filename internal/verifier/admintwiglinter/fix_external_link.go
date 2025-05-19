package admintwiglinter

import (
	"github.com/shyim/go-version"

	"github.com/allincart-org/allincart-cli/internal/html"
)

type ExternalLinkFixer struct{}

func init() {
	AddFixer(ExternalLinkFixer{})
}

func (e ExternalLinkFixer) Check(nodes []html.Node) []CheckError {
	var checkErrors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-external-link" {
			checkErrors = append(checkErrors, CheckError{
				Message:    "sw-external-link is removed, use mt-external-link instead and remove the icon property.",
				Severity:   "error",
				Identifier: "sw-external-link",
				Line:       node.Line,
			})
		}
	})
	return checkErrors
}

func (e ExternalLinkFixer) Supports(v *version.Version) bool {
	return allincart67Constraint.Check(v)
}

func (e ExternalLinkFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-external-link" {
			node.Tag = "mt-external-link"
			var newAttrs html.NodeList
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					if attr.Key == "icon" {
						continue
					}
					newAttrs = append(newAttrs, attr)
				} else {
					// If it's not an html.Attribute (e.g., TwigIfNode), preserve it as is
					newAttrs = append(newAttrs, attrNode)
				}
			}
			node.Attributes = newAttrs
		}
	})
	return nil
}
