package salesreport

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	salesreportDom "github.com/NupalHariz/SalesAn/src/business/domain/sales_report"
	"github.com/NupalHariz/SalesAn/src/business/dto"
	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/NupalHariz/SalesAn/src/business/service/supabase"
	"github.com/nao1215/csv"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/files"
)

type Interface interface {
	UploadReport(ctx context.Context, param dto.UploadReportParam) (string, error)
}

type salesReport struct {
	salesRepostDom salesreportDom.Interface
	auth           auth.Interface
	supabase       supabase.Interface
}

type InitParam struct {
	SalesReportDom salesreportDom.Interface
	Auth           auth.Interface
	Supabase       supabase.Interface
}

func Init(param InitParam) Interface {
	return &salesReport{
		salesRepostDom: param.SalesReportDom,
		auth:           param.Auth,
		supabase:       param.Supabase,
	}
}

func (s *salesReport) UploadReport(ctx context.Context, param dto.UploadReportParam) (string, error) {
	userLogin, err := s.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return "", err
	}
	fileExt := files.GetExtension(param.Report.Filename)

	var valid bool
	switch strings.ToLower(fileExt) {
	case "csv":
		valid, err = s.validateCSV(param.Report)
	}

	if !valid {
		return "", err
	}
	// VALIDATE XSLX (Later)

	param.Report.Filename = fmt.Sprintf("%v-%v", time.Now().String(), param.Report.Filename)
	param.Report.Filename = strings.Replace(param.Report.Filename, " ", "-", -1)

	url, err := s.supabase.Upload(param.Report)
	if err != nil {
		return "", err
	}

	salesReport := entity.SalesReport{
		UserId:  userLogin.ID,
		FileUrl: url,
	}

	//SEND TO QUEUE
	//To Do Later Saja Fokus 3 lainnya dulu -> Get -> Queue

	err = s.salesRepostDom.Create(ctx, salesReport)
	if err != nil {
		return "", err
	}

	successMsg := "processing your file"
	return successMsg, err
}

func (s *salesReport) validateCSV(param *multipart.FileHeader) (bool, error) {
	file, err := param.Open()
	if err != nil {
		return false, err
	}

	r, err := csv.NewCSV(file)
	if err != nil {
		return false, err
	}
	report := make([]entity.Report, 0)
	errorss := r.Decode(&report)

	for i, r := range report {
		_, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			msg := fmt.Errorf("line:%d column date: expected format YYYY-MM-DD, got=%s", i, r.Date)
			errorss = append(errorss, msg)
		}
	}

	if len(errorss) > 0 {
		errMsg := make([]string, len(errorss))
		for i, e := range errorss {
			errMsg[i] = e.Error()
		}

		return false, errors.New(strings.Join(errMsg, ", "))
	}

	return true, nil
}
