package domain

import (
	salesreport "github.com/NupalHariz/SalesAn/src/business/domain/sales_report"
	"github.com/NupalHariz/SalesAn/src/business/domain/user"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/redis"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
)

type Domains struct {
	User        user.Interface
	SalesReport salesreport.Interface
}

type InitParam struct {
	Log   log.Interface
	Db    sql.Interface
	Redis redis.Interface
	Json  parser.JSONInterface
	// TODO: add audit
}

func Init(param InitParam) *Domains {
	return &Domains{
		User:        user.Init(user.InitParam{Db: param.Db, Log: param.Log, Redis: param.Redis, Json: param.Json}),
		SalesReport: salesreport.Init(salesreport.InitParam{Db: param.Db, Log: param.Log}),
	}
}
