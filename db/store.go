package db

import (
	"github.com/go-xorm/xorm"
	"github.com/rs/zerolog/log"
	"quant_trade/model"
	"time"
	"xorm.io/core"
)

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:xxxxxxxxxx@tcp(127.0.0.1:3306)/trade?charset=utf8")
	if err != nil {
		log.Error().Msgf("connect to database error: %v", err)
		return
	}
	err = engine.Sync2(new(model.StockInfo))
	if err != nil {
		log.Error().Msgf("sync2 table error: %v", err)
	}

	err = engine.Sync2(new(model.Level2TradeInfo))
	if err != nil {
		log.Error().Msgf("sync2 table error: %v", err)
	}

	err = engine.Sync2(new(model.StockChange))
	if err != nil {
		log.Error().Msgf("sync2 table error: %v", err)
	}

	err = engine.Sync2(new(model.StockMoneyFlow))
	if err != nil {
		log.Error().Msgf("sync2 table error: %v", err)
	}
}

// 存储到mysql中
func SaveStock(stockInfo interface{}) {
	engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	engine.SetTableMapper(core.SnakeMapper{})
	affected, err := engine.Insert(stockInfo)
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("affected is: %v", affected)
}

func SaveStockMoneyFlow(info []model.StockMoneyFlow) {
	engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	engine.SetTableMapper(core.SnakeMapper{})
	affected, err := engine.Insert(info)
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("affected is: %v", affected)
}

func SaveLevel2TradeInfo(tradeInfo []model.Level2TradeInfo) {
	engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	engine.SetTableMapper(core.SnakeMapper{})
	affected, err := engine.Insert(tradeInfo)
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("affected is: %v", affected)
}

func SaveChangeInfo(changeInfo model.StockChange) {
	engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	engine.SetTableMapper(core.SnakeMapper{})
	affected, err := engine.Insert(changeInfo)
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("affected is: %v", affected)
}

func GetChangeInfo() (changeInfo []model.StockChange, err error) {
	//engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	err = engine.Where("(big_mai_pan+rocket_launch+quantity_buy)>?", 6).And("time_string=?", time.Now().Format("2006-01-02")).Find(&changeInfo)
	//err = engine.Where("(big_mai_pan+rocket_launch+quantity_buy)>?", 6).And("time_string=?", "2021-12-14").Find(&changeInfo)
	if err != nil {
		log.Info().Msgf("fail to get data: %v", err)
		return
	}
	return
}

func GetStockInfo(stockInfo model.StockInfo) (info model.StockInfo, err error) {
	//engine.ShowSQL(true) // 显示SQL的执行, 便于调试分析
	info.TimeString = stockInfo.TimeString
	info.StockId = stockInfo.StockId
	_, err = engine.Get(&info)
	if err != nil {
		log.Info().Msgf("fail to get data: %v", err)
		return
	}
	return
}
