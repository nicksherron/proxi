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
	"fmt"

	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var (
	watch int

	statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Check server stats",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
			stats()
		},
	}
)

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.PersistentFlags().StringVarP(&address, "url", "u", fmt.Sprintf("http://%v", listenAddr()), "Url of running ProxyPool server.")
	statsCmd.PersistentFlags().IntVarP(&watch, "watch", "w", 0, "Time in seconds to refresh. Set to 0 if you don't to watch.")
}
