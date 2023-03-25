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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xueqianLu/routegen/cmd/utils"
	"github.com/xueqianLu/routegen/config"
	"github.com/xueqianLu/routegen/contracts"
	"github.com/xueqianLu/routegen/database"
	"github.com/xueqianLu/routegen/log"
	"github.com/xueqianLu/routegen/types"
	"github.com/zhihu/norm"
	"io/ioutil"

	"github.com/spf13/cobra"
)

const (
	urlFlag = "url"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("please enter import file")
			return
		}
		datafile := args[0]
		if utils.Exists(datafile) {
			log.Info("import from file ", datafile)
		} else {
			log.Errorf("file (%s) not exist", datafile)
			return
		}
		url, _ := cmd.PersistentFlags().GetString(urlFlag)
		if err := ImportHandler(datafile, url); err != nil {
			log.Error("import data failed")
		} else {
			log.Info("import finished")
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().String(urlFlag, "https://rpc.ankr.com/bsc", "rpc url")
}

func prepare(db *norm.DB) error {
	createSchema := "" +
		"CREATE TAG IF NOT EXISTS token(name string, address string);" +
		"CREATE EDGE IF NOT EXISTS pair(dex string, tracked string, paircontract string, token0 string, token1 string);"
	_, err := db.Execute(createSchema)
	return err
}

func ImportHandler(datafile string, url string) error {
	db := database.NewDb(config.GetConfig())
	if err := prepare(db); err != nil {
		log.WithField("err", err).Fatalf("prepare db failed")
		panic(err)
	}
	//return nil
	client, err := ethclient.Dial(url)
	if err != nil {
		log.WithField("err", err).Fatalf("dial rpc failed")
		return err
	}
	data, err := ioutil.ReadFile(datafile)
	if err != nil {
		log.WithField("err", err).Fatalf("read data file failed")
		return err
	}
	var dexlist = make([]*types.DexData, 0)
	err = json.Unmarshal(data, &dexlist)
	if err != nil {
		log.WithField("err", err).Fatalf("read data file failed")
		return err
	}
	for _, dex := range dexlist {
		for _, pair := range dex.Pairs {
			name0 := contracts.GetTokenName(client, pair.Token0)
			name1 := contracts.GetTokenName(client, pair.Token1)
			_ = database.InsertToken(db, name0, pair.Token0)
			_ = database.InsertToken(db, name1, pair.Token1)
			_ = database.InsertPair(db, dex.Name, pair.Address, pair.Token0, pair.Token1)
		}
	}
	db.Close()
	return nil
}
