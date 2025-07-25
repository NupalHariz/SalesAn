package salesreport

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

func (s *salesReport) createSQL(ctx context.Context, param entity.SalesReport) (entity.SalesReport, error) {
	var salesReport entity.SalesReport
	s.log.Info(ctx, fmt.Sprintf("create sales report with body %v", param))

	stmt, err := s.db.PrepareNamed(ctx, "iNewSalesReport", insertSalesReport)
	if err != nil {
		return salesReport, err
	}
	defer stmt.Close()

	err = stmt.Get(&salesReport, param)

	s.log.Info(ctx, fmt.Sprintf("success to create sales report with body %v", param))

	return salesReport, nil
}

func (s *salesReport) getListSQL(ctx context.Context) ([]entity.SalesReport, error) {
	var salesReports []entity.SalesReport

	s.log.Debug(ctx, "get all report list")

	rows, err := s.db.Query(ctx, "raSalesReport", readSalesReport)
	if err != nil {
		return salesReports, err
	}
	defer rows.Close()

	for rows.Next() {
		var salesReport entity.SalesReport

		err := rows.StructScan(&salesReport)
		if err != nil {
			return salesReports, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
		}

		salesReports = append(salesReports, salesReport)
	}

	s.log.Debug(ctx, "success to get sales report list")

	return salesReports, nil
}

func (s *salesReport) updateSQL(ctx context.Context, updateParam entity.SalesReportUpdateParam, salesReportParam entity.SalesReportParam) error {
	s.log.Debug(ctx, fmt.Sprintf("update sales report with param %v and body %v", updateParam, salesReportParam))

	qb := query.NewSQLQueryBuilder(s.db, "param", "db", &query.Option{})

	queryExt, queryArgs, err := qb.BuildUpdate(&updateParam, &salesReportParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	tx, err := s.db.Leader().BeginTx(ctx, "txSalesReport", sql.TxOptions{})
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	res, err := tx.Exec("uSalesReport", updateSalesReport+queryExt, queryArgs...)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no sales report updated")
	}

	if err = tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	s.log.Debug(ctx, fmt.Sprintf("success to update sales report with param %v and body %v", updateParam, salesReportParam))

	return nil
}

func (s *salesReport) getSQL(ctx context.Context, param entity.SalesReportParam) (entity.SalesReport, error) {
	var salesReport entity.SalesReport

	s.log.Debug(ctx, fmt.Sprintf("get sales summary with param %v", param))

	qb := query.NewSQLQueryBuilder(s.db, "param", "db", &query.Option{})
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return salesReport, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := s.db.QueryRow(ctx, "rSalesReport", readSalesReport+queryExt, queryArgs...)
	if err != nil {
		return salesReport, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&salesReport); err != nil && errors.Is(sql.ErrNotFound, err) {
		return salesReport, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return salesReport, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	s.log.Debug(ctx, fmt.Sprintf("success to get sales summary with param %v", param))

	return salesReport, nil
}
