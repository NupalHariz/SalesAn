package usecase

import (
	"github.com/NupalHariz/SalesAn/src/business/domain"
	"github.com/NupalHariz/SalesAn/src/business/service/supabase"
	salesreport "github.com/NupalHariz/SalesAn/src/business/usecase/sales_report"
	"github.com/NupalHariz/SalesAn/src/business/usecase/user"
	"github.com/NupalHariz/SalesAn/src/handler/pubsub/publisher"
	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/hash"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
)

type Usecases struct {
	User        user.Interface
	SalesReport salesreport.Interface
}

type InitParam struct {
	Dom       *domain.Domains
	Json      parser.JSONInterface
	Log       log.Interface
	Hash      hash.Interface
	Auth      auth.Interface
	Supabase  supabase.Interface
	Publisher publisher.Interface
}

func Init(param InitParam) *Usecases {
	return &Usecases{
		User:        user.Init(user.InitParam{UserDomain: param.Dom.User, Auth: param.Auth, Hash: param.Hash}),
		SalesReport: salesreport.Init(salesreport.InitParam{SalesReportDom: param.Dom.SalesReport, SalesSummaryDom: param.Dom.SalesSummary, ProductSummaryDom: param.Dom.ProductSummary, DailySalesSummaryDom: param.Dom.DailySalesSummary, Auth: param.Auth, Supabase: param.Supabase, Publisher: param.Publisher, Json: param.Json}),
	}
}
