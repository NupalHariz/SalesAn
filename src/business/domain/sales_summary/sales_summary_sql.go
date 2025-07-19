package salessumarry

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

func (s *salesSummary) createSQL(ctx context.Context, param entity.SalesSummary) error {
	s.log.Info(ctx, fmt.Sprintf("insert to sales_summaries with param: %v", param))

	tx, err := s.db.Leader().BeginTx(ctx, "txSalesSummary", sql.TxOptions{})
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	res, err := tx.NamedExec("iNewSalesSummary", insertSalesSummary, param)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no sales_summary created")
	}

	if err = tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	s.log.Info(ctx, fmt.Sprintf("success to create sales summary with param: %v", param))

	return nil
}

func (s *salesSummary) getSQL(ctx context.Context, param entity.SalesSummaryParam) (entity.SalesSummary, error) {
	var salesSummary entity.SalesSummary

	s.log.Debug(ctx, fmt.Sprintf("get sales summary with param %v", param))

	qb := query.NewSQLQueryBuilder(s.db, "param", "db", &query.Option{})
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return salesSummary, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := s.db.QueryRow(ctx, "rSalesSummary", readSalesSummary+queryExt, queryArgs...)
	if err != nil {
		return salesSummary, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&salesSummary); err != nil && errors.Is(sql.ErrNotFound, err) {
		return salesSummary, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return salesSummary, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	s.log.Debug(ctx, fmt.Sprintf("success to get sales summary with param %v", param))

	return salesSummary, nil
}
