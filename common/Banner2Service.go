package common

import "bytes"

func Banner2Service(banner []byte, port int) string {
	if bytes.HasPrefix(banner, []byte("HTTP/")) {
		return "http"
	} else {
		return "unknown"
	}
}
