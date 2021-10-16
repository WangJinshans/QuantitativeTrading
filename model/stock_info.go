package model

type StockInfo struct {
	StockId      string
	CurrentPrice float64
	HighestPrice float64
	LowestPrice  float64
	HighestRate  float64
	LowestRate   float64 // 通过最高、低价与开盘价计算得出
	CurrentRate  float64
	ChangeRate   float64 // 换手率
}
