/*
Copyright © 2020 nicksherron <nsherron90@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"syscall"
	"time"

	isatty "github.com/mattn/go-isatty"
	"github.com/nicksherron/proxi/internal"
	"github.com/spf13/cobra"
)

var (
	traceProfile string
	cpuProfile   string
	memProfile   string
	pingDB       bool
	serverCmd    = &cobra.Command{
		Use:   "server",
		Short: "Download then check proxies and start rest api server for querying results.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
			internal.DbInit()
			oldLimit, newLimit := internal.IncrFdLimit()
			if newLimit != 0 {
				log.Printf("Increased maximum number of open files to %v (it was originally set to %v).",
					newLimit, oldLimit)
			}
			if pingDB {
				internal.DbPing()
				return
			}
			if cpuProfile != "" || memProfile != "" || traceProfile != "" {
				profileInit()
			}
			internal.StartupMessage()
			if internal.DownloadProxiesVar {
				go internal.DownloadInit()
			} else if internal.CheckProxiesVar {
				go internal.CheckInit()
			}
			internal.API()
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&traceProfile, "trace", "", "Write trace profile to file.")
	serverCmd.PersistentFlags().StringVar(&cpuProfile, "cpu", "", "Write cpu profile to file.")
	serverCmd.PersistentFlags().StringVar(&memProfile, "mem", "", "Write memory profile to file.")
	serverCmd.PersistentFlags().BoolVar(&internal.CheckProxiesVar, "check", false, "Check proxies after downloading.")
	serverCmd.PersistentFlags().BoolVar(&internal.DownloadProxiesVar, "download", false, "Download proxies after server start.")
	serverCmd.PersistentFlags().StringVarP(&internal.Addr, "addr", "a", listenAddr(), "Ip and port to listen and serve on.")
	serverCmd.PersistentFlags().StringVar(&internal.MaxmindFilePath, "maxmind-file", maxmindPath(), "Maxmind country db file. Downloads if default doesn't exist.")
	serverCmd.PersistentFlags().StringVar(&internal.DbPath, "db", dbPath(), "Sqlite3 backend storage file location.")
	serverCmd.PersistentFlags().StringVar(&internal.LogFile, "log", logPath(), "Set filepath for HTTP log.")
	serverCmd.PersistentFlags().IntVarP(&internal.Deadline, "deadline", "d", 60, "Deadline time for downloads in seconds. Set to 0 if you don't want any deadline.")
	serverCmd.PersistentFlags().IntVar(&internal.FileLimitMax, "ulimit", 2048, "Number of allowed file handles per process.")
	serverCmd.PersistentFlags().DurationVarP(&internal.Timeout, "timeout", "t", 30*time.Second, "Specify request time out for checking proxies.")
	serverCmd.PersistentFlags().IntVarP(&internal.Workers, "workers", "w", 100, "Number of (goroutines) concurrent requests to make for checking proxies.")
	serverCmd.PersistentFlags().BoolVar(&pingDB, "ping", false, "Ping db and exit.")
	serverCmd.PersistentFlags().BoolVarP(&internal.Progress, "progress", "p", isTerminal(os.Stderr), "Show proxy test progress bar.")
}

func listenAddr() string {
	var a string
	if os.Getenv("PROXYPOOL_ADDRESS") != "" {
		a = os.Getenv("PROXYPOOL_ADDRESS")
		return a
	}
	a = "0.0.0.0:4444"
	return a

}

func maxmindPath() string {
	maxmindFile := "GeoLite2-Country.mmdb"
	f := filepath.Join(dataHome(), maxmindFile)
	return f
}

func dbPath() string {
	dbFile := "data.db"
	f := filepath.Join(dataHome(), dbFile)
	return f
}

func logPath() string {
	logFile := "server.log"
	f := filepath.Join(configHome(), logFile)
	return f
}

func profileInit() {

	go func() {
		defer os.Exit(1)
		if traceProfile != "" {
			f, err := os.Create(traceProfile)
			if err != nil {
				log.Fatal("could not create trace profile: ", err)
			}
			defer f.Close()
			if err := trace.Start(f); err != nil {
				log.Fatal("could not start trace profile: ", err)
			}
			defer trace.Stop()
		}

		if cpuProfile != "" {
			f, err := os.Create(cpuProfile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			defer f.Close()
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
		}

		defer func() {
			if memProfile != "" {
				mf, err := os.Create(memProfile)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				defer mf.Close()
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(mf); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
	}()
}

func isTerminal(f *os.File) bool {
	if runtime.GOOS == "windows" {
		return false
	}

	fd := f.Fd()
	return os.Getenv("TERM") != "dumb" && (isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd))
}
