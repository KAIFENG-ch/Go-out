package Go_out

import (
	"path"
)

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}
	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) == '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("the string is null")
	}
	return str[len(str)-1]
}

func cleanPath(path string) string {
	// todo
	return ""
}
