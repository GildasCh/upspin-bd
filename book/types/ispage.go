package types

import "strings"

var imageExtensions = []string{
	"jpg", "jpeg",
	"png",
	"gif",
	"bmp",
}

func IsImage(name string, isDir bool) bool {
	if isDir {
		return false
	}

	name = strings.ToLower(name)

	for _, ext := range imageExtensions {
		if strings.HasSuffix(name, "."+ext) {
			return true
		}
	}
	return false
}
