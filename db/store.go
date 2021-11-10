package db

import (
	"github.com/go-xorm/xorm"
	"github.com/rs/zerolog/log"
	"quant_trade/model"
	"xorm.io/core"
)

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:wangjinshan123..@tcp(127.0.0.1:3306)/trade?charset=utf8")
	if err != nil {
		log.Error().Msgf("connect to database error: %v", err)
		return
	}
	err = engine.Sync2(new(model.StockInfo))
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
