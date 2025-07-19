package dailysalessummary

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

func (d *dailySalesSummary) createSQL(ctx context.Context, param []entity.DailySalesSummary) error {
	d.log.Info(ctx, fmt.Sprintf("insert to daily sales summary with param: %v", param))

	tx, err := d.db.Leader().BeginTx(ctx, "txDailySalesSummary", sql.TxOptions{})
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	res, err := tx.NamedExec("iNewDailySalesSummary", insertDailySalesSummary, param)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no daily sales summary created")
	}

	if err = tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	d.log.Info(ctx, fmt.Sprintf("success to create daily sales summary with param: %v", param))

	return nil
}

func (p *dailySalesSummary) getListSQL(ctx context.Context, param entity.DailySalesSummaryParam) ([]entity.DailySalesSummary, error) {
	var dailySalesSummaries []entity.DailySalesSummary

	p.log.Debug(ctx, fmt.Sprintf("read daily sales summary list with param %v", param))

	qb := query.NewSQLQueryBuilder(p.db, "param", "db", &query.Option{})
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return dailySalesSummaries, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := p.db.Query(ctx, "raProductSummary", readDailySalesSummary+queryExt, queryArgs...)
	if err != nil {
		return dailySalesSummaries, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}
	defer rows.Close()

	for rows.Next(){
		var dailySalesSummary entity.DailySalesSummary
		err := rows.StructScan(&dailySalesSummary)
		if err != nil {
			p.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}

		dailySalesSummaries = append(dailySalesSummaries, dailySalesSummary)
	}

	return dailySalesSummaries, nil
}

