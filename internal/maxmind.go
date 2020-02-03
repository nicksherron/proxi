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
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/oschwald/geoip2-golang"
)

var (
	// MaxmindFilePath is the path to maxmind country db file. Used to resolve the geolocation of proxies.
	MaxmindFilePath    string
	maxmindDownloadURL = "https://httpbin.net/GeoLite2-Country.mmdb"
	validMaxmind       bool
)

func downloadFile(f string) (string, error) {
	out, err := os.Create(f)
	defer out.Close()

	// Get the data
	resp, err := http.Get(maxmindDownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return out.Name(), nil
}

func maxmindDb() (*geoip2.Reader, error) {
	f := MaxmindFilePath
	var M string

	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		M, err = downloadFile(f)
		if err != nil {
			return nil, err
		}
	} else {
		M = f
	}
	db, err := geoip2.Open(M)
	if err != nil {
		return nil, err
	}
	return db, nil

}
