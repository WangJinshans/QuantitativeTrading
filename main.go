package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
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
	data_center.GetBkInfo()

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
}
