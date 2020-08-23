package common

import (
	"github.com/imfht/req"
	"github.com/rs/zerolog/log"
	"regexp"
)

var headRegexp = regexp.MustCompile(`(?i)charset=(.*?)$`)

func getHeaderEncoding(resp *req.Resp) string {
	ct := resp.Response().Header.Get("Content-Type")
	charSet := headRegexp.FindStringSubmatch(ct)
	if len(charSet) > 0 {
		return charSet[1]
	}
	return ""
}

// todo: add html encoding here.
func GetHtmlEncoding(resp *req.Resp) string {
	return ""
}
func GetResponseEncodingRemoved(resp *req.Resp) string {
	f := []func(resp2 *req.Resp) string{getHeaderEncoding}
	for _, v := range f {
		charset := v(resp)
		if len(charset) > 0 {
			log.Trace().Msgf("html encoding is %s .", charset)
			return charset
		}
	}
	return ""
}
