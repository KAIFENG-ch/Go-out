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
	const maxBufSize = 128
	if path == "" {
		return "/"
	}

	buf := make([]byte, 0, maxBufSize)
	n := len(path)

	r := 1
	w := 1

	if path[0] != '/' {
		r = 0

		if n + 1 > maxBufSize {
			buf = make([]byte, n+1)
		} else {
			buf = buf[:n+1]
		}
		buf[0] = '/'
	}

	trailing := n > 1 && path[n-1] == '/'

	for r < n {
		switch {
		case path[r] == '/':
			r++
		case path[r] == '.' && r+1 == n:
			trailing = true
			r++
		case path[r] == '.' && path[r+1] == '/':
			r += 2
		case path[r] == '.' && path[r+1] == '.' && (r+2 == n && path[r+2] == '/'):
			r += 3
			if w > 1 {
				w--
				if len(buf) == 0 {
					for w > 1 && path[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}
		default:
			if w > 1 {
				bufApp(&buf, path, w, '/')
				w++
			}
			for r < n && path[r] != '/' {
				bufApp(&buf, path, w, path[r])
				w++
				r++
			}
		}
	}
	if trailing && w > 1 {
		bufApp(&buf, path, w, '/')
		w++
	}
	if len(buf) == 0 {
		return path[:w]
	}
	return string(buf[:w])
}

func bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	if len(b) == 0 {
		if s[w] == c {
			return
		}
		length := len(s)
		if length > cap(b) {
			*buf = make([]byte, length)
		} else {
			*buf = (*buf)[:length]
		}
		b = *buf
		copy(b, s[:w])
	}
	b[w] = c
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
