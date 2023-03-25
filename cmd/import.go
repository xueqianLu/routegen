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
	"github.com/zhihu/norm"
	"io/ioutil"

	"github.com/spf13/cobra"
)

const (
	urlFlag    = "url"
	initDBFlag = "initdb"
)

type ImportToken struct {
	Address string `json:"id"`
	Name    string `json:"name"`
}

type ImportPairInfo struct {
	Address      string      `json:"id"`
	Name         string      `json:"name"`
	TrackedValue string      `json:"trackedReserveBNB"`
	Token0       ImportToken `json:"token0"`
	Token1       ImportToken `json:"token1"`
}

type ImportPairs struct {
	Pairs []ImportPairInfo `json:"pairs"`
}

type ImportData struct {
	Name string      `json:"name"`
	Fee  string      `json:"fee"`
	Data ImportPairs `json:"data"`
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Error("please enter import file")
			return
		}
		url, _ := cmd.PersistentFlags().GetString(urlFlag)
		initdb, _ := cmd.PersistentFlags().GetBool(initDBFlag)

		db := database.NewDb(config.GetConfig())
		if initdb {
			if err := prepare(db); err != nil {
				log.WithField("err", err).Fatalf("prepare db failed")
				panic(err)
			}
			log.Infof("init db finished")
		}

		for _, datafile := range args {
			if utils.Exists(datafile) {
				log.Info("import from file ", datafile)
			} else {
				log.Errorf("file (%s) not exist", datafile)
				continue
			}
			if err := ImportHandler(db, datafile, url); err != nil {
				log.Error("import data from %s failed", datafile)
			} else {
				log.Infof("import data from %s finished", datafile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().String(urlFlag, "https://rpc.ankr.com/bsc", "rpc url")
	importCmd.PersistentFlags().Bool(initDBFlag, false, "init database")
}

func prepare(db *norm.DB) error {
	createSchema := "" +
		"CREATE TAG IF NOT EXISTS token(name string, address string);" +
		"CREATE EDGE IF NOT EXISTS pair(dex string, tracked string, fee string, pairaddress string, token0 string, token1 string);" +
		"CREATE TAG INDEX token_index on token();" +
		"CREATE EDGE INDEX pair_index on pair();"
	_, err := db.Execute(createSchema)
	return err
}

func ImportHandler(db *norm.DB, datafile string, url string) error {

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
	var dexInfo = new(ImportData)
	err = json.Unmarshal(data, &dexInfo)
	if err != nil {
		//log.WithField("err", err).Fatalf("read data file failed")
		return err
	}
	dexName := dexInfo.Name
	for _, pair := range dexInfo.Data.Pairs {
		var name0, name1 = pair.Token0.Name, pair.Token1.Name
		if len(name0) == 0 {
			name0 = contracts.GetTokenName(client, pair.Token0.Address)
		}
		if len(name1) == 0 {
			name1 = contracts.GetTokenName(client, pair.Token1.Address)
		}

		_ = database.InsertToken(db, name0, pair.Token0.Address)
		_ = database.InsertToken(db, name1, pair.Token1.Address)
		_ = database.InsertPair(db, dexName, pair.Address, dexInfo.Fee, pair.TrackedValue, pair.Token0.Address, pair.Token1.Address)
	}
	//for _, dex := range dexlist {
	//	for _, pair := range dex.Pairs {
	//		name0 := contracts.GetTokenName(client, pair.Token0)
	//		name1 := contracts.GetTokenName(client, pair.Token1)
	//		_ = database.InsertToken(db, name0, pair.Token0)
	//		_ = database.InsertToken(db, name1, pair.Token1)
	//		_ = database.InsertPair(db, dex.Name, pair.Address, pair.Token0, pair.Token1)
	//	}
	//}
	db.Close()
	return nil
}
