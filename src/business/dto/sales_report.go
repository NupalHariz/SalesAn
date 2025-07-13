package dto

import "mime/multipart"

type UploadReportParam struct {
	Report *multipart.FileHeader `form:"report"`
}

type GetReportList struct {
	Id      int64  `json:"id"`
	FileUrl string `json:"file_url"`
	Status  string `json:"status"`
}