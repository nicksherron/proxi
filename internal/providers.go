/*
 * Copyright © 2020 nicksherron <nsherron90@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package internal

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

var (
	// Matches ip and port
	reProxy       = regexp.MustCompile(`(?ms)(?P<ip>(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?))(?:.*?(?:(?:(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?))|(?P<port>\d{2,5})))`)
	templateProxy = "http://${ip}:${port}\n"
)

func FreeproxylistsP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "freeproxylists.com"
		w            sync.WaitGroup
		fplReID      = regexp.MustCompile(`(?m)href\s*=\s*['"](?P<type>[^'"]*)/(?P<id>\d{10})[^'"]*['"]`)
		fplUrls      = []string{
			"http://www.freeproxylists.com/anonymous.html",
			"http://www.freeproxylists.com/elite.html",
		}
	)
	for _, u := range fplUrls {
		body, err := get(u)
		if err != nil {
			continue
		}
		template := "http://www.freeproxylists.com/load_${type}_${id}.html\n"
		matches := findAllTemplate(fplReID, body, template)
		for _, match := range matches {
			w.Add(1)
			go func(endpoint string) {
				defer w.Done()
				ipList, err := get(endpoint)
				if err != nil {
					return
				}
				matched := findAllTemplate(reProxy, ipList, templateProxy)

				for _, proxy := range matched {
					if proxy == "" {
						continue
					}
					p := Proxy{Proxy: proxy, Source: source}
					mu.Lock()
					foundProxies = append(foundProxies, &p)
					mu.Unlock()
				}
			}(match)
		}
		w.Wait()
	}

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func WebanetlabsP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "webanetlabs.net"
		w            sync.WaitGroup
		re           = regexp.MustCompile(`(?m)href\s*=\s*['"]([^'"]*proxylist_at_[^'"]*)['"]`)
		url          = "https://webanetlabs.net/publ/24"
	)
	body, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(re, body) {
		w.Add(1)
		go func(page string) {
			defer w.Done()
			// https://webanetlabs.net/freeproxyweb/proxylist_at_02.11.2019.txt
			u := "https://webanetlabs.net" + page
			ipList, err := get(u)
			if err != nil {
				return
			}
			for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(href)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func CheckerproxyP() Proxies {

	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "checkerproxy.net"
		w            sync.WaitGroup
		re           = regexp.MustCompile(`(?m)href\s*=\s*['"](/archive/\d{4}-\d{2}-\d{2})['"]`)
		url          = "https://checkerproxy.net/"
	)
	body, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(re, body) {
		w.Add(1)
		go func(endpoint string) {

			defer w.Done()
			u := "https://checkerproxy.net/api" + endpoint

			res, err := http.Get(u)
			if err != nil {
				return
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return
			}
			result := gjson.GetBytes(body, "#.addr")

			result.ForEach(func(key, value gjson.Result) bool {
				proxy := fmt.Sprintf("http://%v", value.String())

				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
				return true // keep iterating
			})

		}(href)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func ProxyListP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "proxy-list.org"
		ipBase64     = regexp.MustCompile(`Proxy\('([\w=]+)'\)`)
		w            sync.WaitGroup
	)
	w.Add(10)
	for i := 1; i < 11; i++ {
		go func(page int) {
			defer w.Done()
			u := fmt.Sprintf("http://proxy-list.org/english/index.php?p=%v", page)
			ipList, err := get(u)
			if err != nil {
				return
			}
			for _, match := range findSubmatchRange(ipBase64, ipList) {
				if match == "" {
					continue
				}
				decoded, err := base64.StdEncoding.DecodeString(match)
				check(err)
				proxy := fmt.Sprintf("http://%v", string(decoded))
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(i)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func AliveproxyP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "aliveproxy.com"
		w            sync.WaitGroup
		suffixes     = []string{
			//"socks5-list",
			"high-anonymity-proxy-list",
			"anonymous-proxy-list",
			"fastest-resolver",
			"us-proxy-list",
			"gb-proxy-list",
			"fr-proxy-list",
			"de-proxy-list",
			"jp-proxy-list",
			"ca-proxy-list",
			"ru-proxy-list",
			"proxy-list-port-80",
			"proxy-list-port-81",
			"proxy-list-port-3128",
			"proxy-list-port-8000",
			"proxy-list-port-8080",
		}
		re = regexp.MustCompile(`(?P<ip>(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])):(?P<port>[0-9]{2,5})`)
	)

	for _, href := range suffixes {
		w.Add(1)
		u := fmt.Sprintf("http://www.aliveproxy.com/%v/", href)
		go func(endpoint string) {
			defer w.Done()
			ipList, err := get(endpoint)
			if err != nil {
				return
			}
			for _, proxy := range findAllTemplate(re, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(u)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func FeiyiproxyP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		source       = "feiyiproxy.com"
		baseUrl      = "http://www.feiyiproxy.com/?page_id=1457"
	)
	ipList, err := get(baseUrl)
	if err != nil {
		return Proxies{}
	}
	for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
		if proxy == "" {
			continue
		}
		p := Proxy{Proxy: proxy, Source: source}
		foundProxies = append(foundProxies, &p)
	}
	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
			fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))
		}

	}
	return foundProxies
}

func YipP() Proxies {
	start := time.Now()
	var (
		largest      int
		foundProxies Proxies
		mu           sync.Mutex
		source       = "7yip.cn"
		w            sync.WaitGroup
		ints         []int
		reHref       = regexp.MustCompile(`(?ms)<li><a href="\?action=china&page=(\d+)">\d?</a></li>`)
		url          = "https://www.7yip.cn/free/?page=1"
	)

	body, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(reHref, body) {
		i, err := strconv.Atoi(href)
		if err != nil {
			continue
		}
		ints = append(ints, i)
	}
	if len(ints) == 0 {
		return Proxies{}
	}
	sort.Ints(ints)
	largest = ints[len(ints)-1]
	largest++
	for i := 1; i < largest; i++ {
		w.Add(1)
		go func(page int) {
			defer w.Done()
			u := fmt.Sprintf("https://www.7yip.cn/free/?page=%v", page)
			ipList, err := get(u)
			if err != nil {
				return
			}
			for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(i)
	}
	w.Wait()
	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
			fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))
		}

	}
	return foundProxies
}

func Ip3366P() Proxies {
	start := time.Now()
	var (
		largest      int
		foundProxies Proxies
		mu           sync.Mutex
		source       = "ip3366.net"
		w            sync.WaitGroup
		ints         []int
		reHref       = regexp.MustCompile(`(?ms)<a href="\?stype=1&page=(\d+)">`)
		url          = "http://www.ip3366.net/free/?stype=1&page=1"
	)

	body, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(reHref, body) {
		i, err := strconv.Atoi(href)
		if err != nil {
			continue
		}
		ints = append(ints, i)
	}
	if len(ints) == 0 {
		return Proxies{}
	}
	sort.Ints(ints)
	largest = ints[len(ints)-1]
	largest++
	for i := 1; i < largest; i++ {
		w.Add(1)
		go func(page int) {
			defer w.Done()
			u := fmt.Sprintf("http://www.ip3366.net/free/?stype=1&page=%v", page)
			ipList, err := get(u)
			if err != nil {
				return
			}
			for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(i)
	}
	w.Wait()
	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
			fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))
		}

	}
	return foundProxies
}

func KuaidailiP(resultLimit int) Proxies {
	start := time.Now()
	var (
		largest      int
		foundProxies Proxies
		mu           sync.Mutex
		source       = "kuaidaili.com"
		w            sync.WaitGroup
		ints         []int
		reHref       = regexp.MustCompile(`(?m)<a href="/free/inha/(\d+)/">`)
		url          = "https://www.kuaidaili.com/free/inha/1/"
	)
	body, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(reHref, body) {
		i, err := strconv.Atoi(href)
		if err != nil {
			continue
		}
		ints = append(ints, i)
	}
	if len(ints) == 0 {
		return Proxies{}
	}
	sort.Ints(ints)
	// Limit urls to visit if too many results to wait for when testing and what not.
	if resultLimit != 0 {
		largest = resultLimit
	} else {
		largest = ints[len(ints)-1]
		largest++
	}
	counter := 0
	for i := 1; i < largest; i++ {
		w.Add(1)
		counter++
		go func(page int) {
			defer w.Done()
			u := fmt.Sprintf("https://www.kuaidaili.com/free/inha/%v/", page)
			ipList, err := get(u)
			if err != nil {
				return
			}
			for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(i)
		// only 25 goroutines at a time. (1170 urls to get)
		if counter >= 25 {
			w.Wait()
			counter = 0
		}
	}
	w.Wait()
	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))
	}
	return foundProxies
}

func ProxylistMeP(resultLimit int) Proxies {
	start := time.Now()
	var (
		largest      int
		foundProxies Proxies
		mu           sync.Mutex
		source       = "proxylist.me"
		w            sync.WaitGroup
		ints         []int
		reHref       = regexp.MustCompile(`(?m)href\s*=\s*['"][^'"]*/?page=(\d+)['"]`)
		re           = regexp.MustCompile(`>(?P<ip>(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])):(?P<port>[0-9]{2,5})<`)
		url          = "https://proxylist.me/"
	)
	body, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(reHref, body) {
		i, err := strconv.Atoi(href)
		if err != nil {
			continue
		}
		ints = append(ints, i)
	}
	if len(ints) == 0 {
		return Proxies{}
	}
	sort.Ints(ints)
	// Limit urls to visit if too many results to wait for when testing and what not.
	if resultLimit != 0 {
		largest = resultLimit
	} else {
		largest = ints[len(ints)-1]
		largest++
	}

	counter := 0
	for i := 1; i < largest; i++ {
		w.Add(1)
		counter++
		go func(page int) {
			defer w.Done()
			u := fmt.Sprintf("https://proxylist.me/?page=%v", page)
			ipList, err := get(u)
			if err != nil {
				return
			}
			ipList = strings.ReplaceAll(strings.ReplaceAll(ipList, " ", ""), "\n", "")
			ipList = strings.ReplaceAll(ipList, "</a></td><td>", ":")
			for _, proxy := range findAllTemplate(re, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(i)
		// only 25 goroutines at a time. (1170 urls to get)
		if counter >= 25 {
			w.Wait()
			counter = 0
		}
	}
	w.Wait()
	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func ProxylistDownloadP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		source       = "proxy-list.download"
	)
	body, err := get("https://www.proxy-list.download/api/v1/get?type=http")
	if err != nil {
		return Proxies{}
	}
	for _, proxy := range findAllTemplate(reProxy, body, templateProxy) {
		if proxy == "" {
			continue
		}
		p := Proxy{Proxy: proxy, Source: source}
		foundProxies = append(foundProxies, &p)
	}

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func BlogspotP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "blogspot.com"
		w            sync.WaitGroup
		re           = regexp.MustCompile(`(?m)<a href\s*=\s*['"]([^'"]*\.\w+/\d{4}/\d{2}/[^'"#]*)['"]>`)
		domains      = []string{
			"sslproxies24.blogspot.com",
			"proxyserverlist-24.blogspot.com",
			"freeschoolproxy.blogspot.com",
			"googleproxies24.blogspot.com",
		}
	)
	for _, domain := range domains {
		w.Add(1)
		go func(endpoint string) {
			u := fmt.Sprintf("http://%v/", endpoint)
			defer w.Done()
			urlList, err := get(u)
			if err != nil {
				return
			}
			for _, href := range findSubmatchRange(re, urlList) {
				w.Add(1)
				go func(endpoint string) {
					ipList, err := get(endpoint)
					if err != nil {
						return
					}
					defer w.Done()
					for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
						if proxy == "" {
							continue
						}
						p := Proxy{Proxy: proxy, Source: source}
						mu.Lock()
						foundProxies = append(foundProxies, &p)
						mu.Unlock()
					}
				}(href)
			}
		}(domain)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func ProxP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "prox.com"
		w            sync.WaitGroup
		re           = regexp.MustCompile(`href\s*=\s*['"]([^'"]?proxy_list_high_anonymous_[^'"]*)['"]`)
		url          = "http://www.proxz.com/proxy_list_high_anonymous_0.html"
	)
	urlList, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(re, urlList) {
		w.Add(1)
		go func(endpoint string) {
			u := fmt.Sprintf("http://www.proxz.com/%v", endpoint)
			ipList, err := get(u)
			if err != nil {
				return
			}
			defer w.Done()
			for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(href)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

	}

	return foundProxies
}

func MyProxyP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		mu           sync.Mutex
		source       = "my-proxy.com"
		w            sync.WaitGroup
		re           = regexp.MustCompile(`(?m)href\s*=\s*['"]([^'"]?free-[^'"]*)['"]`)
		url          = "https://www.my-proxy.com/free-proxy-list.html"
	)

	urlList, err := get(url)
	if err != nil {
		return Proxies{}
	}
	for _, href := range findSubmatchRange(re, urlList) {
		w.Add(1)
		go func(endpoint string) {
			u := fmt.Sprintf("https://www.my-proxy.com/%v", endpoint)
			defer w.Done()
			ipList, err := get(u)
			if err != nil {
				return
			}
			for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
				if proxy == "" {
					continue
				}
				p := Proxy{Proxy: proxy, Source: source}
				mu.Lock()
				foundProxies = append(foundProxies, &p)
				mu.Unlock()
			}
		}(href)
	}
	w.Wait()

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
			fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

		}

	}
	return foundProxies
}

func XseoP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		source       = "xseo.in"
		baseUrl      = "http://xseo.in/freeproxy"
		re           = regexp.MustCompile(`(?P<ip>(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])):(?P<port>[0-9]{2,5})`)
	)
	ipList, err := get(baseUrl)
	if err != nil {
		return Proxies{}
	}
	for _, proxy := range findAllTemplate(re, ipList, templateProxy) {
		if proxy == "" {
			continue
		}
		p := Proxy{Proxy: proxy, Source: source}
		foundProxies = append(foundProxies, &p)
	}

	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
			fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))

		}

	}
	return foundProxies
}

func GithubClarketmP() Proxies {
	start := time.Now()
	var (
		foundProxies Proxies
		source       = "github.com/clarketm"
		baseUrl      = "https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt"
	)
	ipList, err := get(baseUrl)
	if err != nil {
		return Proxies{}
	}
	for _, proxy := range findAllTemplate(reProxy, ipList, templateProxy) {
		if proxy == "" {
			continue
		}
		p := Proxy{Proxy: proxy, Source: source}
		foundProxies = append(foundProxies, &p)
	}
	if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
		if os.Getenv("PROXYPOOL_PROVIDER_DEBUG") == "1" {
			fmt.Printf("\n%v\t%v\t%v\n", time.Since(start), source, len(foundProxies))
		}

	}
	return foundProxies
}
