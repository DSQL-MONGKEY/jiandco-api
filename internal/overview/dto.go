package overview

type DailyRevenue struct {
	Date    string `json:"date"`
	Revenue int64  `json:"revenue"`
	Orders  int64  `json:"orders"`
}

type OrderStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type Response struct {
	Month                string             `json:"month"`
	Revenue              int64              `json:"revenue"`
	Expenses             int64              `json:"expenses"`
	NetProfit            int64              `json:"net_profit"`
	TotalOrders          int64              `json:"total_orders"`
	PaidOrders           int64              `json:"paid_orders"`
	AverageOrderValue    int64              `json:"average_order_value"`
	OrderStatusBreakdown []OrderStatusCount `json:"order_status_breakdown"`
	Daily                []DailyRevenue     `json:"daily"`
}
