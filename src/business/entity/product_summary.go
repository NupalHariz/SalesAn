package entity

type ProductSummary struct {
	Id          int64  `db:"id"`
	ReportId    int64  `db:"report_id"`
	ProductName string `db:"product_name"`
	Quantity    int64  `db:"quantity"`
	Revenue     int64  `db:"revenue"`
}
