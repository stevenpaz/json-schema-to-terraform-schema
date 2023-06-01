package internal

import (
	"fmt"
	"os"
)

// WriteFileBytes writes bytes to a file path. If the file already exists, it is
// truncated.
func WriteFileBytes(path string, bytes []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
