package model

type Level2TradeInfo struct {
	TradeId         string
	SaleOrderID     string
	BuyOrderID      string
	SaleOrderVolume int64
	BuyOrderVolume  int64
	TradeTime       string
	TradeDay        string
	StockId         string
}
