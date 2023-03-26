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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xueqianLu/routegen/cmd/utils"
	"github.com/xueqianLu/routegen/config"
	"github.com/xueqianLu/routegen/database"
	"github.com/xueqianLu/routegen/log"
	"github.com/xueqianLu/routegen/types"
	"github.com/zhihu/norm"
	"io/ioutil"
	"os"
)

const (
	outputFlag = "out"
	maxOpFlag  = "op"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump all route for given token pairs",
	Run: func(cmd *cobra.Command, args []string) {
		db := database.NewDb(config.GetConfig())
		tokenMap := make(map[string]bool)
		tokenList := make([]string, 0)
		for _, datafile := range args {
			if utils.Exists(datafile) {
				log.Info("import from file ", datafile)
			} else {
				log.Errorf("file (%s) not exist", datafile)
				continue
			}

			data, err := ioutil.ReadFile(datafile)
			if err != nil {
				log.WithField("err", err).Error("read data file failed")
				continue
			}
			var dexInfo = new(ImportData)
			err = json.Unmarshal(data, &dexInfo)
			if err != nil {
				log.WithField("err", err).Error("unmarshal file failed")
				continue
			}
			for _, pair := range dexInfo.Data.Pairs {
				tokenMap[pair.Token0.Address] = true
				tokenMap[pair.Token1.Address] = true
			}
		}
		for token, _ := range tokenMap {
			tokenList = append(tokenList, token)
		}
		output, _ := cmd.PersistentFlags().GetString(outputFlag)
		op, _ := cmd.PersistentFlags().GetInt(maxOpFlag)

		if err := DumpHandler(db, tokenList, output, op); err != nil {
			log.Errorf("dump token route failed with err:(%s)", err)
		} else {
			log.Info("dump token route finished")
		}
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.PersistentFlags().String(outputFlag, "dump.txt", "out put filename")
	dumpCmd.PersistentFlags().Int(maxOpFlag, 4, "max jump for token swap route")
}

func convertPathToString(routes []*types.TokenRoute) []string {
	paths := make([]string, 0)
	for _, route := range routes {
		pair := "["
		tokens := ""
		// ["","",""]
		for i, step := range route.Steps {
			if i == 0 {
				pair += step.Pair
			} else {
				pair += ","
				pair += step.Pair
			}
			t := fmt.Sprintf("[\"%s\",\"%s\"],", step.Src, step.Dst)
			tokens += t
		}
		tokens = tokens[:len(tokens)-1]
		pair += "]"
		path := fmt.Sprintf("[%s,%s]\n", pair, tokens)
		paths = append(paths, path)
	}
	return paths

}

func DumpHandler(db *norm.DB, tokens []string, dumpfile string, maxOp int) error {
	fp, err := os.OpenFile(dumpfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm) // 读写方式打开
	if err != nil {
		log.WithField("err", err).WithField("file", dumpfile).Error("open file failed")
		return err
	}
	defer fp.Close()
	log.Infof("dump token route, token count %d", len(tokens))

	for i := 0; i < len(tokens); i++ {
		for j := 0; j < len(tokens); j++ {
			if i == j {
				continue
			}
			paths := database.QueryRouteWithMaxJump(db, tokens[i], tokens[j], maxOp)
			data := convertPathToString(paths)
			for _, str := range data {
				_, err = fp.WriteString(str)
				if err != nil {
					log.WithField("err", err).Error("write to file failed")
					return err
				}

			}
		}
	}
	return nil
}
