package entity

type SalesSummary struct {
	Id                int64  `db:"id"`
	ReportId          int64  `db:"report_id"`
	TotalTransaction  int64  `db:"total_transaction"`
	Success           int64  `db:"success"`
	Failed            int64  `db:"failed"`
	TotalRevenue      int64  `db:"total_revenue"`
	MostPaymentMethod string `db:"most_payment_method"`
}

type SalesSummaryParam struct {
	ReportId int64 `db:"report_id" param:"report_id"`
}
