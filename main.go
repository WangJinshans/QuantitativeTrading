package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math"
	"net/http"
	"quant_trade/model"
	"regexp"
	"strings"
	"time"
	"xorm.io/core"
)

var client http.Client
var engine *xorm.Engine

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

func Split808(segment []byte) (messages [][]byte, residueBytes []byte, invalidMessages [][]byte) {
	//segment := []byte{0x77, 0x01, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x7e, 0x02,0x7e, 0x02, 0x03, 0x05}
	//segment := []byte{0x77, 0x01, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07}
	old7e := []byte{0x7e, 0x02}
	source7e := []byte{0x7e}
	old7d := []byte{0x7e, 0x01}
	source7d := []byte{0x7d}
	rawPackages := bytes.Split(segment, old7e)
	last := len(rawPackages)
	for index, pkg := range rawPackages {
		fmt.Printf("pkg is: %x\n", pkg)
		if bytes.Equal(pkg, []byte("")) {
			continue
		}
		data := bytes.Replace(pkg, old7e, source7e, -1)
		data = bytes.Replace(pkg, old7d, source7d, -1)
		if index == 0 {
			fmt.Printf("invalid is: %x\n", data)
			var rawData []byte
			rawData = append(rawData, 0x7e, 0x02)
			rawData = append(rawData, pkg...)
			invalidMessages = append(invalidMessages, rawData)
			continue
		}
		if index == last-1 {
			fmt.Printf("residueBytes is: %x\n", data)
			residueBytes = append(residueBytes, 0x7e, 0x02)
			residueBytes = append(residueBytes, pkg...)
			continue
		}
		messages = append(messages, data)
		log.Info().Msgf("pkg is: %x", data)
	}
	return
}

func Split808Fix(segment []byte) (messages [][]byte, residueBytes []byte, invalidMessages [][]byte) {
	startFlag := []byte{0x7e, 0x02}
	var index int
	var indexList []int

	for i := 0; i < len(segment)-2; i++ {
		sf := segment[i : i+2]
		if bytes.Equal(sf, startFlag) {
			indexList = append(indexList, index)
			index += 2
		}
		segment = segment[1:]
		index += 1
	}

	fmt.Printf("index list is: %v\n", indexList)
	//messages, residueBytes, invalidMessages = SplitPackage(segment, indexList)
	return
}

func SplitPackage(segment []byte, indexList []int) (messages [][]byte, residueBytes []byte, invalidMessages [][]byte) {
	//segment := []byte{0x77, 0x01, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x7e, 0x02,0x7e, 0x02, 0x03, 0x05}
	//segment := []byte{0x7e, 0x02, 0x02, 0x03, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x7e, 0x02,0x7e, 0x02, 0x03, 0x05}
	old7e := []byte{0x7e, 0x02}
	source7e := []byte{0x7e}
	old7d := []byte{0x7e, 0x01}
	source7d := []byte{0x7d}
	var entireList []int
	if len(indexList)%2 != 0 {
		// 有剩余
		left := indexList[len(indexList)-1]
		entireList = indexList[:len(indexList)-1]
		residueBytes = append(residueBytes, segment[left:]...)
	} else {
		// 包完整
		entireList = indexList
	}

	for index := 0; index < len(entireList)-1; index++ {
		if index == 0 {
			// 前面的数据写入无效数据buffer
			data := segment[:entireList[index]]
			invalidMessages = append(invalidMessages, data)
		}
		if index%2 == 0 {
			// 起始 -- 结束 中间数据为完整数据包
			data := segment[entireList[index]:entireList[index+1]]
			pkg := bytes.Replace(data, old7e, source7e, -1)
			data = bytes.Replace(data, old7d, source7d, -1)
			messages = append(messages, pkg)
		} else {
			// 结束 -- 起始  中间的数据写入无效buffer
			data := segment[entireList[index]+2 : entireList[index+1]]
			if data != nil {
				invalidMessages = append(invalidMessages, data)
			}
		}
	}
	return
}

func main() {
	//segment := []byte{0x77, 0x01, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x7e, 0x02, 0x7e, 0x02, 0x03, 0x05}
	segment := []byte{0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02}
	segment = []byte{0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x7e, 0x02}
	segment = []byte{0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x7e, 0x02, 0x7e, 0x02, 0x03, 0x05}
	segment = []byte{0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x77, 0x7e, 0x02, 0x02, 0x7e, 0x02}
	//segment = []byte{0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07, 0x7e, 0x02, 0x77, 0x7e, 0x02, 0x02, 0x7e, 0x02, 0x77}

	//segment := []byte{0x77, 0x01, 0x7e, 0x02,   0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07}
	//messages, residueBytes, invalidMessages := Split808(segment)
	messages, residueBytes, invalidMessages := Split808Fix(segment)
	fmt.Printf("message is: %x\n", messages)
	fmt.Printf("residueBytes is: %x\n", residueBytes)
	fmt.Printf("invalidMessages is: %x\n", invalidMessages)
	fmt.Println("-------------------------------------------------------------------------------------------------")
	//segment = []byte{0x77, 0x01, 0x7e, 0x02, 0x7e, 0x02, 0x02, 0x03, 0x05, 0x06, 0x07}
	//if residueBytes != nil {
	//	segment = append(residueBytes, segment...)
	//}
	//messages, residueBytes, invalidMessages = Split808(segment)
	//fmt.Printf("message is: %x\n", messages)
	//fmt.Printf("residueBytes is: %x\n", residueBytes)
	//fmt.Printf("invalidMessages is: %x\n", invalidMessages)

	//infoList := GetStockList()
	//log.Info().Msg("------------------------------------------------------------------------------------------------------------------")
	//FilterData(infoList)
	//GetBkInfo()
}

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

func GetStockList() (stockList []*model.StockInfo) {
	var page = 1
	var total int
	for {
		time.Sleep(500 * time.Millisecond)

		webUrl := fmt.Sprintf("http://51.push2.eastmoney.com/api/qt/clist/get?cb=jQuery1124012243073358859502_1634284246158&pn=%d&pz=30&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152", page)
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
			code, ok := info["f12"]
			if !ok {
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

func FilterData(stockList []*model.StockInfo) {

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
		if info.CurrentRate > 10.5 || info.HighestRate > 10.5 || math.Abs(info.LowestRate) > 10.5 {
			// 筛出部分创业板的数据
			continue
		}

		if info.CurrentPrice > 50 || info.CurrentPrice < 5 {
			// 价格相对较高或低
			continue
		}
		if info.CurrentRate < 3 || info.CurrentRate > 5 {
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
