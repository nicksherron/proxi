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

// getCmd represents the stats command
var (
	numProxies int
	anon       bool
	country    string
	getAll     bool
	getCmd     = &cobra.Command{
		Use:   "get",
		Short: "Return one or more proxies from db that passed checks.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().Parse(args)
			getProxy()
		},
	}
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringVarP(&address, "url", "u", fmt.Sprintf("http://%v", listenAddr()), "Url of running ProxyPool server.")
	getCmd.PersistentFlags().IntVarP(&numProxies, "num", "n", 1, "Number of proxies to return.")
	getCmd.PersistentFlags().BoolVar(&anon, "anon", false, "Only return anonymous proxies.")
	getCmd.PersistentFlags().StringVarP(&country, "country", "c", "", "Filter by country. Format is 'US', 'CH' etc.")
	getCmd.PersistentFlags().BoolVar(&getAll, "all", false, "Return all proxies ignoring filters or status. Warning! may produce lots of results.")

}
