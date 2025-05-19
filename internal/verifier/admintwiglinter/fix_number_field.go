package admintwiglinter

import (
	"strings"

	"github.com/shyim/go-version"

	"github.com/allincart-org/allincart-cli/internal/html"
)

type NumberFieldFixer struct{}

func init() {
	AddFixer(NumberFieldFixer{})
}

func (n NumberFieldFixer) Check(nodes []html.Node) []CheckError {
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-number-field" {
			errs = append(errs, CheckError{
				Message:    "sw-number-field is removed, use mt-number-field instead. Please review conversion for props, events and label slot.",
				Severity:   "error",
				Identifier: "sw-number-field",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (n NumberFieldFixer) Supports(v *version.Version) bool {
	return allincart67Constraint.Check(v)
}

func (n NumberFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-number-field" {
			node.Tag = "mt-number-field"
			var newAttrs html.NodeList

			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case ColonValueAttr:
						newAttrs = append(newAttrs, html.Attribute{
							Key:   ":model-value",
							Value: attr.Value,
						})
					case VModelValueAttr:
						attr.Key = VModelAttr
						newAttrs = append(newAttrs, attr)
					case "@update:value":
						newAttrs = append(newAttrs, html.Attribute{
							Key:   "@change",
							Value: attr.Value,
						})
					default:
						newAttrs = append(newAttrs, attr)
					}
				} else {
					// If it's not an html.Attribute (e.g., TwigIfNode), preserve it as is
					newAttrs = append(newAttrs, attrNode)
				}
			}
			node.Attributes = newAttrs

			var label string
			var remainingChildren html.NodeList
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == TemplateTag {
					for _, a := range elem.Attributes {
						if attr, ok := a.(html.Attribute); ok {
							if attr.Key == LabelSlotAttr {
								var sb strings.Builder
								for _, inner := range elem.Children {
									sb.WriteString(strings.TrimSpace(inner.Dump(0)))
								}
								label = sb.String()
								goto SkipChild
							}
						}
					}
				}
				remainingChildren = append(remainingChildren, child)
			SkipChild:
			}
			node.Children = remainingChildren
			if label != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "label",
					Value: label,
				})
			}
		}
	})
	return nil
}
