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

func getX(u string) (string, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
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

func getKuaidaili(u string) (string, error) {
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

// DownloadProxies downloads proxies from providers.
func DownloadProxies() Proxies {
	log.Println("\rStarting proxy downloads...")
	wgD.Add(16)
	var providerProxies Proxies

	ctxTimeout := time.Minute
	// Download from providers
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := FreeproxylistsP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := WebanetlabsP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := CheckerproxyP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := ProxyListP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := AliveproxyP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := KuaidailiP(ctx, 0)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := FeiyiproxyP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := YipP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := Ip3366P(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := ProxylistMeP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := ProxylistDownloadP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := BlogspotP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := ProxP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := MyProxyP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := XseoP(ctx)
		mutex.Lock()
		providerProxies = append(providerProxies, results...)
		mutex.Unlock()
	}()
	go func() {
		defer wgD.Done()
		ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()
		results := GithubClarketmP(ctx)
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
