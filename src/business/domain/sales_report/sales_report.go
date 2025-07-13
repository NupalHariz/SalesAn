package salesreport

import (
	"context"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

type Interface interface {
	Create(ctx context.Context, param entity.SalesReport) error
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

func (s *salesReport) Create(ctx context.Context, param entity.SalesReport) error {
	err := s.createSQL(ctx, param)
	if err != nil {
		return err
	}

	return nil
}
