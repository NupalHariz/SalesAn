package dto

import "mime/multipart"

type UploadReportParam struct {
	Report *multipart.FileHeader `form:"report"`
}
