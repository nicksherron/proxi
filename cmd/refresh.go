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

// refreshCmd represents the stats command
var (
	refreshCmd = &cobra.Command{
		Use:   "refresh",
		Short: "Re-download and check proxies.",
		Long:  "Re-download and check proxies if the server is not already busying downloading or checking. Returns busy, if so.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
			getRefresh()
		},
	}
)

func init() {
	rootCmd.AddCommand(refreshCmd)
	refreshCmd.PersistentFlags().StringVarP(&address, "url", "u", fmt.Sprintf("http://%v", listenAddr()), "Url of running ProxyPool server.")
}
