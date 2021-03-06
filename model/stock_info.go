package model

// 分时
type StockMinuteInfo struct {
	Id           int64 `xorm:"pk autoincr"`
	StockId      string
	StockName    string
	TimeStamp    int64
	AveragePrice float64 // 均价
	CurrentPrice float64
	HighestPrice float64
	LowestPrice  float64
	HighestRate  float64
	LowestRate   float64 // 通过最高、低价与开盘价计算得出
	CurrentRate  float64
	Diff         float64 // 距离前高
}

// 筛出的目标
type TargetStockInfo struct {
	Id           int64 `xorm:"pk autoincr"`
	StockId      string
	StockName    string
	TimeStamp    int64
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

// 每天的涨跌
type StockInfo struct {
	Id             int64 `xorm:"pk autoincr"`
	StockId        string
	StockName      string
	StockMarket    string  // 市场类型
	TimeString     string  // 时间
	YesterdayPrice float64 // 昨日收盘价
	CurrentPrice   float64
	HighestPrice   float64
	LowestPrice    float64
	HighestRate    float64
	LowestRate     float64 // 通过最高、低价与开盘价计算得出
	CurrentRate    float64
	ChangeRate     float64 // 换手率
	Amplitude      float64 // 振幅
	Diff           float64 // 距离前高
}

type SimpleStockInfo struct {
	StockMarket  string
	StockName    string
	StockId      string
	Rate         float64
	CurrentPrice float64
}

type StockMoneyFlow struct {
	StockId        string
	MainStream     int64 // 主力
	Super          int64
	Big            int64
	Middle         int64
	Small          int64
	MainStreamRate float64 // 主力
	SuperRate      float64
	BigRate        float64
	MiddleRate     float64
	SmallRate      float64
	TimeString     string
}
