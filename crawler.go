package coolCrawler

import (
	"bufio"
	"bytes"
	"coolCrawler/common"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/imfht/req"
	"github.com/joeguo/tldextract"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Fetcher struct {
	InputFileName          string `short:"i" long:"input-file" description:"输入文件名" default:"-"`
	OutputFileName         string `short:"o" long:"output-file" description:"输出文件名" default:"-"`
	ProcessNum             int    `short:"p" long:"process-num" description:"并发数" default:"100"`
	Timeout                int32  `short:"t" long:"timeout" description:"最大超时数(s)" default:"30"`
	Retries                int    `short:"r" long:"retries" description:"最大重试次数" default:"2"`
	WithTitle              bool   `long:"with-title" description:"是否输出Title"`
	WithHTML               bool   `long:"with-html" description:"是否输出HTML"`
	WithTld                bool   `long:"with-tld" description:"是否输出TLD"`
	WithIP                 bool   `long:"with-ip" description:"是否输出IP"`
	IconMode               bool   `long:"icon-mode" description:"输出Body的base64和hash"`
	WithHeaders            bool   `long:"with-headers" description:"是否输出Headers"`
	DefaultHTTPS           bool   `long:"default-https" description:"没有协议号的域名默认使用https"`
	UserAgent              string `long:"user-agent" description:"User-Agent" default:"Mozilla/5.0 (compatible;Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)"`
	WithCert               bool   `long:"with-cert" description:"是否输出HTTPS证书"`
	WithLinks              bool   `long:"with-links" description:"是否输出链接信息"`
	PreScan                bool   `long:"pre-scan" description:"探测前先端口扫描"`
	Ports                  string `long:"ports" description:"扫描的端口，用 ,分割" default:"80,8080,443"`
	OutPutTable            bool   `long:"table" description:"输出 table而不是json"`
	FilterBinaryExtensions bool   `long:"filter-binary" description:"是否过滤已知的二进制后缀URL"`
	Endpoint               string `long:"endp" description:"endpoint" default:"/"`
	NoLog                  bool   `long:"no-log" description:"不输出log信息"`
}

var (
	cache              = "/tmp/tld.cache"
	extract, _         = tldextract.New(cache, false)
	disabledExtentions = []string{".3ds", ".3g2", ".3gp", ".7z", ".DS_Store", ".a", ".aac", ".adp", ".ai", ".aif", ".aiff", ".apk", ".ar", ".asf", ".au", ".avi", ".bak", ".bin", ".bk", ".bmp", ".btif", ".bz2", ".cab", ".caf", ".cgm", ".cmx", ".cpio", ".cr2", ".dat", ".deb", ".djvu", ".dll", ".dmg", ".dmp", ".dng", ".doc", ".docx", ".dot", ".dotx", ".dra", ".dsk", ".dts", ".dtshd", ".dvb", ".dwg", ".dxf", ".ear", ".ecelp4800", ".ecelp7470", ".ecelp9600", ".egg", ".eol", ".eot", ".epub", ".exe", ".f4v", ".fbs", ".fh", ".fla", ".flac", ".fli", ".flv", ".fpx", ".fst", ".fvt", ".g3", ".gif", ".gz", ".h261", ".h263", ".h264", ".ico", ".ief", ".image", ".img", ".ipa", ".iso", ".jar", ".jpeg", ".jpg", ".jpgv", ".jpm", ".jxr", ".ktx", ".lvp", ".lz", ".lzma", ".lzo", ".m3u", ".m4a", ".m4v", ".mar", ".mdi", ".mid", ".mj2", ".mka", ".mkv", ".mmr", ".mng", ".mov", ".movie", ".mp3", ".mp4", ".mp4a", ".mpeg", ".mpg", ".mpga", ".mxu", ".nef", ".npx", ".o", ".oga", ".ogg", ".ogv", ".otf", ".pbm", ".pcx", ".pdf", ".pea", ".pgm", ".pic", ".png", ".pnm", ".ppm", ".pps", ".ppt", ".pptx", ".ps", ".psd", ".pya", ".pyc", ".pyo", ".pyv", ".qt", ".rar", ".ras", ".raw", ".rgb", ".rip", ".rlc", ".rz", ".s3m", ".s7z", ".scm", ".scpt", ".sgi", ".shar", ".sil", ".smv", ".so", ".sub", ".swf", ".tar", ".tbz2", ".tga", ".tgz", ".tif", ".tiff", ".tlz", ".ts", ".ttf", ".uvh", ".uvi", ".uvm", ".uvp", ".uvs", ".uvu", ".viv", ".vob", ".war", ".wav", ".wax", ".wbmp", ".wdp", ".weba", ".webm", ".webp", ".whl", ".wm", ".wma", ".wmv", ".wmx", ".woff", ".woff2", ".wvx", ".xbm", ".xif", ".xls", ".xlsx", ".xlt", ".xm", ".xpi", ".xpm", ".xwd", ".xz", ".z", ".zip", ".zipx"}
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	common.SetUlimitMax()
}

// enrich HTTP的response: ip\Cert\tld
func (fetcher *Fetcher) EnrichResponse(response Response) Response {
	if !(fetcher.WithTld || fetcher.WithIP) { // do nothing
		return response
	}
	parsedUrl, _ := url.Parse(response.SourceURL)
	var host string
	if strings.Contains(host, ":") {
		host = host[0:strings.Index(host, ":")]
	} else {
		host = parsedUrl.Host
	}
	if fetcher.WithIP {
		if ip, err := getRemoteIP(host); err == nil {
			response.IP = ip
		}
	}
	if fetcher.WithTld {
		response.Tld = getTld(response.URL)
	}
	return response
}

func HasDisableExtension(url string) bool {
	for _, item := range disabledExtentions {
		if strings.HasSuffix(strings.ToLower(url), item) {
			return true
		}
	}
	return false
}
func (fetcher *Fetcher) DoHTTPRequest(targetUrl string) Response {
	r := req.New()
	r.Client().Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r.MaxReadSize = 1 * 1024 * 1024 // 1mb
	r.SetTimeout(time.Duration(fetcher.Timeout) * time.Second)
	var (
		rawResp *req.Resp
		err     error
	)
	if fetcher.Retries < 0 {
		fetcher.Retries = 0
	}
	if fetcher.FilterBinaryExtensions && HasDisableExtension(targetUrl) {
		return Response{Succeed: false, URL: targetUrl, SourceURL: targetUrl, Time: JSONTime(time.Now()), ErrorReason: "disabled extensions"}
	}
	for i := 0; i <= fetcher.Retries; i++ {
		rawResp, err = r.Get(targetUrl, req.Header{"User-Agent": fetcher.UserAgent})
		if err == nil {
			break
		}
	}
	if err != nil {
		errorReason := err.Error()
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			errorReason = "Timeout"
		}
		if strings.Contains(strings.ToLower(err.Error()), "connection refused") {
			errorReason = "PortClosed."
		}
		log.Warn().Msgf("failed get %s. error reason: %s", targetUrl, errorReason)
		return Response{Succeed: false, ErrorReason: err.Error(), URL: targetUrl, SourceURL: targetUrl, Time: JSONTime(time.Now())}
	}
	html, _ := rawResp.ToString()
	statusCode := rawResp.Response().StatusCode
	response := Response{
		URL:        rawResp.Request().URL.String(),
		StatusCode: statusCode,
		Succeed:    true,
		Time:       JSONTime(time.Now()),
		SourceURL:  targetUrl,
	}
	if strings.HasPrefix(rawResp.Request().URL.String(), "https://") {
		//var certInterface map[string]interface{}
		//inrec, _ := json.Marshal(rawResp.Response().TLS.PeerCertificates[0])
		//err := json.Unmarshal(inrec, &certInterface)
		if len(rawResp.Response().TLS.PeerCertificates) > 0 {
			response.Cert = rawResp.Response().TLS.PeerCertificates[0].DNSNames // only echo dns name
		}
	}
	if fetcher.IconMode {
		encoded := base64.StdEncoding.EncodeToString(rawResp.Bytes())
		response.B64Content = encoded
		response.Hash = fmt.Sprintf("%x", sha1.Sum(rawResp.Bytes()))
		return response
	}
	if fixedHtml, err := FixEncoding(rawResp.Bytes(), rawResp.Response().Header.Get("Content-Type")); err == nil {
		html = fixedHtml
	}
	if fetcher.WithLinks {
		rawLinks := getLinks(html, targetUrl)
		response.Links = rawLinks
	}
	if fetcher.WithHTML {
		response.Html = html
	}
	if fetcher.WithTitle {
		title := strings.TrimSpace(getTitle(html))
		response.Title = title
	}
	if fetcher.WithHeaders {
		buf := new(bytes.Buffer)
		rawResp.Response().Header.Write(buf)
		response.Headers = buf.String()
	}
	log.Info().Msgf("HTTP Request Succeed %s. [title: %s]", targetUrl, response.Title)
	return fetcher.EnrichResponse(response)
}

//func (fetcher *Fetcher) WriteWithTimeout(conn net.Conn) ([]byte, error) {
//
//}

func (fetcher *Fetcher) DialPortService(hostPort string) (isOpen bool, Service string) {
	d := net.Dialer{Timeout: time.Duration(fetcher.Timeout) * time.Second}
	conn, err := d.Dial("tcp", hostPort)
	if err != nil {
		return false, ""
	} else {
		defer conn.Close()
		_, port, _ := net.SplitHostPort(hostPort)
		if strings.Contains(port, "443") {
			if err := tls.Client(conn, &tls.Config{InsecureSkipVerify: true}).Handshake(); err == nil {
				return true, "https"
			} else {
				log.Warn().Msg(err.Error())
			}
		}
		return true, "http"
	}
}

func (fetcher *Fetcher) Crawl(input chan string, output chan Response, group *sync.WaitGroup) {
	defer group.Done()
	// input chan.
	for {
		select {
		case inputUrl, ok := <-input:
			if !ok {
				return
			} else {
				if fetcher.PreScan { // pre scan mode. 1.  scan ip port before send request.
					if strings.Contains(inputUrl, ":") {
						isOpen, Service := fetcher.DialPortService(inputUrl)
						if isOpen {
							log.Debug().Msgf("found [%s] %s port open ", Service, inputUrl)
						}
						if isOpen && strings.Contains(Service, "http") {
							response := fetcher.DoHTTPRequest(fmt.Sprintf("%s://%s", Service, inputUrl))
							output <- response
						}
						continue
					}
				}
				if !strings.HasPrefix(strings.ToLower(inputUrl), "http:") && !strings.HasPrefix(strings.ToLower(inputUrl), "https:") {
					service := "http"
					if fetcher.DefaultHTTPS {
						service = "https"
					}
					inputUrl = fmt.Sprintf("%s://%s", service, inputUrl)
				}
				response := fetcher.DoHTTPRequest(fmt.Sprintf("%s%s", inputUrl, fetcher.Endpoint))
				output <- response
			}
		}
	}
}
func (fetcher *Fetcher) DNSInfo(input chan string, output chan DNSInfo, group *sync.WaitGroup) {
	defer group.Done()
	// input chan.
	for {
		select {
		case inputDomain, ok := <-input:
			if ok {
				if response, err := getRemoteIP(inputDomain); err == nil {
					output <- DNSInfo{
						Domain: inputDomain,
						IP:     response,
					}
				}
			} else {
				return
			}
		}
	}
}

func (fetcher *Fetcher) OutputWorker(output chan Response, group *sync.WaitGroup) {
	defer group.Done()
	//var pipe io.Writer
	var pipe *os.File
	var err error
	if fetcher.OutputFileName == "-" {
		pipe = os.Stdout
	} else {
		pipe, err = os.OpenFile(fetcher.OutputFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal()
		}
	}
	if fetcher.OutPutTable {
		fetcher.OutputTable(pipe, output)
	} else {
		fetcher.OutPutJson(pipe, output)
	}
	defer pipe.Close()
	defer pipe.Sync()
}

func (fetcher *Fetcher) Process() {
	inputChan := make(chan string, fetcher.ProcessNum*4)
	outputChan := make(chan Response, fetcher.ProcessNum*4)
	fetchWg := sync.WaitGroup{}
	outputWg := sync.WaitGroup{}
	if fetcher.ProcessNum <= 0 {
		fetcher.ProcessNum = 10
	}
	for i := 0; i < fetcher.ProcessNum; i++ {
		go fetcher.Crawl(inputChan, outputChan, &fetchWg) // start all workers
		fetchWg.Add(1)
	}
	outputWg.Add(1)
	go fetcher.OutputWorker(outputChan, &outputWg)
	var scanner *os.File
	if fetcher.InputFileName == "-" {
		scanner = os.Stdin
	} else {
		var err error
		scanner, err = os.Open(fetcher.InputFileName)
		if err != nil {
			log.Fatal()
		}
	}
	f := bufio.NewScanner(scanner)
	for f.Scan() {
		inputUrl := f.Text()
		if fetcher.PreScan {
			if common.IsCIDR(inputUrl) {
				ips, err := common.GenerateIP(inputUrl)
				if err == nil {
					log.Info().Msgf("convert cidr [%s] -> [%d] ip", inputUrl, len(ips))
					for _, ip := range ips {
						for _, port := range strings.Split(fetcher.Ports, ",") {
							inputChan <- fmt.Sprintf("%s:%s", ip, port)
						}
					}
				} else {
					log.Err(err)
				}
			} else if common.IsIp(inputUrl) {
				for _, port := range strings.Split(fetcher.Ports, ",") {
					inputChan <- fmt.Sprintf("%s:%s", inputUrl, port)
				}
			} else {
				inputChan <- inputUrl
			}
		} else {
			inputChan <- inputUrl
		}

	}
	close(inputChan)
	fetchWg.Wait()
	close(outputChan)
	outputWg.Wait()
	log.Debug().Msg("exit.")
}

func (fetcher *Fetcher) NeedFetch() bool { // 是否需要发送HTTP请求
	if fetcher.WithHTML || fetcher.WithTitle || fetcher.WithLinks {
		return true
	} else {
		return false
	}
}
