package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/query"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

func (u *user) createSQL(ctx context.Context, inputParam entity.UserInputParam) (entity.User, error) {
	user := entity.User{}

	u.log.Debug(ctx, fmt.Sprintf("create user with body: %v", inputParam))

	tx, err := u.db.Leader().BeginTx(ctx, "txUser", sql.TxOptions{})
	if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamed("iNewUser", insertUser)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLPrepareStmt, err.Error())
	}
	defer stmt.Close()

	err = stmt.Get(&user, inputParam)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	if err := tx.Commit(); err != nil {
		return user, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	u.log.Debug(ctx, fmt.Sprintf("success create user with body: %v", inputParam))

	return user, nil
}

func (u *user) getSQL(ctx context.Context, param entity.UserParam) (entity.User, error) {
	user := entity.User{}

	u.log.Debug(ctx, fmt.Sprintf("get user with body: %v", param))

	param.QueryOption.DisableLimit = true
	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &param.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&param)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := u.db.Follower().QueryRow(ctx, "rUser", readUser+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return user, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	if err := row.StructScan(&user); err != nil && errors.Is(err, sql.ErrNotFound) {
		return user, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	} else if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	}

	u.log.Debug(ctx, fmt.Sprintf("success get user with body: %v", param))

	return user, nil
}

func (u *user) getListSQL(ctx context.Context, param entity.UserParam) ([]entity.User, *entity.Pagination, error) {
	users := []entity.User{}

	u.log.Debug(ctx, fmt.Sprintf("get user list with body: %v", param))

	qb := query.NewSQLQueryBuilder(u.db , "param", "db", &param.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&param)
	if err != nil {
		return users, nil, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := u.db.Follower().Query(ctx, "rUserList", readUser+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return users, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		user := entity.User{}
		err := rows.StructScan(&user)
		if err != nil {
			return users, nil, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
		}

		users = append(users, user)
	}

	pg := entity.Pagination{
		CurrentPage:     param.PaginationParam.Page,
		CurrentElements: int64(len(users)),
		SortBy:          param.SortBy,
	}

	if !param.QueryOption.DisableLimit && len(users) > 0 && param.IncludePagination {
		err := u.db.Follower().Get(ctx, "cUserList", countUser+countExt, &pg.TotalElements, countArgs...)
		if err != nil {
			return users, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(param.Limit)

	u.log.Debug(ctx, fmt.Sprintf("success get user list with body: %v", param))

	return users, &pg, nil
}

func (u *user) updateSQL(ctx context.Context, updateParam entity.UserUpdateParam, selectParam entity.UserParam) error {
	u.log.Debug(ctx, fmt.Sprintf("update user %v with body: %v", selectParam.ID, updateParam))

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &selectParam.QueryOption)
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	tx, err := u.db.Leader().BeginTx(ctx, "txUser", sql.TxOptions{})
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	res, err := tx.Exec("uUser", updateUser+queryUpdate, args...)
	if err != nil && strings.Contains(err.Error(), entity.DuplicateEntryErrMessage) {
		return errors.NewWithCode(codes.CodeSQLUniqueConstraint, err.Error())
	} else if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	} else if rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no user updated")
	}

	if err := tx.Commit(); err != nil {
		return errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	u.log.Debug(ctx, fmt.Sprintf("success update user %v with body: %v", selectParam.ID, updateParam))

	return nil
}
