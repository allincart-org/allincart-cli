package verifier

import (
	"context"
	"fmt"
	"strings"

	"github.com/allincart-org/allincart-cli/extension"
)

type ToolList []Tool

var availableTools = ToolList{}

func AddTool(tool Tool) {
	availableTools = append(availableTools, tool)
}

func GetTools() ToolList {
	return availableTools
}

type ToolConfig struct {
	// Path to the tool directory
	ToolDirectory string

	InputWasDirectory bool

	// The minimum version of Allincart that is supported
	MinAllincartVersion string
	// The maximum version of Allincart that is supported
	MaxAllincartVersion string
	// The version of Allincart that is checked against
	CheckAgainst string
	// The root directory of the extension/project
	RootDir string
	// Contains a list of directories that are considered as source code
	SourceDirectories []string
	// Contains a list of identifiers that are ignored
	ValidationIgnores []ToolConfigIgnore
	// Contains a list of directories that are considered as admin code
	AdminDirectories []string
	// Contains a list of directories that are considered as storefront code
	StorefrontDirectories []string

	Extension extension.Extension
}

type ToolConfigIgnore struct {
	Identifier string
	Path       string
	Message    string
}

type Tool interface {
	Name() string
	Check(ctx context.Context, check *Check, config ToolConfig) error
	Fix(ctx context.Context, config ToolConfig) error
	Format(ctx context.Context, config ToolConfig, dryRun bool) error
}

func (tl ToolList) Only(only string) (ToolList, error) {
	if only == "" {
		return tl, nil
	}

	var filteredTools []Tool
	requestedTools := strings.Split(only, ",")

	for _, requestedTool := range requestedTools {
		requestedTool = strings.TrimSpace(requestedTool)
		found := false

		for _, t := range tl {
			if t.Name() == requestedTool {
				filteredTools = append(filteredTools, t)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("tool with name %q not found, possible tools: %s", requestedTool, tl.PossibleString())
		}
	}

	return filteredTools, nil
}

func (tl ToolList) PossibleString() string {
	var possibleTools []string
	for _, t := range tl {
		possibleTools = append(possibleTools, t.Name())
	}

	return strings.Join(possibleTools, ",")
}
