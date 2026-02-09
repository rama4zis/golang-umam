package dto

type ReportTodayResponse struct {
	TotalRevenue       int `json:"total_revenue"`
	TotalTransaction   int `json:"total_transaksi"`
	BestSellingProduct struct {
		Name     string `json:"nama"`
		Quantity int    `json:"qty_terjual"`
	} `json:"produk_terlaris"`
}
