package rest

import (
	"github.com/NupalHariz/SalesAn/src/business/dto"
	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
)

func (r *rest) UploadReport(ctx *gin.Context) {
	var param dto.UploadReportParam

	if err := r.Bind(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	data, err := r.uc.SalesReport.UploadReport(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeAccepted, data, nil)
}

func (r *rest) GetReportList(ctx *gin.Context) {
	data, err := r.uc.SalesReport.ListReport(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, data, nil)
}

func (r *rest) GetSummarize(ctx *gin.Context){
	var param dto.ReportParam
	err := r.BindUri(ctx, &param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	res, err := r.uc.SalesReport.GetSummarizeReport(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, res, nil)
}
