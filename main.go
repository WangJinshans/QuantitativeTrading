package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"github.com/rs/zerolog/log"
	"quant_trade/data_center"
	"quant_trade/db"
	"quant_trade/model"
	"time"
)

func init() {
}

func WatchStock(ctx context.Context) {

	c := cron.New()           //精确到秒
	spec := "00 30 11 * * ? " //cron表达式 每天11:30

	c.AddFunc(spec, func() {
		infoList := data_center.GetStockList()
		targetStockList := data_center.FilterStock(infoList)

		for _, info := range targetStockList {
			stock := new(model.TargetStockInfo)
			stock.Id = info.Id
			stock.StockId = info.StockId
			stock.StockName = info.StockName
			stock.TimeStamp = time.Now().Unix()
			stock.CurrentPrice = info.CurrentPrice
			stock.HighestPrice = info.HighestPrice
			stock.LowestPrice = info.LowestPrice
			stock.HighestRate = info.HighestRate
			stock.LowestRate = info.LowestRate
			stock.CurrentRate = info.CurrentRate
			stock.ChangeRate = info.ChangeRate
			stock.Amplitude = info.Amplitude
			stock.Diff = info.Diff
			db.SaveStock(stock)
		}
		if len(targetStockList) > 0 {
			go WatchTargetStock(targetStockList)
		}
	})
	c.Start()
}

// 早盘大幅拉涨的票
func WatchHeadStock(ctx context.Context) {

	c := cron.New()           //精确到秒
	spec := "00 40 09 * * ? " //cron表达式 每天11:30

	c.AddFunc(spec, func() {
		infoList := data_center.GetStockList()
		targetStockList := data_center.FilterHeadStock(infoList)

		for _, info := range targetStockList {
			stock := new(model.TargetStockInfo)
			stock.Id = info.Id
			stock.StockId = info.StockId
			stock.StockName = info.StockName
			stock.TimeStamp = time.Now().Unix()
			stock.CurrentPrice = info.CurrentPrice
			stock.HighestPrice = info.HighestPrice
			stock.LowestPrice = info.LowestPrice
			stock.HighestRate = info.HighestRate
			stock.LowestRate = info.LowestRate
			stock.CurrentRate = info.CurrentRate
			stock.ChangeRate = info.ChangeRate
			stock.Amplitude = info.Amplitude
			stock.Diff = info.Diff
			db.SaveStock(stock)
		}
		if len(targetStockList) > 0 {
			go WatchTargetStock(targetStockList)
		}
	})
	c.Start()
}

func WatchHeadStockAgain(ctx context.Context) {

	c := cron.New()           //精确到秒
	spec := "00 50 09 * * ? " //cron表达式 每天11:30

	c.AddFunc(spec, func() {
		infoList := data_center.GetStockList()
		targetStockList := data_center.FilterHeadStock(infoList)

		for _, info := range targetStockList {
			stock := new(model.TargetStockInfo)
			stock.Id = info.Id
			stock.StockId = info.StockId
			stock.StockName = info.StockName
			stock.TimeStamp = time.Now().Unix()
			stock.CurrentPrice = info.CurrentPrice
			stock.HighestPrice = info.HighestPrice
			stock.LowestPrice = info.LowestPrice
			stock.HighestRate = info.HighestRate
			stock.LowestRate = info.LowestRate
			stock.CurrentRate = info.CurrentRate
			stock.ChangeRate = info.ChangeRate
			stock.Amplitude = info.Amplitude
			stock.Diff = info.Diff
			db.SaveStock(stock)
		}
		if len(targetStockList) > 0 {
			go WatchTargetStock(targetStockList)
		}
	})
	c.Start()
}

func WatchTargetStock(targetStockList []*model.StockInfo) {

	for _, stock := range targetStockList {
		WatchSingleStock(stock)
	}
}

func WatchSingleStock(stock *model.StockInfo) {
	timer := time.NewTimer(1 * time.Minute)
	for {
		select {
		case <-timer.C:
			if time.Now().Hour() > 15 {
				return
			}
			data_center.GetStockMinuteInfo(stock)
			timer.Reset(1 * time.Second)
		}
	}
}

func main() {

	//infoList := GetStockList()
	//log.Info().Msg("------------------------------------------------------------------------------------------------------------------")
	//FilterStock(infoList)

	//data_center.GetBkInfo()

	//signChan := make(chan os.Signal)
	//signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
	//
	////ctx := context.Background()
	////go WatchStock(ctx)
	//stock := &model.StockInfo{
	//	StockId:        "301086",
	//	StockMarket:    "0",
	//	YesterdayPrice: 96.67,
	//}
	//data_center.GetStockMinuteInfo(stock)

	//<-signChan

	messages, residueBytes, invalidMessages := Split([]byte(""))
	for _, item := range messages {
		log.Info().Msgf("message is: %x", item)
	}
	log.Info().Msgf("len message is: %x", len(messages))
	log.Info().Msgf("residueBytes is: %x", residueBytes)
	log.Info().Msgf("invalidMessages is: %x", invalidMessages)

}

func Split(segment []byte) (messages [][]byte, residueBytes []byte, invalidMessages [][]byte) {

	var indexList []int
	var startFlag = []byte("##")
	var index int = 0

	s := "2323024c45575043413030304a463236303133340101009114060b10131801c4b6232375218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa"
	s = "5043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa"
	s = "3322222323024c45575043413030304a463236303133340101009114060b10131801c4b6232375218ba64c45575043413023234a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa23233233"
	s = "3322222323024c45575043413030304a463236303133340101009114060b10131801c4b6232375218ba64c45575043413055555555555555555523234a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa23233233"
	//s = "3322222323024c45575043413030304a463236303133340101009114060b10131801c4b6232375218ba64c455750434130555555555232355555555523234a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa2323024c45575043413030304a463236303133340101009114060b10131801c4b6000075218ba64c45575043413030304a4632363031333431313131313131313131313131313131313131313131313131313131313131313131313130303030303030303030303030303030303030303030303030303030303030303030303000025f82890f6d7a8a8739000100020fdcaf003c224047a058822e710005ec987d020600d68f8922a2fa23233233"

	segment, err := hex.DecodeString(s)
	if err != nil {
		log.Info().Msgf("error is: %v", err)
		return
	}

	for i := 0; i < len(segment)-1; i += 1 {
		sf := segment[index : index+2]
		if bytes.Equal(sf, startFlag) {
			indexList = append(indexList, index)
			index += 2
			continue
		}
		index += 1
	}

	if len(indexList) > 0 {
		if indexList[0] > 0 {
			// 前面的干扰数据 考虑2323为原始数据部分,头部的2323丢失的可能性(基本不存在的原因,头部的2323会被添加到当前字节中)
			data := segment[0:indexList[0]]
			invalidMessages = append(invalidMessages, data)
			//segment = segment[indexList[0]:] 写在此处影响切片index的判断
		}
	} else {
		// 没有2323
		residueBytes = append(residueBytes, segment...)
		return
	}

	var flag int // 2323标识的位置
	for flag < len(indexList) {
		var seg []byte
		if len(indexList) == 1 {
			// 一个2323
			seg = segment[indexList[flag]:]
			data, left, err := verifySplitPackage(seg)
			log.Info().Msgf("data is: %x, left is: %x, error is: %v", data, left, err)
			if err != nil {
				// 不完整的一个包
				residueBytes = append(residueBytes, seg...)
			} else {
				messages = append(messages, data)
			}
			if len(left) > 0 {
				// 一个完整包到下个2323之前的数据为无效数据
				invalidMessages = append(invalidMessages, left)
			}
			flag += 1
		} else {
			// 多个2323 seg 存储每个2323到下个2323的数据
			nextFlag := flag + 1
			for nextFlag <= len(indexList) {

				log.Info().Msgf("flag is: %d, next flag is: %d", flag, nextFlag)
				if nextFlag == len(indexList) {
					seg = segment[indexList[flag]:] // 最后一条
					log.Info().Msgf("last seg is: %x", seg)
				} else {
					seg = segment[indexList[flag]:indexList[nextFlag]]
					log.Info().Msgf("seg is: %x", seg)
				}

				data, left, err := verifySplitPackage(seg)
				if err != nil {
					// 不完整的一个包
					if nextFlag == len(indexList) {
						// 最后一个写入剩余
						residueBytes = append(residueBytes, seg...)
						log.Info().Msgf("residueBytes error: %x", residueBytes)
					}
					nextFlag += 1
					continue
				}
				if len(left) > 0 {
					// 一个完整包到下个2323之前的数据为无效数据
					log.Info().Msgf("left append: %x", left)
					invalidMessages = append(invalidMessages, left)
				}
				messages = append(messages, data)
				flag = nextFlag
				nextFlag += 1
			}
			flag += 1
		}
	}
	return messages, residueBytes, invalidMessages
}

// 再次切割 如果bs长度超过其中length的长度, 返回pkg以及无用的剩余部分
// err 判断是否存在完整包的情况
func verifySplitPackage(bs []byte) (pkg []byte, left []byte, err error) {
	if len(bs) < 25 {
		err = errors.New("length error")
		return
	}
	payloadLen := int(binary.BigEndian.Uint16(bs[22:24]))
	length := payloadLen + 25
	if length > len(bs) {
		err = errors.New("bad length")
		return
	}
	data := bs[:length]
	err = checkBCC(data)
	if err != nil {
		return
	}
	pkg = data
	if length < len(bs) {
		left = bs[length:]
	}
	return
}

func checkBCC(bs []byte) (err error) {
	var checksum byte
	for _, b := range bs {
		checksum ^= b
	}
	if checksum == 0 {
		return
	}
	err = errors.New("bcc error")
	return
}
