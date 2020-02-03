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
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/icrowley/fake"
)

var (
	mutex  = &sync.Mutex{}
	busy   bool
	wgD    sync.WaitGroup
	reader io.ReadCloser
)

func findSubmatchRange(regex *regexp.Regexp, str string) []string {
	var matched []string
	for _, matches := range regex.FindAllString(str, -1) {
		match := regex.FindStringSubmatch(matches)[1]
		matched = append(matched, match)
	}
	return matched
}

func findAllTemplate(pattern *regexp.Regexp, html string, template string) []string {
	var (
		results []string
		result  []byte
	)

	for _, matches := range pattern.FindAllStringSubmatchIndex(html, -1) {
		result = pattern.ExpandString(result, template, html, matches)
	}
	for _, newLine := range strings.Split(string(result), "\n") {
		results = append(results, newLine)
	}
	return results
}

func get(u string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()
	client := &http.Client{
		//Timeout: 60 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return "", err
	}
	//req.Header.Set("X-Forwarded-For", fake.IPv4())
	req.Header.Set("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	body, err := ioutil.ReadAll(reader)

	if err != nil {
		return "", err
	}

	return string(body), nil

}

func getX(u string) (string, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Forwarded-For", fake.IPv4())
	req.Header.Set("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// DownloadProxies downloads proxies from providers.
func DownloadProxies() Proxies {
	log.Println("\rStarting proxy downloads...")
	wgD.Add(16)
	var providerProxies Proxies

	// Download from providers
	go func() {
		defer wgD.Done()
		results := FreeproxylistsP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := WebanetlabsP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := CheckerproxyP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := ProxyListP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := AliveproxyP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := KuaidailiP(10)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := FeiyiproxyP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := YipP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := Ip3366P()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := ProxylistMeP(10)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := ProxylistDownloadP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := BlogspotP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := ProxP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := MyProxyP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := XseoP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		results := GithubClarketmP()
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()

	wgD.Wait()
	return providerProxies

}

// DownloadInit initializes DownloadProxies and saves results to GormDB
func DownloadInit() {
	busy = true
	validMaxmind = true
	providerResults := DownloadProxies()
	ipDB, err := maxmindDb()
	if err != nil {
		validMaxmind = false
	}
	var tmpfile *os.File
	dumpResults := false
	if os.Getenv("PROXYPOOL_DUMP") == "1" {
		tmpfile, err = ioutil.TempFile("", "proxi-dump.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		dumpResults = true

	}
	dbPrepWrite()
	for _, v := range providerResults {
		outIP := strings.Split(strings.ReplaceAll(v.Proxy, "http://", ""), ":")[0]
		ip := net.ParseIP(outIP)
		if ip == nil {
			continue
		}
		if validMaxmind {
			country, err := ipDB.Country(ip)
			check(err)
			v.Country = country.Country.IsoCode
		}
		loadDb(v)
		if dumpResults {
			fmt.Fprintln(tmpfile, v)
		}
	}
	CheckInit()
	fmt.Println("Done")
}
