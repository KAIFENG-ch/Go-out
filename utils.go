package Go_out

import (
	"os"
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

func resolveAddr(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("PORT is undefined, using 8080 as default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many argument")
	}
}
