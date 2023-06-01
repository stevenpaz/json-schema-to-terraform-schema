package internal

import (
	"fmt"

	"mvdan.cc/gofumpt/format"
)

// FormatCode formats the given Go code using gofumpt.
func FormatGoCode(data []byte) ([]byte, error) {
	formattedCode, err := format.Source(data, format.Options{
		LangVersion: "1.20",
		ExtraRules:  false,
		ModulePath:  "",
	})
	if err != nil {
		return nil, fmt.Errorf("error formatting code: %w", err)
	}

	return formattedCode, nil
}
