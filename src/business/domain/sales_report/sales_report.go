package salesreport

import (
	"context"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

type Interface interface {
	Create(ctx context.Context, param entity.SalesReport) (entity.SalesReport, error)
	GetList(ctx context.Context) ([]entity.SalesReport, error)
	Update(ctx context.Context, updateParam entity.SalesReportUpdateParam, saleReportParam entity.SalesReportParam) error
	Get(ctx context.Context, param entity.SalesReportParam) (entity.SalesReport, error)
}

type salesReport struct {
	db  sql.Interface
	log log.Interface
}

type InitParam struct {
	Db  sql.Interface
	Log log.Interface
}

func Init(param InitParam) Interface {
	return &salesReport{
		db:  param.Db,
		log: param.Log,
	}
}

func (s *salesReport) Create(ctx context.Context, param entity.SalesReport) (entity.SalesReport, error) {
	salesReport, err := s.createSQL(ctx, param)
	if err != nil {
		return salesReport, err
	}

	return salesReport, nil
}

func (s *salesReport) GetList(ctx context.Context) ([]entity.SalesReport, error) {
	salesReports, err := s.getListSQL(ctx)
	if err != nil {
		return salesReports, err
	}

	return salesReports, nil
}

func (s *salesReport) Update(ctx context.Context, updateParam entity.SalesReportUpdateParam, saleReportParam entity.SalesReportParam) error {
	err := s.updateSQL(ctx, updateParam, saleReportParam)
	if err != nil {
		return err
	}

	return nil
}

func (s *salesReport) Get(ctx context.Context, param entity.SalesReportParam) (entity.SalesReport, error) {
	salesReport, err := s.getSQL(ctx, param)
	if err != nil {
		return salesReport, err
	}

	return salesReport, nil
}
