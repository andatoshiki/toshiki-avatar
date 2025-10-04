package errors

import "fmt"

// Error types for CLI and HTTP
var (
	ErrNoAvatars = fmt.Errorf("no avatars found in directory or list file")
	ErrInvalidImageType = fmt.Errorf("invalid image type requested")
	ErrAvatarNotFound = fmt.Errorf("avatar not found")
	ErrPathTraversal = fmt.Errorf("avatar path is outside allowed directory")
	ErrFileUnreadable = fmt.Errorf("avatar file is not readable")
	ErrFlagConflict = fmt.Errorf("cannot use both -d and -l flags together or both empty")
)
