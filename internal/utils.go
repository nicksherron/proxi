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

	"github.com/fatih/color"
	"github.com/nicksherron/proxi/internal/fdlimit"
)

var (
	//FileLimitMax sets max open file descriptors value
	FileLimitMax int
)

// TODO: Make proper error and logging. This is stupid.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// IncrFdLimit attempts to increase max file handles per process
func IncrFdLimit() (int, uint64){
	var newLimit uint64
	oldLimit, err := fdlimit.Current()
	if err != nil {
		log.Fatalf("Failed to retrieve file descriptor allowance: %v", err)
	}
	if oldLimit < FileLimitMax {
		if newLimit, err = fdlimit.Raise(uint64(FileLimitMax)); err != nil {
			log.Fatalf("Failed to raise file descriptor allowance: %v", err)
		}
	}
	return oldLimit, newLimit
}



//--------------------------------------------------------------------------

// StartupMessage prints startup banner
func StartupMessage() {
	banner := fmt.Sprintf(`
                     _           
                    (_)       
 _ __  _ __ _____  ___ 		Endpoint: %v	
| '_ \| '__/ _ \ \/ / |		Api docs: http://%v/swagger/index.html
| |_) | | | (_) >  <| |		Version:  %v	
| .__/|_|  \___/_/\_\_|		
| |                    
|_|             
`, Addr, Addr, Version)
	color.HiGreen(banner)

	//apiBanner := fmt.Sprintf("\nListening and serving HTTP on %v", Addr)
	//color.HiGreen("%v\n", apiBanner)
	//fmt.Printf("Visit http://%v/swagger/index.html for api docs.", Addr)
	//color.HiGreen(docsBanner)
	fmt.Print("\n")
	log.Printf("Listening and serving HTTP on %v", Addr)
	fmt.Print("\n\n\n")

}
