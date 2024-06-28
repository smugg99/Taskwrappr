// scripter.go
package taskwrappr

import (
	"os"
)

type Script struct {
	Path           string
	Content        string
	CleanedContent string
}

func NewScript(filePath string) (*Script, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Script{
		Path:    filePath,
		Content: string(content),
	}, nil
}
