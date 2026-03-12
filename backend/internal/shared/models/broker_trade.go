package models

import "time"

// BrokerBuy represents broker buying transaction data
type BrokerBuy struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BrokerCode    string    `gorm:"type:varchar(10);not null;index" json:"broker_code"`
	StockCode     string    `gorm:"type:varchar(10);not null;index" json:"stock_code"`
	TradeDate     time.Time `gorm:"type:date;not null;index" json:"trade_date"`
	BuyLot        float64   `gorm:"type:decimal(20,2)" json:"buy_lot"`
	BuyLotValue   string    `gorm:"type:varchar(50)" json:"buy_lot_value"`
	BuyValue      float64   `gorm:"type:decimal(20,2)" json:"buy_value"`
	BuyValueValue string    `gorm:"type:varchar(50)" json:"buy_value_value"`
	BuyAvgPrice   float64   `gorm:"type:decimal(20,10)" json:"buy_avg_price"`
	Type          string    `gorm:"type:varchar(50);index" json:"type"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BrokerBuy) TableName() string {
	return "broker_buys"
}

// BrokerSell represents broker selling transaction data
type BrokerSell struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BrokerCode     string    `gorm:"type:varchar(10);not null;index" json:"broker_code"`
	StockCode      string    `gorm:"type:varchar(10);not null;index" json:"stock_code"`
	TradeDate      time.Time `gorm:"type:date;not null;index" json:"trade_date"`
	SellLot        float64   `gorm:"type:decimal(20,2)" json:"sell_lot"`
	SellLotValue   string    `gorm:"type:varchar(50)" json:"sell_lot_value"`
	SellValue      float64   `gorm:"type:decimal(20,2)" json:"sell_value"`
	SellValueValue string    `gorm:"type:varchar(50)" json:"sell_value_value"`
	SellAvgPrice   float64   `gorm:"type:decimal(20,10)" json:"sell_avg_price"`
	Type           string    `gorm:"type:varchar(50);index" json:"type"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BrokerSell) TableName() string {
	return "broker_sells"
}
