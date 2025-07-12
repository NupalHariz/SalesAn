package user

import (
	"context"
	"fmt"
	"time"

	userDomain "github.com/NupalHariz/SalesAn/src/business/domain/user"
	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/hash"
	"github.com/reyhanmichiels/go-pkg/v2/null"
	"github.com/reyhanmichiels/go-pkg/v2/query"
)

var Now = time.Now

type Interface interface {
	Register(ctx context.Context, inputParam entity.UserInputParam) (entity.User, error)
	SignIn(ctx context.Context, param entity.UserLoginParam) (entity.UserLoginResponse, error)
	Get(ctx context.Context, param entity.UserParam) (entity.User, error)
	RefreshToken(ctx context.Context, param entity.RefreshTokenParam) (entity.UserLoginResponse, error)
}

type user struct {
	user userDomain.Interface
	auth auth.Interface
	hash hash.Interface
}

type InitParam struct {
	UserDomain userDomain.Interface
	Auth       auth.Interface
	Hash       hash.Interface
}

func Init(param InitParam) Interface {
	return &user{
		user: param.UserDomain,
		auth: param.Auth,
		hash: param.Hash,
	}
}

func (u *user) Register(ctx context.Context, inputParam entity.UserInputParam) (entity.User, error) {
	user := entity.User{}

	if inputParam.Password != inputParam.ConfirmPassword {
		return user, errors.NewWithCode(codes.CodeBadRequest, "confirmation password failed")
	}

	isUserExist := false
	_, err := u.user.Get(ctx, entity.UserParam{
		Email: inputParam.Email,
		QueryOption: query.Option{
			IsActive: true,
		},
	})
	if err != nil && errors.GetCode(err) != codes.CodeSQLRecordDoesNotExist {
		return user, err
	} else if err == nil {
		isUserExist = true
	}

	if isUserExist {
		return user, errors.NewWithCode(codes.CodeConflict, "email already used")
	}

	hashedPassword, err := u.hash.Bcrypt().GenerateFromText(inputParam.Password)
	if err != nil {
		return user, err
	}

	inputParam.CreatedAt = null.TimeFrom(Now())
	inputParam.Password = hashedPassword
	user, err = u.user.Create(ctx, inputParam)
	if err != nil {
		return user, err
	}

	fmt.Printf("\n\n\nType of RefreshToken: %T\n\n\n", user.RefreshToken)

	return user, nil
}

func (u *user) SignIn(ctx context.Context, param entity.UserLoginParam) (entity.UserLoginResponse, error) {
	userLoginResponse := entity.UserLoginResponse{}

	user, err := u.user.Get(ctx, entity.UserParam{
		Email: param.Email,
		QueryOption: query.Option{
			IsActive: true,
		},
	})
	if err != nil && errors.GetCode(err) == codes.CodeSQLRecordDoesNotExist {
		return userLoginResponse, errors.NewWithCode(codes.CodeUnauthorized, "invalid email or password")
	} else if err != nil && errors.GetCode(err) != codes.CodeSQLRecordDoesNotExist {
		return userLoginResponse, err
	}

	isPasswordSame := u.hash.Bcrypt().CompareHashWithText(user.Password, param.Password)
	if !isPasswordSame {
		return userLoginResponse, errors.NewWithCode(codes.CodeUnauthorized, "invalid email or password")
	}

	accessToken, refreshToken, err := u.issueToken(ctx, user.ID)
	if err != nil {
		return userLoginResponse, err
	}

	userLoginResponse = entity.UserLoginResponse{
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return userLoginResponse, nil
}

func (u *user) Get(ctx context.Context, param entity.UserParam) (entity.User, error) {
	param.QueryOption.IsActive = true
	user, err := u.user.Get(ctx, param)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *user) RefreshToken(ctx context.Context, param entity.RefreshTokenParam) (entity.UserLoginResponse, error) {
	response := entity.UserLoginResponse{}

	user, err := u.user.Get(ctx, entity.UserParam{
		RefreshToken: param.RefreshToken,
		QueryOption: query.Option{
			IsActive: true,
		},
	})
	if err != nil && errors.GetCode(err) == codes.CodeSQLRecordDoesNotExist {
		return response, errors.NewWithCode(codes.CodeNotFound, "refresh token not found")
	} else if err != nil && errors.GetCode(err) != codes.CodeSQLRecordDoesNotExist {
		return response, err
	}

	err = u.auth.ValidateRefreshToken(param.RefreshToken)
	if err != nil {
		return response, err
	}

	accessToken, refreshToken, err := u.issueToken(ctx, user.ID)
	if err != nil {
		return response, err
	}

	response = entity.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (u *user) issueToken(ctx context.Context, userID int64) (string, string, error) {
	accessToken, err := u.auth.CreateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := u.auth.CreateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	err = u.user.Update(ctx, entity.UserUpdateParam{
		RefreshToken: refreshToken,
	}, entity.UserParam{
		ID: userID,
		QueryOption: query.Option{
			DisableLimit: true,
		},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}
