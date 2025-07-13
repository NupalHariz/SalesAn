package salesreport

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

func (s *salesReport) createSQL(ctx context.Context, param entity.SalesReport) error {
	s.log.Info(ctx, fmt.Sprintf("create sales report with body %v", param))

	tx, err := s.db.Leader().BeginTx(ctx, "txSalesReport", sql.TxOptions{})
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	res, err := tx.NamedExec("iNewSalesReport", insertSalesReport, param)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no sales report created")
	}

	if err = tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	s.log.Info(ctx, fmt.Sprintf("success to create sales report with body %v", param))

	return nil
}

func (s *salesReport) getListSQL(ctx context.Context) ([]entity.SalesReport, error) {
	var salesReports []entity.SalesReport

	s.log.Debug(ctx, "get all report list")

	rows, err := s.db.Query(ctx, "raSalesReport", readSalesReportList)
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