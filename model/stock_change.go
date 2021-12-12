package model

// 股票异动
type StockChange struct {
	StockId                string
	StockName              string
	BigMaiPan              int    // 类型 有大买盘
	RocketLaunch           int    // 类型 火箭发射
	QuantityBuy            int    // 类型 大笔买入
	BigMaiPanChangeTime    string // 异动时间点
	RocketLaunchChangeTime string // 异动时间点
	QuantityBuyChangeTime  string // 异动时间点
	TimeString             string // 异动时间
}
