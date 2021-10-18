package model

type StockInfo struct {
	Id           int64 `xorm:"pk autoincr"`
	StockId      string
	StockName    string
	timeString   string // 时间
	CurrentPrice float64
	HighestPrice float64
	LowestPrice  float64
	HighestRate  float64
	LowestRate   float64 // 通过最高、低价与开盘价计算得出
	CurrentRate  float64
	ChangeRate   float64 // 换手率
	Amplitude    float64 // 振幅
	Diff         float64 // 距离前高
}
