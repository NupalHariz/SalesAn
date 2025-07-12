package entity

import "time"

type DailySalesSummary struct {
	Id               int64     `db:"id"`
	ReportId         int64     `db:"report_id"`
	Date             time.Time `db:"date"`
	TotalTransaction int64     `db:"total_transaction"`
	TotalRevenue     int64     `db:"total_revenue"`
}
