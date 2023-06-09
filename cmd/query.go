/*
Copyright © 2023 xueqianLu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xueqianLu/routegen/config"
	"github.com/xueqianLu/routegen/database"
	"github.com/xueqianLu/routegen/log"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Error("please enter token0 and token1")
			return
		}
		token0, token1 := args[0], args[1]
		db := database.NewDb(config.GetConfig())
		paths := database.QueryRoute(db, token0, token1)
		for i, path := range paths {
			route := fmt.Sprintf("path[%d]=", i)
			for n, step := range path.Steps {
				if n == 0 {
					str := fmt.Sprintf("%s ---(%s:%s:%s)---> %s", step.Src, step.Pairs[0].Dex, step.Pairs[0].Pair, step.Pairs[0].Fee, step.Dst)
					route += str
				} else {
					str := fmt.Sprintf(" ---(%s:%s:%s)---> %s", step.Pairs[0].Dex, step.Pairs[0].Pair, step.Pairs[0].Fee, step.Dst)
					route += str
				}
			}
			log.Info(route)
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
