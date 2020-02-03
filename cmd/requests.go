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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

var (
	address string
)

func get(u string) string {
	resp, err := http.Get(u)
	if err != nil {
		fmt.Printf("Request failed for %v are you sure the server is running?\n", u)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func post(u string, data url.Values) string {
	resp, err := http.PostForm(u, data)
	if err != nil {
		fmt.Printf("Request failed for %v are you sure the server is running?\n", u)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func stats() {
	var stat map[string]interface{}
	u := fmt.Sprintf("%v/stats", address)
	if watch == 0 {
		json.Unmarshal([]byte(get(u)), &stat)
		f := colorjson.NewFormatter()
		f.Indent = 2
		s, _ := f.Marshal(stat)
		fmt.Println(string(s))
	} else {
		for {
			fmt.Printf("\033[5;1H")
			fmt.Printf("updating every %d seconds\n\n", watch)
			json.Unmarshal([]byte(get(u)), &stat)
			f := colorjson.NewFormatter()
			f.Indent = 2
			s, _ := f.Marshal(stat)
			fmt.Println(string(s))
			time.Sleep(time.Duration(watch) * time.Second)
		}
	}
}

func getProxy() {

	var (
		proxies []interface{}
		proxy   map[string]interface{}
	)

	if getAll {
		u := fmt.Sprintf("%v/getall", address)
		json.Unmarshal([]byte(get(u)), &proxies)
		f := colorjson.NewFormatter()
		f.Indent = 2
		s, _ := f.Marshal(proxies)
		fmt.Println(string(s))
		return
	}
	v := url.Values{}
	if anon {
		v.Add("anon", "")
	}
	if country != "" {
		v.Add("country", country)
	}

	if numProxies == 1 {
		u := fmt.Sprintf("%v/get?%v", address, v.Encode())
		json.Unmarshal([]byte(get(u)), &proxy)
		f := colorjson.NewFormatter()
		f.Indent = 2
		s, _ := f.Marshal(proxy)
		fmt.Println(string(s))
		return
	}
	u := fmt.Sprintf("%v/get/%v?%v", address, numProxies, v.Encode())
	json.Unmarshal([]byte(get(u)), &proxies)
	f := colorjson.NewFormatter()
	f.Indent = 2
	s, _ := f.Marshal(proxies)
	fmt.Println(string(s))

}

func findProxy(proxy string) {
	var body map[string]interface{}

	u := fmt.Sprintf("%v/find", address)
	v := url.Values{}
	v.Add("proxy", proxy)
	json.Unmarshal([]byte(post(u, v)), &body)
	f := colorjson.NewFormatter()
	f.Indent = 2
	s, _ := f.Marshal(body)
	fmt.Println(string(s))
}

func deleteProxy(proxy string) {

	u := fmt.Sprintf("%v/delete", address)
	v := url.Values{}
	v.Add("proxy", proxy)
	fmt.Println(post(u, v))

}

func getRefresh() {
	u := fmt.Sprintf("%v/refresh", address)
	color.HiGreen(get(u))
}
