package dailysalessummary

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
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
