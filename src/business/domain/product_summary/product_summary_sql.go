package productsummary

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

func (p *productSummary) createSQL(ctx context.Context, param []entity.ProductSummary) error {
	p.log.Info(ctx, fmt.Sprintf("insert to product_summaries with param: %v", param))

	tx, err := p.db.Leader().BeginTx(ctx, "txProductSummary", sql.TxOptions{})
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	res, err := tx.NamedExec("iNewProductSummary", insertProductSummary, param)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no product summary created")
	}

	if err = tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	p.log.Info(ctx, fmt.Sprintf("success to create product summary with param: %v", param))

	return nil
}

func (p *productSummary) getListSQL(ctx context.Context, param entity.ProductSummaryParam) ([]entity.ProductSummary, error) {
	var productSummaries []entity.ProductSummary

	p.log.Debug(ctx, fmt.Sprintf("read product summary list with param %v", param))

	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &query.Option{})
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return productSummaries, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := p.db.Query(ctx, "raProductSummary", readProductSummary+queryExt, queryArgs...)
	if err != nil {
		return productSummaries, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var productSummary entity.ProductSummary
		err := rows.StructScan(&productSummary)
		if err != nil {
			p.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		productSummaries = append(productSummaries, productSummary)
	}

	return productSummaries, nil
}
