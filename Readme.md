# PunkCrawler
PunkCrawler is a WebSpider with high speed(Actually it is a demo project to learn go for me).
## release
```
goreleaser --rm-dist --skip-publish --skip-validate
```

## usage
```bash
echo 202.120.7.213/24  | go run main.go --pre-scan --with-title --timeout=5 --ports 80,443 --table 
```

```output
+----+----------------------------+--------+--------------------------+---------------------+
|  # | URL                        | 状态码 | 标题                     | TIME                |
+----+----------------------------+--------+--------------------------+---------------------+
|  1 | http://210.39.3.5:80/      |    200 |                          | 2020-08-21 17:43:04 |
|  2 | https://210.39.3.5:443/    |    200 |                          | 2020-08-21 17:43:04 |
| 32 | http://210.39.12.250:80/   |    200 | 统一身份认证             | 2020-08-21 17:43:19 |
+----+----------------------------+--------+--------------------------+---------------------+
```
## known issue
日志在Powershell下展示不正常
![image.png](/uploads/77FC6362DCAE448187528D40A454A2C6/image.png)


## main feature
- configurable options for http.
- high speed
- easy to use
- json output

# can be used to
- send lots of http request to a list of urls.
 
## usage
```bash
» go run main.go -h                                                                                                              1 ↵ cat@jinxufang-LC2
Usage:
  main [OPTIONS]

Application Options:
  -i, --input-file=   输入文件名 (default: -)
  -o, --output-file=  输出文件名 (default: -)
  -p, --process_num=  并发数 (default: 100)
  -t, --timeout=      最大超时数(s) (default: 30)
  -r, --retries=      最大重试次数 (default: 2)
      --with-title    是否输出Title
      --with-html     是否输出HTML
      --with-tld      是否输出TLD
      --with-ip       是否输出IP
      --with-headers  是否输出Headers
      --with-links    是否输出链接信息

Help Options:
  -h, --help          Show this help message

» echo http://qq.com | go run main.go --with-title --with-headers | jq                                                               cat@jinxufang-LC2
{
  "url": "http://qq.com",
  "title": "腾讯首页",
  "status_code": 200,
  "time": "2020-05-05 18:06:36",
  "succeed": true,
  "source_url": "http://qq.com",
  "headers": "Cache-Control: max-age=60\r\nConnection: keep-alive\r\nContent-Type: text/html; charset=GB2312\r\nDate: Tue, 05 May 2020 10:06:36 GMT\r\nExpires: Tue, 05 May 2020 10:07:36 GMT\r\nServer: squid/3.5.24\r\nVary: Accept-Encoding\r\n"
}
```
