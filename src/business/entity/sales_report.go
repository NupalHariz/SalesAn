package entity

import (
	"github.com/reyhanmichiels/go-pkg/v2/null"
)

type Status string

type SalesReport struct {
	Id           int64       `db:"id"`
	UserId       int64       `db:"user_id"`
	FileUrl      string      `db:"file_url"`
	ErrorMessage null.String `db:"error_message"`
	StartAt      null.Time   `db:"start_at"`
	CompletedAt  null.Time   `db:"completed_at"`
}

type SalesReportParam struct {
	FileUrl string `db:"file_url"`
}

type SalesReportUpdateParam struct {
	ErrorMessage null.String `db:"error_message"`
	StartAt      null.Time   `db:"start_at"`
	CompletedAt  null.Time   `db:"completed_at"`
}
