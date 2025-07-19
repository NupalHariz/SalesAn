package dto

import (
	"mime/multipart"

	"github.com/NupalHariz/SalesAn/src/business/entity"
)

type UploadReportParam struct {
	Report *multipart.FileHeader `form:"report"`
}

type GetReportList struct {
	Id      int64  `json:"id"`
	FileUrl string `json:"file_url"`
	Status  string `json:"status"`
}

type ReportParam struct {
	ReportId int64 `uri:"report_id"`
}

type SummaryReport struct {
	SalesSummary      entity.SalesSummary        `json:"sales_summary"`
	ProductSummary    []entity.ProductSummary    `json:"product_summary"`
	DailySalesSummary []entity.DailySalesSummary `json:"daily_sales_summary"`
	ErrorMessage      string                     `json:"error_message"`
}
