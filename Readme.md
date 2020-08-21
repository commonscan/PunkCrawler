# PunkCrawler
PunkCrawler is a WebSpider with high speed(Actually it is a demo project to learn go for me).
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
|  3 | http://210.39.2.207:80/    |    200 | 深圳市教师教育网         | 2020-08-21 17:43:04 |
|  4 | https://210.39.2.207:443/  |    200 | 深圳市教师教育网         | 2020-08-21 17:43:05 |
|  5 | http://210.39.3.178:80/    |    404 |                          | 2020-08-21 17:43:06 |
|  6 | https://210.39.3.178:443/  |    404 |                          | 2020-08-21 17:43:06 |
|  7 | http://210.39.5.3:80/      |    200 |                          | 2020-08-21 17:43:07 |
|  8 | http://210.39.5.7:80/      |    500 |                          | 2020-08-21 17:43:07 |
|  9 | http://210.39.5.8:80/      |    200 |                          | 2020-08-21 17:43:07 |
| 10 | https://210.39.4.100:443/  |    200 |                          | 2020-08-21 17:43:07 |
| 11 | http://210.39.5.9:80/      |    404 |                          | 2020-08-21 17:43:07 |
| 12 | http://210.39.5.24:80/     |    503 |                          | 2020-08-21 17:43:07 |
| 13 | http://210.39.5.34:80/     |    404 |                          | 2020-08-21 17:43:07 |
| 14 | http://210.39.5.43:80/     |    500 |                          | 2020-08-21 17:43:07 |
| 16 | https://210.39.5.8:443/    |    200 |                          | 2020-08-21 17:43:08 |
| 17 | http://210.39.8.10:80/     |    404 |                          | 2020-08-21 17:43:11 |
| 18 | http://210.39.8.22:80/     |    200 | 操作失败                 | 2020-08-21 17:43:11 |
| 19 | https://210.39.8.22:443/   |    200 | 操作失败                 | 2020-08-21 17:43:12 |
| 21 | https://210.39.9.3:443/    |    200 | NSFOCUS&nbsp;SAS[H]      | 2020-08-21 17:43:13 |
| 22 | https://210.39.9.2:443/    |    200 |                          | 2020-08-21 17:43:13 |
| 23 | http://210.39.9.1:80/      |    200 | 统一身份认证             | 2020-08-21 17:43:15 |
| 24 | https://210.39.9.1:443/    |    200 | 统一身份认证             | 2020-08-21 17:43:15 |
| 25 | https://210.39.5.24:443/   |    200 | 统一身份认证             | 2020-08-21 17:43:16 |
| 26 | http://210.39.12.36:80/    |    200 | IBM HTTP Server          | 2020-08-21 17:43:16 |
| 27 | http://210.39.12.28:80/    |    200 |                          | 2020-08-21 17:43:16 |
| 28 | http://210.39.12.247:80/   |    404 | 404 Unknown Virtual Host | 2020-08-21 17:43:17 |
| 29 | https://210.39.12.28:443/  |    200 | Welcome to nginx!        | 2020-08-21 17:43:17 |
| 30 | https://210.39.12.247:443/ |    404 | 404 Unknown Virtual Host | 2020-08-21 17:43:18 |
| 31 | https://210.39.12.250:443/ |    200 | 统一身份认证             | 2020-08-21 17:43:19 |
| 32 | http://210.39.12.250:80/   |    200 | 统一身份认证             | 2020-08-21 17:43:19 |
+----+----------------------------+--------+--------------------------+---------------------+
```

Thanks for [zgrab2](https://github.com/zmap/zgrab2)!
tested on go1.14.4.

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
