package database

import (
	"fmt"
	"github.com/vesoft-inc/nebula-go/v3/nebula"
	"github.com/xueqianLu/routegen/config"
	"github.com/xueqianLu/routegen/database/models"
	"github.com/xueqianLu/routegen/log"
	"github.com/zhihu/norm"
	"github.com/zhihu/norm/constants"
	"github.com/zhihu/norm/dialectors"
	"time"
)

func NewDb(conf *config.Config) *norm.DB {
	dalector := dialectors.MustNewNebulaDialector(dialectors.DialectorConfig{
		Addresses: []string{conf.DbHost},
		Timeout:   time.Second * 5,
		Space:     conf.DbSpace,
		Username:  conf.DbUser,
		Password:  conf.DbPasswd,
	})
	db := norm.MustOpen(dalector, norm.Config{})
	return db
}

func InsertToken(db *norm.DB, name string, address string) error {
	token := &models.Token{
		Name:    name,
		Address: address,
	}
	err := db.InsertVertex(token)
	if err != nil {
		log.WithField("err", err).Error("insert token failed")
	}
	return err
}

func InsertPair(db *norm.DB, dexname string, pairaddr string, token0, token1 string) error {
	pair := &models.Pair{
		EModel: norm.EModel{
			Src:       token0,
			SrcPolicy: constants.PolicyNothing,
			Dst:       token1,
			DstPolicy: constants.PolicyNothing,
		},
		DexName:     dexname,
		Token1:      token1,
		Token0:      token0,
		PairAddress: pairaddr,
	}
	err := db.InsertEdge(pair)
	if err != nil {
		log.WithField("err", err).Error("insert pair failed")
	}
	return err
}

func getValueofValue(value *nebula.Value) string {
	if value.NVal != nil {
		return fmt.Sprintf("value.NVal=%v", value.NVal)
	}
	if value.BVal != nil {
		return fmt.Sprintf("value.BVal=%v", value.BVal)
	}
	if value.IVal != nil {
		return fmt.Sprintf("value.IVal=%v", value.IVal)
	}
	if value.FVal != nil {
		return fmt.Sprintf("value.FVal=%v", value.FVal)
	}
	if value.SVal != nil {
		return fmt.Sprintf("%v", string(value.SVal))
	}
	if value.DVal != nil {
		return fmt.Sprintf("value.DVal=%v", value.DVal)
	}
	if value.TVal != nil {
		return fmt.Sprintf("value.TVal=%v", value.TVal)
	}
	if value.DtVal != nil {
		return fmt.Sprintf("value.DtVal=%v", value.DtVal)
	}
	if value.VVal != nil {
		return fmt.Sprintf("value.VVal=%v", value.VVal)
	}
	if value.EVal != nil {
		return fmt.Sprintf("value.EVal=%v", value.EVal)
	}
	if value.PVal != nil {
		return fmt.Sprintf("value.PVal=%v", value.PVal)
	}
	if value.LVal != nil {
		return fmt.Sprintf("value.LVal=%v", value.LVal)
	}
	if value.MVal != nil {
		return fmt.Sprintf("value.MVal=%v", value.MVal)
	}
	if value.UVal != nil {
		return fmt.Sprintf("value.UVal=%v", value.UVal)
	}
	if value.GVal != nil {
		return fmt.Sprintf("value.GVal=%v", value.GVal)
	}
	if value.GgVal != nil {
		return fmt.Sprintf("value.GgVal=%v", value.GgVal)
	}
	if value.DuVal != nil {
		return fmt.Sprintf("value.DuVal=%v", value.DuVal)
	}
	return ""
}

func getValueofTag(tag *nebula.Tag) string {
	v := ""
	s := fmt.Sprintf("name:%s", string(tag.Name))
	v += s
	for k, p := range tag.Props {
		s = fmt.Sprintf("[%s]=%s\n", k, getValueofValue(p))
		v += s
	}
	return v
}

func getValueofTags(tags []*nebula.Tag) string {
	v := ""
	for i, tag := range tags {
		s := fmt.Sprintf("t[%d]=%s\n", i, getValueofTag(tag))
		v += s
	}
	return v
}

func printVertex(v *nebula.Vertex) {
	log.Infof("vertex.VID = %s", getValueofValue(v.Vid))
	if len(v.Tags) > 0 {
		log.Infof("vertex.Tags = %s", getValueofTags(v.Tags))
	}
}

func printStep(index int, s *nebula.Step) {
	//Dst *Vertex `thrift:"dst,1" db:"dst" json:"dst"`
	//Type EdgeType `thrift:"type,2" db:"type" json:"type"`
	//Name []byte `thrift:"name,3" db:"name" json:"name"`
	//Ranking EdgeRanking `thrift:"ranking,4" db:"ranking" json:"ranking"`
	//Props map[string]*Value `thrift:"props,5" db:"props" json:"props"`
	log.Infof("step[%d].dst = %s", index, getValueofValue(s.Dst.Vid))
	//printVertex(s.Dst)
	//log.Infof("step[%d].edgetype = %d", index, s.Type)
	//log.Infof("step[%d].rank = %d", index, s.Ranking)
	//log.Infof("step[%d].name = %s", index, string(s.Name))
	for k, v := range s.Props {
		log.Infof("step[%d].prop[%s]=%v", index, k, getValueofValue(v))
	}
	//log.Infof("step[%d].props = %v", index, s.Props)

}

func GetDstFromStep(step *nebula.Step) string {
	return getValueofValue(step.Dst.GetVid())
}

func GetPairInfoFromStep(step *nebula.Step, routeStep *RouteStep) {
	if dex, exist := step.Props[PairProp_dex]; exist {
		routeStep.Dex = getValueofValue(dex)
	}
	if pairAddr, exist := step.Props[PairProp_paircontract]; exist {
		routeStep.Pair = getValueofValue(pairAddr)
	}
	if fee, exist := step.Props[PairProp_fee]; exist {
		routeStep.Fee = getValueofValue(fee)
	}
}

type RouteStep struct {
	Pair string `json:"pair"`
	Dex  string `json:"dex"`
	Src  string `json:"from"`
	Dst  string `json:"to"`
	Fee  string `json:"fee"`
}

type TokenRoute struct {
	Steps []RouteStep `json:"steps"`
}

func ParsePathInfo(path *nebula.Path) []RouteStep {
	src := path.GetSrc()
	steps := path.GetSteps()
	routePath := make([]RouteStep, len(steps))
	srcToken := getValueofValue(src.Vid)
	for i, step := range steps {
		routeStep := RouteStep{
			Src: srcToken,
		}
		routeStep.Dst = GetDstFromStep(step)
		GetPairInfoFromStep(step, &routeStep)
		srcToken = routeStep.Dst
		routePath[i] = routeStep
	}
	return routePath
}

func QueryRoute(db *norm.DB, token0, token1 string) []*TokenRoute {
	nql := fmt.Sprintf("FIND NOLOOP PATH WITH PROP FROM \"%s\" TO \"%s\" OVER * YIELD path AS p", token0, token1)
	result := make([]map[string]interface{}, 0)
	res, err := db.Debug().Execute(nql)
	if err != nil {
		log.WithField("err", err).Error("query route failed")
		return []*TokenRoute{}
	} else {
		//log.WithField("rows", len(res.GetRows())).Info("query route")
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			log.WithField("err", err).Error("parse route failed")
			return []*TokenRoute{}
		}
		paths := make([]*TokenRoute, 0, len(result))

		for _, vpath := range result {
			// vpath only have one key (AS p)
			for _, v := range vpath {
				if path, ok := v.(*nebula.Path); ok {
					steps := ParsePathInfo(path)
					tokenRoute := new(TokenRoute)
					tokenRoute.Steps = steps
					paths = append(paths, tokenRoute)
				}
			}
		}
		return paths
	}
}
