package common

import (
	"github.com/imfht/req"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
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

var disabledExtentions = []string{".3ds", ".3g2", ".3gp", ".7z", ".DS_Store", ".a", ".aac", ".adp", ".ai", ".aif", ".aiff", ".apk", ".ar", ".asf", ".au", ".avi", ".bak", ".bin", ".bk", ".bmp", ".btif", ".bz2", ".cab", ".caf", ".cgm", ".cmx", ".cpio", ".cr2", ".dat", ".deb", ".djvu", ".dll", ".dmg", ".dmp", ".dng", ".doc", ".docx", ".dot", ".dotx", ".dra", ".dsk", ".dts", ".dtshd", ".dvb", ".dwg", ".dxf", ".ear", ".ecelp4800", ".ecelp7470", ".ecelp9600", ".egg", ".eol", ".eot", ".epub", ".exe", ".f4v", ".fbs", ".fh", ".fla", ".flac", ".fli", ".flv", ".fpx", ".fst", ".fvt", ".g3", ".gif", ".gz", ".h261", ".h263", ".h264", ".ico", ".ief", ".image", ".img", ".ipa", ".iso", ".jar", ".jpeg", ".jpg", ".jpgv", ".jpm", ".jxr", ".ktx", ".lvp", ".lz", ".lzma", ".lzo", ".m3u", ".m4a", ".m4v", ".mar", ".mdi", ".mid", ".mj2", ".mka", ".mkv", ".mmr", ".mng", ".mov", ".movie", ".mp3", ".mp4", ".mp4a", ".mpeg", ".mpg", ".mpga", ".mxu", ".nef", ".npx", ".o", ".oga", ".ogg", ".ogv", ".otf", ".pbm", ".pcx", ".pdf", ".pea", ".pgm", ".pic", ".png", ".pnm", ".ppm", ".pps", ".ppt", ".pptx", ".ps", ".psd", ".pya", ".pyc", ".pyo", ".pyv", ".qt", ".rar", ".ras", ".raw", ".rgb", ".rip", ".rlc", ".rz", ".s3m", ".s7z", ".scm", ".scpt", ".sgi", ".shar", ".sil", ".smv", ".so", ".sub", ".swf", ".tar", ".tbz2", ".tga", ".tgz", ".tif", ".tiff", ".tlz", ".ts", ".ttf", ".uvh", ".uvi", ".uvm", ".uvp", ".uvs", ".uvu", ".viv", ".vob", ".war", ".wav", ".wax", ".wbmp", ".wdp", ".weba", ".webm", ".webp", ".whl", ".wm", ".wma", ".wmv", ".wmx", ".woff", ".woff2", ".wvx", ".xbm", ".xif", ".xls", ".xlsx", ".xlt", ".xm", ".xpi", ".xpm", ".xwd", ".xz", ".z", ".zip", ".zipx"}

func UrlHasDisableExtension(url string) bool {
	for _, item := range disabledExtentions {
		if strings.HasSuffix(strings.ToLower(url), item) {
			return true
		}
	}
	return false
}
