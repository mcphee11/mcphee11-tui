package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents from %s to %s: %w", src, dst, err)
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info %s: %w", src, err)
	}
	// Preserve file permissions
	if err = os.Chmod(dst, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set permissions for destination file %s: %w", dst, err)
	}

	return nil
}

// CopyDir recursively copies a directory from src to dst.
func CopyDir(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source directory info %s: %w", src, err)
	}

	// Create the destination directory with the same permissions
	if err := os.MkdirAll(dst, sourceInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory %s: %w", src, err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(sourcePath, destPath); err != nil {
				return err // Propagate the error up
			}
		} else {
			if err := CopyFile(sourcePath, destPath); err != nil {
				return err // Propagate the error up
			}
		}
	}
	return nil
}
