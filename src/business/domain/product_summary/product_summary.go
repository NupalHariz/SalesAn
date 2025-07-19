package productsummary

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

type Interface interface {
	Create(ctx context.Context, param []entity.ProductSummary) error
	GetList(ctx context.Context, param entity.ProductSummaryParam) ([]entity.ProductSummary, error)
}

type productSummary struct {
	db  sql.Interface
	log log.Interface
}

type InitParam struct {
	Db  sql.Interface
	Log log.Interface
}

func Init(param InitParam) Interface {
	return &productSummary{
		db:  param.Db,
		log: param.Log,
	}
}

func (p *productSummary) Create(ctx context.Context, param []entity.ProductSummary) error {
	err := p.createSQL(ctx, param)
	if err != nil {
		return err
	}

	return nil
}

func (p *productSummary) GetList(ctx context.Context, param entity.ProductSummaryParam) ([]entity.ProductSummary, error){
	fmt.Println("\n\n\nPARAM: ", param)

	productSummaries, err := p.getListSQL(ctx, param)
	if err != nil {
		return productSummaries, err
	}

	return productSummaries, nil
}
