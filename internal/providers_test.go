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
	"flag"
	"regexp"
	"testing"
)

var (
	re          = regexp.MustCompile(`http://((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])):([0-9]{2,5})`)
	testResults = flag.Bool("verify", false, "test for whether providers return results instead of just checking format.")

)

// TODO: Not sure if these are very idiomatic. Maybe use test table for providers instead of separate functions but still perform each test regardless of success.

func TestKuaidailiP(t *testing.T) {
	results := KuaidailiP(100)
	name := "KuaidailiP(100)"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}


func TestFeiyiproxyP(t *testing.T) {
	results := FeiyiproxyP()
	name := "FeiyiproxyP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}


func TestYipP(t *testing.T) {
	results := YipP()
	name := "YipP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}


func TestIp3366P(t *testing.T) {
	results := Ip3366P()
	name := "Ip3366P()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}

func TestGithubClarketmP(t *testing.T) {
	results := GithubClarketmP()
	name := "GithubClarketmP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestProxP(t *testing.T) {
	results := ProxP()
	name := "ProxP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestProxyListP(t *testing.T) {
	results := ProxyListP()
	name := "ProxyListP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestMyProxyP(t *testing.T) {
	results := MyProxyP()
	name := "MyProxyP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestXseoP(t *testing.T) {
	results := XseoP()
	name := "XseoP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestProxylistDownloadP(t *testing.T) {
	results := ProxylistDownloadP()
	name := "ProxylistDownloadP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestFreeproxylistsP(t *testing.T) {
	results := FreeproxylistsP()
	name := "FreeproxylistsP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestAliveproxyP(t *testing.T) {
	results := AliveproxyP()
	name := "AliveproxyP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}

func TestBlogspotP(t *testing.T) {
	results := BlogspotP()
	name := "BlogspotP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}


func TestWebanetlabsP(t *testing.T) {
	flag.Parse()
	results := WebanetlabsP()
	name := "WebanetlabsP()"
	if len(results) == 0 {

		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
func TestCheckerproxyP(t *testing.T) {
	results := CheckerproxyP()
	name := "CheckerproxyP()"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}

func TestProxylistMe(t *testing.T) {
	results := ProxylistMeP(2)
	name := "ProxylistMeP(2)"
	if len(results) == 0 {
		if *testResults {
			t.Errorf("%s didn't return results.", name)
		}
		return
	}
	if !re.MatchString(results[0].Proxy) {
		t.Errorf("%s sample = %v; expected url pattern matching http://121.139.218.165:31409", name, results[0].Proxy)
	} else {
		t.Logf("%s sample = %v \t found = %v", name, results[0].Proxy, len(results))
	}
}
