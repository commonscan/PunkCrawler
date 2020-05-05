package coolCrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
	"golang.org/x/net/html/charset"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var (
	titleRegex, _ = regexp.Compile("(?is)<title>(.*?)</title>")
	IgnoredExt    = []string{".3ds", ".3g2", ".3gp", ".7z", ".DS_Store", ".a", ".aac", ".adp", ".ai", ".aif", ".aiff", ".apk", ".ar", ".asf", ".au", ".avi", ".bak", ".bin", ".bk", ".bmp", ".btif", ".bz2", ".cab", ".caf", ".cgm", ".cmx", ".cpio", ".cr2", ".dat", ".deb", ".djvu", ".dll", ".dmg", ".dmp", ".dng", ".doc", ".docx", ".dot", ".dotx", ".dra", ".dsk", ".dts", ".dtshd", ".dvb", ".dwg", ".dxf", ".ear", ".ecelp4800", ".ecelp7470", ".ecelp9600", ".egg", ".eol", ".eot", ".epub", ".exe", ".f4v", ".fbs", ".fh", ".fla", ".flac", ".fli", ".flv", ".fpx", ".fst", ".fvt", ".g3", ".gif", ".gz", ".h261", ".h263", ".h264", ".ico", ".ief", ".image", ".img", ".ipa", ".iso", ".jar", ".jpeg", ".jpg", ".jpgv", ".jpm", ".jxr", ".ktx", ".lvp", ".lz", ".lzma", ".lzo", ".m3u", ".m4a", ".m4v", ".mar", ".mdi", ".mid", ".mj2", ".mka", ".mkv", ".mmr", ".mng", ".mov", ".movie", ".mp3", ".mp4", ".mp4a", ".mpeg", ".mpg", ".mpga", ".mxu", ".nef", ".npx", ".o", ".oga", ".ogg", ".ogv", ".otf", ".pbm", ".pcx", ".pdf", ".pea", ".pgm", ".pic", ".png", ".pnm", ".ppm", ".pps", ".ppt", ".pptx", ".ps", ".psd", ".pya", ".pyc", ".pyo", ".pyv", ".qt", ".rar", ".ras", ".raw", ".rgb", ".rip", ".rlc", ".rz", ".s3m", ".s7z", ".scm", ".scpt", ".sgi", ".shar", ".sil", ".smv", ".so", ".sub", ".swf", ".tar", ".tbz2", ".tga", ".tgz", ".tif", ".tiff", ".tlz", ".ts", ".ttf", ".uvh", ".uvi", ".uvm", ".uvp", ".uvs", ".uvu", ".viv", ".vob", ".war", ".wav", ".wax", ".wbmp", ".wdp", ".weba", ".webm", ".webp", ".whl", ".wm", ".wma", ".wmv", ".wmx", ".woff", ".woff2", ".wvx", ".xbm", ".xif", ".xls", ".xlsx", ".xlt", ".xm", ".xpi", ".xpm", ".xwd", ".xz", ".z", ".zip", ".zipx"}
)

func getTld(uri string) string {
	parsedTld := extract.Extract(uri)
	var rootDomain string
	if len(parsedTld.Tld) > 0 {
		rootDomain = fmt.Sprintf("%s.%s", parsedTld.Root, parsedTld.Tld)
	} else {
		rootDomain = fmt.Sprintf(parsedTld.Root)
	}
	return rootDomain
}
func getTitle(text string) string {
	matchResult := titleRegex.FindStringSubmatch(text)
	if len(matchResult) > 0 {
		return matchResult[1]
	} else {
		return ""
	}
}
func getRemoteIP(host string) ([]string, error) {
	if net.ParseIP(host).To4() != nil {
		return []string{net.ParseIP(host).To4().String()}, nil
	}
	hosts, err := net.LookupHost(host)
	if err != nil {
		return []string{}, err
	} else {
		return hosts, err
	}
}

func joinURLs(baseURL, hyperlink string) string {
	parse, err := url.Parse(hyperlink)
	if err != nil {
		return ""
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	nextURLToCrawl := base.ResolveReference(parse)
	return nextURLToCrawl.String()
}

// 输入一个HTML，输出一个Link数组
func getLinks(text string, baseUrl string) (links []string) {
	rootNode, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return
	}
	dupeSet := map[string]struct{}{}
	rootNode.Find("a").Each(func(i int, s *goquery.Selection) {
		uri, exists := s.Attr("href")
		if !exists {
			return
		} else {
			for _, i := range IgnoredExt {
				if strings.HasSuffix(uri, i) {
					return
				}
			}
			if strings.HasPrefix(uri, "javascript:") {
				return
			}
			if !strings.HasPrefix(uri, "http") {
				uri = joinURLs(baseUrl, uri)
			}
			if _, ok := dupeSet[uri]; !ok {
				links = append(links, uri)
				dupeSet[uri] = struct{}{}
			}
		}
	})
	return links
}

func FixEncoding(input string, contentType string) (string, error) {
	// e encoding.Encoding, name string, certain bool
	_, encodingName, _ := charset.DetermineEncoding([]byte(input), contentType)
	return iconv.ConvertString(input, encodingName, "utf8")
}
