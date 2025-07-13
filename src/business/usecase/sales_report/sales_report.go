package salesreport

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	salesreportDom "github.com/NupalHariz/SalesAn/src/business/domain/sales_report"
	"github.com/NupalHariz/SalesAn/src/business/dto"
	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/NupalHariz/SalesAn/src/business/service/supabase"
	"github.com/nao1215/csv"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/files"
	"github.com/xuri/excelize/v2"
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
	case "xlam", "xlsm", "xlsx", "xltx", "xltm":
		valid, err = s.validateExceel(param.Report)
	default:
		return "", errors.New("invalid file extension")
	}

	if !valid {
		return "", err
	}

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

func (s *salesReport) validateExceel(param *multipart.FileHeader) (bool, error) {
	file, err := param.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()

	rFile, err := excelize.OpenReader(file)
	if err != nil {
		return false, err
	}

	sheetName := rFile.GetSheetName(0)
	if sheetName == "" {
		return false, errors.New("sheet can't be empty")
	}

	rows, err := rFile.GetRows(sheetName)
	if err != nil {
		return false, err
	}

	var validationError []string

	columnNames := map[int]string{
		0: "InvoiceID",
		1: "Date",
		2: "CustomerName",
		3: "Item",
		4: "Quantity",
		5: "UnitPrice",
		6: "Total",
		7: "Status",
		8: "PaymentMethod",
	}

	for i, row := range rows[1:] {
		line := i + 1

		if len(row) < 9 {
			return false, errors.New("there should be 9 column")
		}

		emptyCol := make(map[int]bool)

		for colID, val := range row {
			if strings.TrimSpace(val) == "" {
				emptyCol[colID] = true

				validationError = append(validationError, fmt.Sprintf(
					"line:%d column %s has an empty columny", line, columnNames[colID],
				))
			}
		}

		if !emptyCol[1] {
			_, dErr := time.Parse("2006-01-02", row[1])
			if dErr != nil {
				validationError = append(validationError, fmt.Sprintf(
					"line:%d column date: expected format YYYY-MM-DD, got=%s", line, row[1],
				))
			}
		}

		if !emptyCol[4] {
			quantity, qErr := strconv.Atoi(row[4])
			if qErr != nil || quantity < 0 {
				validationError = append(validationError, fmt.Sprintf(
					"line:%d column Quantity: must be numeric and > 0, got = %s", line, row[4],
				))
			}
		}

		if !emptyCol[5] {
			unitPrice, uErr := strconv.Atoi(row[5])
			if uErr != nil || unitPrice < 0 {
				validationError = append(validationError, fmt.Sprintf(
					"line:%d column Unit Price: must be numeric and > 0, got = %s", line, row[5],
				))
			}
		}

		if !emptyCol[6] {
			total, tErr := strconv.Atoi(row[6])
			if tErr != nil || total < 0 {
				validationError = append(validationError, fmt.Sprintf(
					"line:%d column Total: must be numeric and > 0, got = %s", line, row[6],
				))
			}
		}
	}

	if len(validationError) > 0 {
		return false, errors.New(strings.Join(validationError, ", "))
	}

	return true, nil
}
