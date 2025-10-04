package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// IsReadable checks if a file is readable by attempting to open it.
func IsReadable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

// IsPathWithinBase checks if the given path is within the base directory.
func IsPathWithinBase(base, target string) bool {
	absBase, err1 := filepath.Abs(base)
	absTarget, err2 := filepath.Abs(target)
	if err1 != nil || err2 != nil {
		return false
	}
	return len(absTarget) >= len(absBase) && absTarget[:len(absBase)] == absBase
}

// NormalizePath ensures path separators are consistent across OSes.
func NormalizePath(path string) string {
	if runtime.GOOS == "windows" {
		return filepath.ToSlash(path)
	}
	return path
}
