package salesreport

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	dailySalesSummaryDom "github.com/NupalHariz/SalesAn/src/business/domain/daily_sales_summary"
	productSummaryDom "github.com/NupalHariz/SalesAn/src/business/domain/product_summary"
	salesReportDom "github.com/NupalHariz/SalesAn/src/business/domain/sales_report"
	salesSumarryDom "github.com/NupalHariz/SalesAn/src/business/domain/sales_summary"
	"github.com/NupalHariz/SalesAn/src/business/dto"
	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/NupalHariz/SalesAn/src/business/service/supabase"
	"github.com/NupalHariz/SalesAn/src/handler/pubsub/publisher"
	"github.com/nao1215/csv"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	errorPkg "github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/files"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/xuri/excelize/v2"
)

type Interface interface {
	UploadReport(ctx context.Context, param dto.UploadReportParam) (string, error)
	ListReport(ctx context.Context) ([]dto.GetReportList, error)
	SummarizeReport(ctx context.Context, payload entity.PubSubMessage) error
}

type salesReport struct {
	salesReportDom       salesReportDom.Interface
	salesSummaryDom      salesSumarryDom.Interface
	productSummaryDom    productSummaryDom.Interface
	dailySalesSummaryDom dailySalesSummaryDom.Interface
	auth                 auth.Interface
	supabase             supabase.Interface
	publisher            publisher.Interface
	json                 parser.JSONInterface
}

type InitParam struct {
	SalesReportDom       salesReportDom.Interface
	SalesSummaryDom      salesSumarryDom.Interface
	ProductSummaryDom    productSummaryDom.Interface
	DailySalesSummaryDom dailySalesSummaryDom.Interface
	Auth                 auth.Interface
	Supabase             supabase.Interface
	Publisher            publisher.Interface
	Json                 parser.JSONInterface
}

func Init(param InitParam) Interface {
	return &salesReport{
		salesReportDom:       param.SalesReportDom,
		salesSummaryDom:      param.SalesSummaryDom,
		productSummaryDom:    param.ProductSummaryDom,
		dailySalesSummaryDom: param.DailySalesSummaryDom,
		auth:                 param.Auth,
		supabase:             param.Supabase,
		publisher:            param.Publisher,
		json:                 param.Json,
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
		valid, err = s.validateExcel(param.Report)
	default:
		return "", errorPkg.NewWithCode(codes.CodeBadRequest, "invalid file extension")
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

	salesReport, err = s.salesReportDom.Create(ctx, salesReport)
	if err != nil {
		return "", err
	}

	err = s.publisher.Publish(ctx, entity.ExchangeSalesReport, entity.KeySalesReport, salesReport)
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
	validationError := r.Decode(&report)

	for i, r := range report {
		_, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			msg := fmt.Errorf("line:%d column date: expected format YYYY-MM-DD, got=%s", i, r.Date)
			validationError = append(validationError, msg)
		}
	}

	if len(validationError) > 0 {
		errMsg := make([]string, len(validationError))
		for i, e := range validationError {
			errMsg[i] = e.Error()
		}

		return false, errorPkg.NewWithCode(codes.CodeBadRequest, strings.Join(errMsg, ", "))
	}

	return true, nil
}

func (s *salesReport) validateExcel(param *multipart.FileHeader) (bool, error) {
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
		return false, errorPkg.NewWithCode(codes.CodeBadRequest, "sheet can not be empty")
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
			return false, errorPkg.NewWithCode(codes.CodeBadRequest,
				"There must be at least 9 columns: InvoiceID, Date, CustomerName, Item, Quantity, UnitPrice, Total, Status, and PaymentMethod.",
			)
		}

		emptyCol := make(map[int]bool)

		for colID, val := range row {
			if strings.TrimSpace(val) == "" {
				emptyCol[colID] = true

				validationError = append(validationError, fmt.Sprintf(
					"line:%d column %s has an empty column", line, columnNames[colID],
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
					"line:%d column UnitPrice: must be numeric and > 0, got = %s", line, row[5],
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
		return false, errorPkg.NewWithCode(codes.CodeBadRequest, strings.Join(validationError, ", "))
	}

	return true, nil
}

func (s *salesReport) ListReport(ctx context.Context) ([]dto.GetReportList, error) {
	var res []dto.GetReportList

	salesReports, err := s.salesReportDom.GetList(ctx)
	if err != nil {
		return res, err
	}

	for _, s := range salesReports {
		var salesReport dto.GetReportList
		var status string

		if s.StartAt.IsNullOrZero() {
			status = "Waiting"
		} else if !s.StartAt.IsNullOrZero() && s.CompletedAt.IsNullOrZero() {
			status = "Processing"
		} else if !s.CompletedAt.IsNullOrZero() && !s.ErrorMessage.IsNullOrZero() {
			status = "Failed"
		} else {
			status = "Success"
		}

		salesReport = dto.GetReportList{
			Id:      s.Id,
			FileUrl: s.FileUrl,
			Status:  status,
		}

		res = append(res, salesReport)
	}

	_ = s.publisher.Publish(ctx, entity.ExchangeSalesReport, entity.KeyHi, "Naufal Haris, KING OF THE KINGS")

	return res, nil
}

func (s *salesReport) SummarizeReport(ctx context.Context, payload entity.PubSubMessage) error {
	var salesReport entity.SalesReport
	var err error
	err = s.json.Unmarshal([]byte(payload.Payload), &salesReport)
	if err != nil {
		return err
	}

	err = s.salesReportDom.Update(
		ctx,
		entity.SalesReportUpdateParam{StartAt: null.TimeFrom(time.Now())},
		entity.SalesReportParam{FileUrl: salesReport.FileUrl},
	)
	if err != nil {
		return err
	}

	resp, err := http.Get(salesReport.FileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	var reports []entity.Report

	fileExt := files.GetExtension(salesReport.FileUrl)

	switch fileExt {
	case "csv":
		reader, err := csv.NewCSV(resp.Body)
		if err != nil {
			return err
		}

		reports = make([]entity.Report, 0)
		errs := reader.Decode(&reports)
		if len(errs) > 0 {
			errMsg := make([]string, len(errs))
			for i, e := range errs {
				errMsg[i] = e.Error()
			}

			return errorPkg.NewWithCode(codes.CodeBadRequest, strings.Join(errMsg, ", "))
		}
	case "xlam", "xlsm", "xlsx", "xltx", "xltm":
		f, err := excelize.OpenReader(resp.Body)
		if err != nil {
			return errorPkg.NewWithCode(codes.CodeInternalServerError, "failed to read Excel file")
		}

		sheetName := f.GetSheetName(0)
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return err
		}

		reports = make([]entity.Report, 0)
		for _, row := range rows[1:] {
			if len(row) < 9 {
				continue
			}

			quantity, err := strconv.ParseInt(row[4], 10, 64)
			if err != nil {
				continue
			}

			unitPrice, err := strconv.ParseInt(row[5], 10, 64)
			if err != nil {
				continue
			}

			total, err := strconv.ParseInt(row[6], 10, 64)
			if err != nil {
				continue
			}

			report := entity.Report{
				InvoiceId:     row[0],
				Date:          row[1],
				CustomerName:  row[2],
				Item:          row[3],
				Quantity:      quantity,
				UnitPrice:     unitPrice,
				Total:         total,
				Status:        row[7],
				PaymentMethod: row[8],
			}

			reports = append(reports, report)
		}
	default:
		return errorPkg.NewWithCode(codes.CodeBadRequest, "invalid file extension")
	}

	if len(reports) == 0 {
		return errorPkg.NewWithCode(codes.CodeInternalServerError, "no reports found to summarize")
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 3)
	wg.Add(3)

	go func() {
		defer wg.Done()
		salesSummary := s.summarizeSalesSummaries(salesReport.Id, reports)

		errCh <- s.salesSummaryDom.Create(ctx, salesSummary)
	}()

	go func() {
		defer wg.Done()
		productSummary := s.summarizeProductSummaries(salesReport.Id, reports)

		errCh <- s.productSummaryDom.Create(ctx, productSummary)
	}()

	go func() {
		defer wg.Done()

		dailySalesSummary := s.summarizeDailySales(salesReport.Id, reports)

		errCh <- s.dailySalesSummaryDom.Create(ctx, dailySalesSummary)
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			err2 := s.salesReportDom.Update(
				ctx,
				entity.SalesReportUpdateParam{CompletedAt: null.TimeFrom(time.Now()), ErrorMessage: null.StringFrom(err.Error())},
				entity.SalesReportParam{FileUrl: salesReport.FileUrl},
			)

			return err2
		}
	}

	err = s.salesReportDom.Update(
		ctx,
		entity.SalesReportUpdateParam{CompletedAt: null.TimeFrom(time.Now())},
		entity.SalesReportParam{FileUrl: salesReport.FileUrl},
	)

	return nil
}

func (s *salesReport) summarizeSalesSummaries(reportId int64, reports []entity.Report) entity.SalesSummary {
	var salesSummary entity.SalesSummary
	var revenue int64
	mapPaymentMethod := make(map[string]int64)

	for _, report := range reports {
		if strings.ToUpper(report.Status) == "SUCCESS" {
			salesSummary.Success++
			revenue = revenue + report.Total
		} else {
			salesSummary.Failed++
		}

		mapPaymentMethod[report.PaymentMethod]++
	}

	var mostUsed string
	var maxCount int64
	for method, count := range mapPaymentMethod {
		if count > maxCount {
			maxCount = count
			mostUsed = method
		}
	}

	salesSummary.TotalTransaction = int64(len(reports))
	salesSummary.TotalRevenue = revenue
	salesSummary.MostPaymentMethod = mostUsed
	salesSummary.ReportId = reportId

	return salesSummary
}

func (s *salesReport) summarizeProductSummaries(reportId int64, reports []entity.Report) []entity.ProductSummary {
	var productSummaries []entity.ProductSummary

	mapProductQuantity := make(map[string]int64)
	mapProductTotalPrice := make(map[string]int64)

	for _, report := range reports {
		if strings.ToUpper(report.Status) == "SUCCESS" {
			mapProductQuantity[report.Item] = mapProductQuantity[report.Item] + report.Quantity
			mapProductTotalPrice[report.Item] = mapProductTotalPrice[report.Item] + report.Total
		}
	}

	for item, quantity := range mapProductQuantity {
		productSummary := entity.ProductSummary{
			ReportId:    reportId,
			ProductName: item,
			Quantity:    quantity,
			Revenue:     mapProductTotalPrice[item],
		}

		productSummaries = append(productSummaries, productSummary)
	}

	return productSummaries
}

func (s *salesReport) summarizeDailySales(reportId int64, reports []entity.Report) []entity.DailySalesSummary {
	var dailySalesSummaries []entity.DailySalesSummary

	mapDateTransaction := make(map[string]int64)
	mapDateRevenue := make(map[string]int64)

	for _, report := range reports {
		if strings.ToUpper(report.Status) == "SUCCESS" {
			mapDateTransaction[report.Date]++
			mapDateRevenue[report.Date] = mapDateRevenue[report.Date] + report.Total
		}
	}

	for dateString, count := range mapDateTransaction {
		date, _ := time.Parse("2006-01-02", dateString)
		dailySaleSummary := entity.DailySalesSummary{
			ReportId:         reportId,
			Date:             date,
			TotalTransaction: count,
			TotalRevenue:     mapDateRevenue[dateString],
		}

		dailySalesSummaries = append(dailySalesSummaries, dailySaleSummary)
	}

	return dailySalesSummaries
}
