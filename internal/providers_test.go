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
	"context"
	"flag"
	"regexp"
	"testing"
	"time"
)

var (
	re           = regexp.MustCompile(`http://((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])):([0-9]{2,5})`)
	testRresults = flag.Bool("verify", false, "test for whether providers return results instead of just checking format.")
)

// TODO: Not sure if these are very idiomatic. Maybe use test table for providers instead of separate functions but still perform each test regardless of success.

func TestKuaidailiP(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := KuaidailiP(ctx)
	name := "KuaidailiP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := FeiyiproxyP(ctx)
	name := "FeiyiproxyP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := YipP(ctx)
	name := "YipP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := Ip3366P(ctx)
	name := "Ip3366P(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := GithubClarketmP(ctx)
	name := "GithubClarketmP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := ProxP(ctx)
	name := "ProxP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := ProxyListP(ctx)
	name := "ProxyListP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := MyProxyP(ctx)
	name := "MyProxyP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := XseoP(ctx)
	name := "XseoP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := ProxylistDownloadP(ctx)
	name := "ProxylistDownloadP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := FreeproxylistsP(ctx)
	name := "FreeproxylistsP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := AliveproxyP(ctx)
	name := "AliveproxyP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := BlogspotP(ctx)
	name := "BlogspotP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := WebanetlabsP(ctx)
	name := "WebanetlabsP(ctx)"
	if len(results) == 0 {

		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := CheckerproxyP(ctx)
	name := "CheckerproxyP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	results := ProxylistMeP(ctx)
	name := "ProxylistMeP(ctx)"
	if len(results) == 0 {
		if *testRresults {
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
