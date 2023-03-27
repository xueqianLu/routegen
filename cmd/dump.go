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
	"github.com/xueqianLu/routegen/tool"
	"github.com/xueqianLu/routegen/types"
	"github.com/zhihu/norm"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const (
	outputFlag  = "out"
	maxOpFlag   = "op"
	routineFlag = "routine"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump all route for given token pairs",
	Run: func(cmd *cobra.Command, args []string) {
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
		routine, _ := cmd.PersistentFlags().GetUint(routineFlag)

		if err := DumpHandler(routine, tokenList, output, op); err != nil {
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
	dumpCmd.PersistentFlags().Uint(routineFlag, 5, "routine count to dump route file")
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

func DumpHandler(routine uint, tokens []string, dumpfile string, maxOp int) error {
	worker := NewWorker(routine)
	worker.Start()
	return worker.DumpRouteToFile(dumpfile, tokens, maxOp)
}

type Worker struct {
	task   *tool.Tasks
	dbpool []*norm.DB
}

func NewWorker(rountines uint) *Worker {
	w := new(Worker)
	task := tool.NewTasks(rountines, w.handler)
	w.task = task
	w.dbpool = make([]*norm.DB, int(rountines))
	for i := 0; i < int(rountines); i++ {
		w.dbpool[i] = database.NewDb(config.GetConfig())
	}
	return w
}

func (w *Worker) handler(t interface{}) {
	item := t.(Item)
	db := w.dbpool[item.index%len(w.dbpool)]
	paths := database.QueryRouteWithMaxJump(db, item.token0, item.token1, item.maxOp)
	log.Infof("got token path %d", len(paths))
	data := convertPathToString(paths)
	for _, str := range data {
		item.response <- str
	}
}

func (w *Worker) Start() {
	w.task.Run()
}

type Item struct {
	token0, token1 string
	maxOp          int
	index          int
	response       chan interface{}
}

func (w *Worker) DumpRouteToFile(dumpfile string, tokens []string, maxOp int) error {
	fp, err := os.OpenFile(dumpfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend|os.ModePerm) // 读写方式打开
	if err != nil {
		log.WithField("err", err).WithField("file", dumpfile).Error("open file failed")
		return err
	}
	defer fp.Close()
	log.Infof("total token %d", len(tokens))

	results := make(chan string, 10000000)
	writeFinished := false

	go func() {
		var count = 0
		for {
			select {
			case s, ok := <-results:
				if !ok {
					writeFinished = true
					log.Infof("total write to file count %d", count)
					return
				}
				_, err = fp.WriteString(s)
				count += 1
				if (count % 20) == 0 {
					log.Infof("write to file count %d", count)
					fp.Sync()
				}
				//log.Infof("consume routine write to file")
				if err != nil {
					log.WithField("err", err).Error("write to file failed")
					return
				}
			}
		}
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < len(tokens); i++ {
		for j := 0; j < len(tokens); j++ {
			if i == j {
				continue
			}
			wg.Add(1)
			go func(token0, token1 string, op int) {
				defer wg.Done()
				res := make(chan interface{})
				item := Item{
					response: res,
					token0:   token0,
					token1:   token1,
					maxOp:    op,
				}
				//log.Debugf("add item to task")
				if e := w.task.AddTask(item); e != nil {
					err = e
				} else {
					data := <-res
					switch msg := (data).(type) {
					case error:
						err = msg
					case string:
						results <- msg
					}
				}
			}(tokens[i], tokens[j], maxOp)
		}
	}
	wg.Wait()
	close(results)
	for !writeFinished {
		log.Info("wait write file finish")
		time.Sleep(time.Second)
	}

	return err
}
