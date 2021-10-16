package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"regexp"
)

var client http.Client

func main() {
	//GetStockList()
	GetBkInfo()
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

func GetStockList() {
	webUrl := "http://51.push2.eastmoney.com/api/qt/clist/get?cb=jQuery1124012243073358859502_1634284246158&pn=1&pz=30&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23&fields=f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152"
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

	//data := string(content)
	//log.Info().Msgf("content is: %s", data)
	reg := regexp.MustCompile("jQuery1124012243073358859502_1634284246158\\(([\\s\\S]+?)\\);")
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

// 存储到mysql中
func SaveData() {

}
