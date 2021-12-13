package data_center

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math"
	"net/http"
	"quant_trade/db"
	"quant_trade/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var client http.Client

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

func GetChanges() (stockChangeMapInfo map[string]map[string]interface{}) {
	stockChangeMapInfo = make(map[string]map[string]interface{})
	var pageIndex int = 0

	for {
		webUrl := fmt.Sprintf("http://push2ex.eastmoney.com/getAllStockChanges?type=8201,8193,64&cb=jQuery112401745104453711004_1639316373215&pageindex=%d&pagesize=64&ut=7eea3edcaed734bea9cbfc24409ed989&dpt=wzchanges", pageIndex)
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

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("get response error: %v", err)
			resp.Body.Close()
			return
		}
		resp.Body.Close()

		reg := regexp.MustCompile("jQuery112401745104453711004_1639316373215\\(([\\s\\S]+?)\\);")
		rs := reg.FindAllSubmatch(content, -1)
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(rs[0][1], &dataMap)
		if err != nil {
			log.Error().Msgf("unmarshal error: %v", err)
			return
		}
		data := dataMap["data"]
		allStock := data.(map[string]interface{})["allstock"]

		stockList := allStock.([]interface{})

		for _, info := range stockList {

			item := info.(map[string]interface{})
			stockName := item["n"].(string)
			if strings.Contains(stockName, "ST") {
				continue
			}
			var changeTime string
			t := int(item["tm"].(float64))
			if t/100000 > 0 {
				changeTime = fmt.Sprintf("%d", t)
			} else {
				changeTime = fmt.Sprintf("0%d", t)
			}
			stockId := item["c"].(string)
			if strings.Contains(stockId, "30") || strings.Contains(stockId, "68") {
				continue
			}
			changeType := fmt.Sprintf("%d", int(item["t"].(float64)))
			_, ok := stockChangeMapInfo[stockId]
			if ok {
				changeTime = "    " + changeTime
				if changeType == "64" {
					stockChangeMapInfo[stockId]["type64"] = stockChangeMapInfo[stockId]["type64"].(int) + 1
					stockChangeMapInfo[stockId]["type64_time"] = changeTime
				} else if changeType == "8201" {
					stockChangeMapInfo[stockId]["type8201"] = stockChangeMapInfo[stockId]["type8201"].(int) + 1
					stockChangeMapInfo[stockId]["type8201_time"] = changeTime
				} else if changeType == "8193" {
					stockChangeMapInfo[stockId]["type8193"] = stockChangeMapInfo[stockId]["type8193"].(int) + 1
					stockChangeMapInfo[stockId]["type8193_time"] = changeTime
				}
			} else {
				mapInfo := make(map[string]interface{})
				if changeType == "64" {
					mapInfo["type64"] = 1
					mapInfo["type64_time"] = changeTime
					mapInfo["type8201"] = 0
					mapInfo["type8201_time"] = ""
					mapInfo["type8193"] = 0
					mapInfo["type8193_time"] = ""
				} else if changeType == "8201" {
					mapInfo["type64"] = 0
					mapInfo["type64_time"] = ""
					mapInfo["type8201"] = 1
					mapInfo["type8201_time"] = changeTime
					mapInfo["type8193"] = 0
					mapInfo["type8193_time"] = ""
				} else if changeType == "8193" {
					mapInfo["type64"] = 0
					mapInfo["type64_time"] = ""
					mapInfo["type8201"] = 0
					mapInfo["type8201_time"] = ""
					mapInfo["type8193"] = 1
					mapInfo["type8193_time"] = changeTime
				}

				mapInfo["stockName"] = stockName
				stockChangeMapInfo[stockId] = mapInfo
			}
		}
		if len(stockList) < 64 {
			return
		}
		pageIndex += 1
	}
}

func GetBkInfo() {
	ts := time.Now().Unix()
	callBack := fmt.Sprintf("jQuery1123022593397568009088_%d", ts)
	webUrl := "http://98.push2.eastmoney.com/api/qt/clist/get?cb=" + callBack + "&pn=1&pz=20&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:90+t:2+f:!50&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f26,f22,f33,f11,f62,f128,f136,f115,f152,f124,f107,f104,f105,f140,f141,f207,f208,f209,f222"
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

	//reg := regexp.MustCompile("jQuery112405597814244474324_1634826682821\\(([\\s\\S]+?)\\);")
	reg := regexp.MustCompile(callBack + "\\(([\\s\\S]+?)\\);")
	rs := reg.FindAllSubmatch(content, -1)
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

	//for _, item := range diffData {
	//	bkId := item.(map[string]interface{})["f12"].(string)
	//	log.Info().Msgf("v is: %v", item)
	//	GetBkStockInfo(bkId)
	//}

	// 板块前5
	var stockList []*model.SimpleStockInfo
	var bkList []string
	for index := 0; index < 5; index++ {
		item := diffData[index]
		bkId := item.(map[string]interface{})["f12"].(string)
		bkName := item.(map[string]interface{})["f14"].(string)
		lst := GetBkStockInfo(bkId)
		if len(lst) > 0 {
			stockList = append(stockList, lst...)
		}
		bkList = append(bkList, bkName)
	}

	log.Info().Msg("===================================================================================================")
	for _, item := range bkList {
		log.Info().Msgf("bk name is: %s", item)
	}

	//log.Info().Msg("===================================================================================================")
	//for _, item := range stockList {
	//	if strings.Contains(item.StockName, "ST") {
	//		continue
	//	}
	//	s := &model.StockInfo{
	//		StockId:     item.StockId,
	//		StockName:   item.StockName,
	//		StockMarket: item.StockMarket,
	//	}
	//	GetStockMinuteInfo(s)
	//	log.Info().Msgf("currentPrice is: %.2f, stock Id is: %s, name is: %s, rate is: %.2f", item.CurrentPrice, item.StockId, item.StockName, item.Rate)
	//}
}

// 获取板块个股的信息
func GetBkStockInfo(bkId string) (stockList []*model.SimpleStockInfo) {
	webUrl := "http://push2.eastmoney.com/api/qt/clist/get?cb=jQuery112308642520074849604_1634823017240&fid=f62&po=1&pz=50&pn=1&np=1&fltt=2&invt=2&fs=b%3A" + bkId + "&fields=f12%2Cf14%2Cf2%2Cf3%2Cf62%2Cf184%2Cf66%2Cf69%2Cf72%2Cf75%2Cf78%2Cf81%2Cf84%2Cf87%2Cf204%2Cf205%2Cf124%2Cf1%2Cf13"
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

	reg := regexp.MustCompile("jQuery112308642520074849604_1634823017240\\(([\\s\\S]+?)\\);")
	rs := reg.FindAllSubmatch(content, -1)
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

	for _, item := range diffData {
		source := item.(map[string]interface{})
		currentPrice, _ := source["f2"].(float64)
		rate, _ := source["f3"].(float64)
		stockId, _ := source["f12"].(string)
		stockMarket := fmt.Sprintf("%.f", source["f13"].(float64))
		name, _ := source["f14"].(string)

		stock := &model.SimpleStockInfo{
			StockMarket:  stockMarket,
			StockName:    name,
			StockId:      stockId,
			Rate:         rate,
			CurrentPrice: currentPrice,
		}
		if rate > 0 && rate < 3 {
			stockList = append(stockList, stock)
		}
	}
	return
}

func GetStockMinuteInfo(stock *model.StockInfo) {
	secId := fmt.Sprintf("%s.%s", stock.StockMarket, stock.StockId)
	log.Info().Msgf("secId is: %s", secId)
	webUrl := fmt.Sprintf("http://push2.eastmoney.com/api/qt/stock/trends2/get?fields1=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,f12,f13&fields2=f51,f52,f53,f54,f55,f56,f57,f58&ndays=1&iscr=0&secid=%s&cb=jQuery112405294932494445095_1634734830512", secId)
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

	reg := regexp.MustCompile("jQuery112405294932494445095_1634734830512\\(([\\s\\S]+?)\\);")
	rs := reg.FindAllSubmatch(content, -1)

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
	trendData := data["trends"].([]interface{})
	var minuteInfoList []*model.StockMinuteInfo
	var price float64 // 当前价格及最后一条数据的价格
	for _, item := range trendData {
		source, ok := item.(string)
		if !ok {
			return
		}

		dataList := strings.Split(source, ",")
		currentTime := dataList[0]
		currentPrice, err := strconv.ParseFloat(dataList[2], 64)
		if err != nil {
			log.Error().Msgf("parse current price error: %v", err)
			continue
		}
		price = currentPrice
		averagePrice, err := strconv.ParseFloat(dataList[len(dataList)-1], 64)
		if err != nil {
			log.Error().Msgf("parse current price error: %v", err)
		}
		info := &model.StockMinuteInfo{
			StockId:      stock.StockId,
			StockName:    stock.StockName,
			AveragePrice: averagePrice,
			CurrentPrice: currentPrice,
		}
		diff := math.Abs((currentPrice - averagePrice) * 100 / averagePrice)
		if diff > 3 {
			if currentPrice > averagePrice {
				log.Info().Msgf("it's time to sell: %s, diff is: %.2f ========================================================", currentTime, diff)
			} else {
				log.Info().Msgf("it's time to buy: %s, diff is: %.2f ========================================================", currentTime, diff)
			}
		}
		minuteInfoList = append(minuteInfoList, info)
		log.Info().Msgf("current time is: %s, current price is: %.2f, average price is: %.2f --------", currentTime, currentPrice, averagePrice)
	}

	maxIndex, minIndex := PeekIndex(minuteInfoList)
	highestPrice := minuteInfoList[maxIndex].CurrentPrice
	lowestPrice := minuteInfoList[minIndex].CurrentPrice

	highDistance := (highestPrice - price) * 100 / price
	lowDistance := (lowestPrice - price) * 100 / price
	log.Info().Msgf("highDistance is: %.2f, lowDistance is: %.2f", highDistance, lowDistance)

	// 买点
	if highDistance > 3 && math.Abs(lowDistance) < 2 {
		log.Info().Msgf("it's time to buy: %s, highDistance is: %.2f,lowDistance is: %.2f", stock.StockId, highDistance, lowDistance)
	}
	// 计算均线的斜率
}

func PeekIndex(dataList []*model.StockMinuteInfo) (maxIndex int, minIndex int) {
	var max float64
	var min float64
	max = dataList[0].CurrentPrice
	min = dataList[0].CurrentPrice

	for index, item := range dataList {
		if item.CurrentPrice > max {
			max = item.CurrentPrice
			maxIndex = index
		}
		if item.CurrentPrice < min {
			min = item.CurrentPrice
			minIndex = index
		}
	}
	return
}

func FilterHeadStock(stockList []*model.StockInfo) (infoList []*model.StockInfo) {

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

		if info.CurrentPrice > 50 || info.CurrentPrice < 3 {
			// 价格相对较高或低
			continue
		}
		if info.CurrentRate < 7 {
			continue
		}
		diff := info.HighestRate - info.ChangeRate
		if diff > 3 {
			infoList = append(infoList, info)
		}
		log.Info().Msg("-----------------------------------------------------------------")
		log.Info().Msgf("name is: %s, code is: %v, currentPrice is: %v", info.StockName, info.StockId, info.CurrentPrice)
		log.Info().Msgf("highestPrice is: %v, lowestPrice is: %v", info.HighestPrice, info.LowestPrice)
		log.Info().Msgf("rate is: %v, amplitude is: %v, diff is: %v", info.CurrentRate, info.Amplitude, diff)
		log.Info().Msgf("changeRate is: %v, highestRate is: %v, lowestRate is: %v", info.ChangeRate, info.HighestRate, info.LowestRate)
	}
	return
}

func FilterStock(stockList []*model.StockInfo) (infoList []*model.StockInfo) {

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
		if info.CurrentRate != 1.27 {
			continue
		}
		diff := info.CurrentRate - info.HighestRate
		log.Info().Msg("-----------------------------------------------------------------")
		log.Info().Msgf("name is: %s, code is: %v, currentPrice is: %v", info.StockName, info.StockId, info.CurrentPrice)
		log.Info().Msgf("highestPrice is: %v, lowestPrice is: %v", info.HighestPrice, info.LowestPrice)
		log.Info().Msgf("rate is: %v, amplitude is: %v, diff is: %v", info.CurrentRate, info.Amplitude, diff)
		log.Info().Msgf("changeRate is: %v, highestRate is: %v, lowestRate is: %v", info.ChangeRate, info.HighestRate, info.LowestRate)
		infoList = append(infoList, info)
	}
	return
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

func SaveStockInfo(stockInfoList []*model.StockInfo) {
	if len(stockInfoList) <= 0 {
		return
	}

	for _, stock := range stockInfoList {
		db.SaveStock(stock)
	}
}

func GetStockList() (stockList []*model.StockInfo) {
	var page = 1
	var total int
	timeString := time.Now().Format("2006-01-02")
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
			//log.Info().Msgf("v is: %v", v)
			info := v.(map[string]interface{})
			codeStr, ok := info["f12"]
			if !ok {
				continue
			}
			code := codeStr.(string)

			marketStr, ok := info["f13"]
			if !ok {
				continue
			}
			market, _ := marketStr.(string)
			//if strings.HasPrefix(code, "300") || strings.HasPrefix(code, "688") {
			//	continue
			//}
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
			//startPrice, ok := info["f17"].(float64)
			//if !ok {
			//	continue
			//}
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
			//log.Info().Msgf("info is: %v", info)
			//log.Info().Msgf("code is: %v, currentPrice is: %v", code, currentPrice)
			//log.Info().Msgf("highestPrice is: %v, lowestPrice is: %v, startPrice is: %v", highestPrice, lowestPrice, startPrice)
			//log.Info().Msgf("rate is: %v, amplitude is: %v, diff is: %v", rate, amplitude, diff)
			//log.Info().Msgf("changeRate is: %v, highestRate is: %v, lowestRate is: %v", changeRate, highestRate, lowestRate)

			stockInfo := &model.StockInfo{
				StockId:      fmt.Sprintf("%s", code),
				StockName:    name,
				StockMarket:  market,
				TimeString:   timeString,
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
	return
}

func getReverseStock() {

}
