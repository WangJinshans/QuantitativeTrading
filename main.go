package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"quant_trade/model"
	"regexp"
	"strings"
	"time"
	"xorm.io/core"
)

var client http.Client
var engine *xorm.Engine

var headers = map[string]string{
	"Accept":           "*/*",
	"Accept-Language":  "zh-CN,zh;q=0.9",
	"Cache-Control":    "no-cache",
	"Host":             "push2.eastmoney.com",
	"Pragma":           "no-cache",
	"Proxy-Connection": "keep-alive",
	"Referer":          "http://data.eastmoney.com/",
	"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Safari/537.36",
}

func init() {
	//var err error
	//engine, err = xorm.NewEngine("mysql", "root:wangjinshan@tcp(127.0.0.1:3306)/trade?charset=utf8")
	//if err != nil {
	//	log.Error().Msgf("connect to database error: %v", err)
	//	return
	//}
	//err = engine.Sync2(new(model.StockInfo))
	//if err != nil {
	//	log.Error().Msgf("sync2 table error: %v", err)
	//}
}

func GetStockList() (stockList []*model.StockInfo) {
	var page = 1
	var total int
	for {
		time.Sleep(500 * time.Millisecond)

		webUrl := fmt.Sprintf("http://51.push2.eastmoney.com/api/qt/clist/get?cb=jQuery1124012243073358859502_1634284246158&pn=%d&pz=30&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f57,f62,f128,f136,f115,f152", page)
		req, err := http.NewRequest(http.MethodGet, webUrl, nil)
		if err != nil {
			log.Error().Msgf("create req error: %v", err)
			break
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Error().Msgf("failed to request: %v", err)
			break
		}

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("get response error: %v", err)
			break
		}
		resp.Body.Close()

		//log.Info().Msgf("content is: %s", string(content))
		reg := regexp.MustCompile("jQuery1124012243073358859502_1634284246158\\(([\\s\\S]+?)\\);")
		rs := reg.FindAllSubmatch(content, -1)
		//log.Info().Msgf("rs is: %s", string(rs[0][1]))
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(rs[0][1], &dataMap)
		if err != nil {
			log.Error().Msgf("unmarshal error: %v", err)
			return
		}
		diffItem, ok := dataMap["data"]
		if !ok || diffItem == nil {
			log.Error().Msg("empty recode...")
			break
		}
		data, ok := diffItem.(map[string]interface{})
		if !ok {
			continue
		}
		d, ok := data["diff"]
		if !ok {
			continue
		}
		diffData := d.([]interface{})
		for _, v := range diffData {
			log.Info().Msgf("v is: %v", v)
			info := v.(map[string]interface{})
			codeStr, ok := info["f12"]
			if !ok {
				continue
			}
			code := codeStr.(string)
			if strings.HasPrefix(code, "300") || strings.HasPrefix(code, "688") {
				continue
			}
			name, ok := info["f14"].(string)
			if !ok {
				continue
			}
			currentPrice, ok := info["f2"].(float64)
			if !ok {
				continue
			}
			highestPrice, ok := info["f15"].(float64)
			if !ok {
				continue
			}
			lowestPrice, ok := info["f16"].(float64)
			if !ok {
				continue
			}
			startPrice, ok := info["f17"].(float64)
			if !ok {
				continue
			}
			yesterdayPrice, ok := info["f18"].(float64)
			if !ok {
				continue
			}
			rate, ok := info["f3"].(float64)
			if !ok {
				continue
			}
			amplitude, ok := info["f7"].(float64)
			if !ok {
				continue
			}
			changeRate, ok := info["f8"].(float64)
			if !ok {
				continue
			}
			highestRate := ((highestPrice - yesterdayPrice) / yesterdayPrice) * 100
			lowestRate := ((lowestPrice - yesterdayPrice) / yesterdayPrice) * 100
			diff := rate - highestRate
			log.Info().Msgf("code is: %v, currentPrice is: %v", code, currentPrice)
			log.Info().Msgf("highestPrice is: %v, lowestPrice is: %v, startPrice is: %v", highestPrice, lowestPrice, startPrice)
			log.Info().Msgf("rate is: %v, amplitude is: %v, diff is: %v", rate, amplitude, diff)
			log.Info().Msgf("changeRate is: %v, highestRate is: %v, lowestRate is: %v", changeRate, highestRate, lowestRate)

			stockInfo := &model.StockInfo{
				StockId:      fmt.Sprintf("%s", code),
				StockName:    name,
				CurrentPrice: currentPrice,
				HighestPrice: highestPrice,
				LowestPrice:  lowestPrice,
				HighestRate:  highestRate,
				LowestRate:   lowestRate,
				CurrentRate:  rate,
				ChangeRate:   changeRate,
				Amplitude:    amplitude,
				Diff:         diff,
			}
			stockList = append(stockList, stockInfo)
		}
		total += len(diffData)
		page += 1
	}
	log.Info().Msgf("total is: %d", total)
	return
}

func GetBkInfo() {
	webUrl := "http://push2.eastmoney.com/api/qt/clist/get?cb=jQuery1123022593397568009088_1634282352811&pn=1&pz=500&po=1&np=1&fields=f12%2Cf13%2Cf14%2Cf62&fid=f62&fs=m%3A90%2Bt%3A2"
	req, err := http.NewRequest(http.MethodGet, webUrl, nil)
	if err != nil {
		log.Error().Msgf("create req error: %v", err)
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msgf("failed to request: %v", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("get response error: %v", err)
		return
	}
	fmt.Printf("content is: %v\n", string(content))

	reg := regexp.MustCompile("jQuery1123022593397568009088_1634282352811\\(([\\s\\S]+?)\\);")
	rs := reg.FindAllSubmatch(content, -1)
	//fmt.Printf("rs is: %s\n", string(rs[0][1]))
	dataMap := make(map[string]interface{})
	err = json.Unmarshal(rs[0][1], &dataMap)
	if err != nil {
		log.Error().Msgf("unmarshal error: %v", err)
		return
	}
	diffItem, ok := dataMap["data"]
	if !ok {
		log.Error().Msg("empty recode...")
		return
	}
	data := diffItem.(map[string]interface{})
	diffData := data["diff"].([]interface{})
	//log.Info().Msgf("data map is: %v", dataMap)
	for _, v := range diffData {
		log.Info().Msgf("v is: %v", v)
	}
}

func getStockDetailInfo(stockId string) {

}

func FilterData(stockList []*model.StockInfo) {

	var codeList []string
	log.Info().Msgf("len is: %d", len(stockList))
	if len(stockList) <= 0 {
		return
	}
	for _, info := range stockList {
		if strings.Contains(info.StockName, "ST") {
			continue
		}
		if info.ChangeRate < 5 || info.ChangeRate > 10 {
			continue
		}

		if info.CurrentPrice > 50 || info.CurrentPrice < 5 {
			// 价格相对较高或低
			continue
		}
		if info.CurrentRate < 2 || info.CurrentRate > 5 {
			// 涨幅适中
			continue
		}
		diff := info.CurrentRate - info.HighestRate
		log.Info().Msg("-----------------------------------------------------------------")
		log.Info().Msgf("name is: %s, code is: %v, currentPrice is: %v", info.StockName, info.StockId, info.CurrentPrice)
		log.Info().Msgf("highestPrice is: %v, lowestPrice is: %v", info.HighestPrice, info.LowestPrice)
		log.Info().Msgf("rate is: %v, amplitude is: %v, diff is: %v", info.CurrentRate, info.Amplitude, diff)
		log.Info().Msgf("changeRate is: %v, highestRate is: %v, lowestRate is: %v", info.ChangeRate, info.HighestRate, info.LowestRate)
		//SaveData(info)
		codeList = append(codeList, info.StockId)
	}
}

// 获取历史数据
func FilterHistoryData(stockList []*model.StockInfo) {

	if len(stockList) <= 0 {
		return
	}
	//for _, info := range stockList {
	//
	//
	//}
}

// 存储到mysql中
func SaveData(stockInfo *model.StockInfo) {
	engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	engine.SetTableMapper(core.SnakeMapper{})
	affected, err := engine.Insert(stockInfo)
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("affected is: %v", affected)
}

func main() {

	infoList := GetStockList()
	log.Info().Msg("------------------------------------------------------------------------------------------------------------------")
	FilterData(infoList)
	//GetBkInfo()
}
