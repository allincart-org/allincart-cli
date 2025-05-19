package admintwiglinter

import (
	"github.com/shyim/go-version"

	"github.com/allincart-org/allincart-cli/internal/html"
)

type ProgressBarFixer struct{}

func init() {
	AddFixer(ProgressBarFixer{})
}

func (p ProgressBarFixer) Check(nodes []html.Node) []CheckError {
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-progress-bar" {
			errors = append(errors, CheckError{
				Message:    "sw-progress-bar is removed, use mt-progress-bar instead.",
				Severity:   "error",
				Identifier: "sw-progress-bar",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (p ProgressBarFixer) Supports(v *version.Version) bool {
	return allincart67Constraint.Check(v)
}

func (p ProgressBarFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-progress-bar" {
			node.Tag = "mt-progress-bar"
			var newAttrs html.NodeList

			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case ValueAttr:
						attr.Key = ModelValueAttr
						newAttrs = append(newAttrs, attr)
					case VModelValueAttr:
						attr.Key = VModelAttr
						newAttrs = append(newAttrs, attr)
					case UpdateValueAttr:
						attr.Key = UpdateModelValueAttr
						newAttrs = append(newAttrs, attr)
					default:
						newAttrs = append(newAttrs, attr)
					}
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
