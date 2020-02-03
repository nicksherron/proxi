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
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

var (
	DB              *sql.DB
	DbPath          string
	connectionLimit int
)

//Model gets embedded into Proxy
type Model struct {
	ID        uint      `json:"-" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `json:"-"`
}

// Proxy represents a proxy record and is used to create the table for our db.
type Proxy struct {
	Model
	AvgResponse  string `gorm:"-" json:"avg_response"`
	ResponseTime int64    `json:"-" gorm:"default:0"`
	CheckCount   uint   `json:"check_count" gorm:"default:0"`
	Country      string `json:"country" `
	FailCount    uint   `json:"fail_count" gorm:"default:0"`
	LastStatus   string `json:"last_status"`
	Proxy        string `json:"proxy" gorm:"type:varchar(100);unique_index"`
	TimeoutCount uint   `json:"timeout_count" gorm:"default:0"`
	Source       string `json:"source"`
	SuccessCount uint   `json:"success_count" gorm:"default:0"`
	Anonymous    bool   `json:"anonymous"`
	LosingStreak uint   `json:"-" gorm:"default:0"`
	Deleted      bool   `json:"-" gorm:"default:false"`
	Judge        string `json:"-"`
}

// Proxies is a slice of Proxy
type Proxies []*Proxy

type TableStats struct {
	Anon            int   `json:"anon"`
	Good            int   `json:"good"`
	Timeout         int   `json:"timeout"`
	Total           int   `json:"total"`
	RecentlyChecked int64 `json:"recently_checked"`
}

// DbInit initializes our db.
func DbInit() {
	// GormDB contains DB connection state
	var gormdb *gorm.DB

	var err error
	if strings.HasPrefix(DbPath, "postgres://") {
		//
		DB, err = sql.Open("postgres", DbPath)
		if err != nil {
			log.Fatal(err)
		}

		gormdb, err = gorm.Open("postgres", DbPath)
		if err != nil {
			log.Fatal(err)
		}
		connectionLimit = 50
	} else {
		DbPath = fmt.Sprintf("file:%v?cache=shared&mode=rwc", DbPath)
		DB, err = sql.Open("sqlite3", DbPath)
		if err != nil {
			log.Fatal(err)
		}
		gormdb, err = gorm.Open("sqlite3", DbPath)
		if err != nil {
			log.Fatal(err)
		}
		DB.Exec("PRAGMA journal_mode=WAL;")
		connectionLimit = 1

	}
	DB.SetMaxOpenConns(connectionLimit)
	gormdb.AutoMigrate(&Proxy{})
	gormdb.Model(&Proxy{}).AddIndex("idx_proxy_compound", "deleted", "last_status", "anonymous", "country")
	// just need gorm for migration.
	gormdb.Close()

	if connectionLimit != 1 {
		DB.Exec(`
			create or replace view proxies_stats as
			select (select count(*) from proxies where deleted = false And last_status = 'good' AND anonymous) as anon,
			       (select count(*) from proxies where deleted = false And last_status = 'good')               as good,
			       (select count(*) from proxies where deleted = false And last_status = 'timeout')            as timeout,
			       (select count(*) from proxies)                                                              as total;`)
	} else {
		DB.Exec(`
			create view if not exists proxies_stats  as
			select (select count(*) from proxies where deleted = false And last_status = 'good' AND anonymous) as anon,
			       (select count(*) from proxies where deleted = false And last_status = 'good')               as good,
			       (select count(*) from proxies where deleted = false And last_status = 'timeout')            as timeout,
			       (select count(*) from proxies)                                                              as total;`)
	}

	dbCacheStats()
}

// DbPing pings DB and either prints "Pong" or does nothing.
func DbPing() {
	DbInit()
	pingErr := DB.Ping()
	if pingErr == nil {
		fmt.Println("Pong")
	}
}

// dbStats
var stats TableStats

func dbPrepWrite() {
	dbCacheStats()
}

func dbCacheStats() {
	err := DB.QueryRow(`select "anon", "good", "timeout", "total" 
							  from proxies_stats`).Scan(&stats.Anon, &stats.Good, &stats.Timeout, &stats.Total)
	if err != nil {
		log.Fatal(err)
	}
}

//--------------------------------------------------------------------------------------

func loadDb(proxy *Proxy)  {
	_, err := DB.Exec(`insert into proxies("created_at", "updated_at", "check_count", "country", "fail_count",
 							"last_status", "proxy", "timeout_count", "source", "success_count", "anonymous", "losing_streak")
 							VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
 							ON CONFLICT (proxy) DO UPDATE SET updated_at = EXCLUDED.updated_at
 							`, time.Now(), time.Now(), proxy.CheckCount, proxy.Country, proxy.FailCount,
		proxy.LastStatus, proxy.Proxy, proxy.TimeoutCount, proxy.Source, proxy.SuccessCount, proxy.Anonymous, proxy.LosingStreak)
	if err != nil {
		log.Fatal(err)
	}
}

func dbInsert(proxy *Proxy) {
	_, err := DB.Exec(`update proxies SET "updated_at" = $1, "check_count" = $2 ,"fail_count" = $3,
 							"last_status" = $4, "timeout_count" = $5, "success_count" = $6, "losing_streak" = $7,
 							 "deleted" = $8,  "anonymous" = $9 , "proxy" = $10, judge = $11, "response_time" = $12 where id = $13`,
		time.Now(), proxy.CheckCount, proxy.FailCount, proxy.LastStatus, proxy.TimeoutCount,
		proxy.SuccessCount, proxy.LosingStreak, proxy.Deleted, proxy.Anonymous, proxy.Proxy, proxy.Judge, proxy.ResponseTime, proxy.ID)
	if err != nil {
		log.Println(err)
	}
}

func dbFind() Proxies {
	var out Proxies
	rows, err := DB.Query(`SELECT "response_time", "id", "check_count", "fail_count","proxy",
 										"timeout_count", "success_count", "losing_streak" FROM proxies where deleted = false`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var row Proxy
		err = rows.Scan(&row.ResponseTime, &row.ID, &row.CheckCount, &row.FailCount, &row.Proxy, &row.TimeoutCount,
			&row.SuccessCount, &row.LosingStreak)
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, &row)
	}
	return out
}

//--------------------------------------------------------------------------------------

func findProxy(p string) interface{} {
	var row Proxy
	err := DB.QueryRow(`select "response_time",  "anonymous",   "check_count",   "country",   "created_at",   "fail_count",   "id",
   						       "last_status",   "proxy",   "source",   "success_count",   "timeout_count",
   						      "updated_at"  from proxies where proxy = `, p).Scan(&row)

	if err != nil {
		log.Fatal(err)
	}

	if row.Proxy != "" {
		row.AvgResponse = fmt.Sprintf("%v", time.Duration(row.ResponseTime)*time.Nanosecond)
		return row
	}
	return nil
}

func getProxyN(num int64, c *gin.Context) Proxies {
	var (
		proxies Proxies
		rows    *sql.Rows
		err     error
	)
	_, anon := c.GetQuery("anon")
	country, countryBool := c.GetQuery("country")
	country = strings.ToUpper(country)

	if anon {
		if countryBool {
			// better performance with sub queries, see https://stackoverflow.com/a/24591688.
			rows, err = DB.Query(`select "response_time", "anonymous",   "check_count",   "country",   "created_at",   "fail_count",   "id",
   						       "last_status",   "proxy",   "source",   "success_count",   "timeout_count",
   						      "updated_at"  from proxies 
   						      where id in (select id from proxies where last_status = 'good' and country = $1 and anonymous order by random() limit $2)`, country, num)
		} else {
			rows, err = DB.Query(`select "response_time", "anonymous",   "check_count",   "country",   "created_at",   "fail_count",   "id",
   						       "last_status",   "proxy",   "source",   "success_count",   "timeout_count",
   						      "updated_at"  from proxies 
   						      where id in (select id from proxies where last_status = 'good' and anonymous order by random() limit $1)`, num)
		}
	} else {
		if countryBool {
			rows, err = DB.Query(`select "response_time", "anonymous",   "check_count",   "country",   "created_at",   "fail_count",
							   "id", "last_status",   "proxy",   "source",   "success_count",   "timeout_count", "updated_at"
							     from proxies where id in (select id from proxies where last_status = 'good' and
							      country = $1 order by random() limit $2)`, country, num)
		} else {
			rows, err = DB.Query(`select "response_time", "anonymous",   "check_count",   "country",   "created_at",   
								"fail_count",   "id","last_status",   "proxy",   "source",   "success_count",
								"timeout_count","updated_at"  from proxies where id in 
								(select id from proxies where last_status = 'good'  order by random() limit $1)`, num)
		}
	}

	defer rows.Close()

	for rows.Next() {
		var row Proxy
		err = rows.Scan(&row.ResponseTime, &row.Anonymous, &row.CheckCount, &row.Country, &row.CreatedAt, &row.FailCount, &row.ID,
			&row.LastStatus, &row.Proxy, &row.Source, &row.SuccessCount, &row.TimeoutCount, &row.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		row.AvgResponse = fmt.Sprintf("%v", time.Duration(row.ResponseTime)*time.Nanosecond)
		proxies = append(proxies, &row)
	}
	return proxies
}

func getProxyAll() Proxies {
	var proxies Proxies
	rows, err := DB.Query(`select "response_time", "anonymous",   "check_count",   "country",   "created_at",   "fail_count",   "id",
   						       "last_status",   "proxy",   "source",   "success_count",   "timeout_count",
   						      "updated_at"  from proxies)`)
	defer rows.Close()
	for rows.Next() {
		var row Proxy
		err = rows.Scan(&row.ResponseTime, &row.Anonymous, &row.CheckCount, &row.Country, &row.CreatedAt, &row.FailCount, &row.ID,
			&row.LastStatus, &row.Proxy, &row.Source, &row.SuccessCount, &row.TimeoutCount, &row.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		//row.AvgResponse = time.Duration(time.Duration(row.ResponseTime) * time.Nanosecond).String()
		proxies = append(proxies, &row)
	}
	return proxies

}

func deleteProxy(p string) int64 {
	row, err := DB.Exec(`delete from proxies where proxy = $1`, p)
	if err != nil {
		log.Fatal(err)
	}
	result, err := row.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func getStats() TableStats {

	if DB.Stats().InUse != DB.Stats().MaxOpenConnections {
		dbCacheStats()
	}

	stats.RecentlyChecked = atomic.LoadInt64(&testCount)
	return stats
}
