package admintwiglinter

import (
	"strings"

	"github.com/shyim/go-version"

	"github.com/allincart-org/allincart-cli/internal/html"
)

type ColorpickerFixer struct{}

func init() {
	AddFixer(ColorpickerFixer{})
}

func (c ColorpickerFixer) Check(nodes []html.Node) []CheckError {
	var checkErrors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-colorpicker" {
			checkErrors = append(checkErrors, CheckError{
				Message:    "sw-colorpicker is removed, use mt-colorpicker instead. Please review conversion for label property.",
				Severity:   "error",
				Identifier: "sw-colorpicker",
				Line:       node.Line,
			})
		}
	})
	return checkErrors
}

func (c ColorpickerFixer) Supports(v *version.Version) bool {
	return allincart67Constraint.Check(v)
}

func (c ColorpickerFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-colorpicker" {
			node.Tag = "mt-colorpicker"

			var newAttrs html.NodeList
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case ColonValueAttr:
						attr.Key = ":model-value"
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

			// Process label slot: extract inner text and add as label attribute.
			label := ""
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok {
					if elem.Tag == TemplateTag {
						for _, a := range elem.Attributes {
							if attr, ok := a.(html.Attribute); ok {
								if attr.Key == LabelSlotAttr {
									var sb strings.Builder
									for _, inner := range elem.Children {
										sb.WriteString(strings.TrimSpace(inner.Dump(0)))
									}
									label = sb.String()
								}
							}
						}
					}
				}
			}
			node.Children = html.NodeList{}
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
