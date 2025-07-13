package entity

import "time"

type Status string

type SalesReport struct {
	Id           int64     `db:"id"`
	UserId       int64     `db:"user_id"`
	FileUrl      string    `db:"file_url"`
	ErrorMessage string    `db:"error_message"`
	StartAt      time.Time `db:"start_at"`
	CompletedAt  time.Time `db:"completed_at"`
}
