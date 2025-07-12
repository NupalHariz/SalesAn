package usecase

import (
	"github.com/NupalHariz/SalesAn/src/business/domain"
	"github.com/NupalHariz/SalesAn/src/business/usecase/user"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/hash"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
)

type Usecases struct {
	User user.Interface
}

type InitParam struct {
	Dom  *domain.Domains
	Json parser.JSONInterface
	Log  log.Interface
	Hash hash.Interface
	Auth auth.Interface
}

func Init(param InitParam) *Usecases {
	return &Usecases{
		User: user.Init(user.InitParam{UserDomain: param.Dom.User, Auth: param.Auth, Hash: param.Hash}),
	}
}
