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
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicksherron/proxi/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// Addr is the listen and serve address for the server.
	Addr string
	// LogFile is the file location to store the servers http request logs.
	LogFile string
	// Version is the current version of the program. In releases this is set as the git tag via build ldflags.
	Version string
)

type proxyLookup struct {
	Proxy string `form:"proxy" json:"proxy" xml:"proxy"  binding:"required"`
}

func getLogFile() *os.File {
	f, err := os.Create(LogFile)
	if err != nil {
		log.Fatal(err)
	}
	return f

}

// LoggerWithFormatter instance a Logger middleware with the specified log format function.
func loggerWithFormatterWriter(f gin.LogFormatter) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: f,
		Output:    getLogFile(),
	})
}

// API is the rest api/swagger docs that listen and serves forever.
func API() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(loggerWithFormatterWriter(func(param gin.LogFormatterParams) string {
		// my own format
		return fmt.Sprintf("[ProxyPool] %v | %3d | %13v | %15s | %-7s  %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	}))

	r.POST("/delete", func(c *gin.Context) {
		var d proxyLookup
		c.ShouldBind(&d)
		result := deleteProxy(d.Proxy)
		c.String(http.StatusOK, string(result))
	})

	r.POST("/find", func(c *gin.Context) {
		var d proxyLookup
		c.ShouldBind(&d)
		result := findProxy(d.Proxy)
		c.IndentedJSON(http.StatusOK, result)
	})

	r.GET("/get", func(c *gin.Context) {
		var ret *Proxy
		result := getProxyN(1, c)
		if len(result) != 0 {
			ret = result[0]
		}
		c.IndentedJSON(http.StatusOK, ret)
	})

	r.GET("/get/:n", func(c *gin.Context) {
		n := c.Param("n")
		num, _ := strconv.Atoi(n)
		result := getProxyN(int64(num), c)
		c.IndentedJSON(http.StatusOK, result)
	})

	r.GET("/getall", func(c *gin.Context) {
		result := getProxyAll()
		c.IndentedJSON(http.StatusOK, result)
	})

	r.GET("/stats", func(c *gin.Context) {
		result := getStats()
		c.IndentedJSON(http.StatusOK, result)
	})

	r.GET("/db", func(c *gin.Context) {
		result := DB.Stats()
		c.IndentedJSON(http.StatusOK, result)
	})


	r.GET("/refresh", func(c *gin.Context) {
		if busy {
			c.String(http.StatusConflict, "busy")
		} else {
			CheckProxiesVar = true
			DownloadProxiesVar = true
			go DownloadProxies()
			c.String(http.StatusOK, "ok")
		}
	})

	r.GET("/busy", func(c *gin.Context) {
		c.String(http.StatusOK, "%v", busy)
	})

	docs.SwaggerInfo.Host = fmt.Sprintf("http://%v", Addr)
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%v/swagger/doc.json", Addr))
	docs.SwaggerInfo.Version = Version
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

	err := r.Run(Addr)
	if err != nil {
		fmt.Println("Error: \t", err)
	}
}
