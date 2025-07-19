package rest

import (
	"github.com/NupalHariz/SalesAn/src/business/dto"
	"github.com/gin-gonic/gin"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
)

// @Summary Upload Report
// @Description Upload a report file (CSV or Excel) to generate a summary
// @Tags Sales Report
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Upload report file"
// @Success 200 {object} entity.HTTPResp{}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/sales-report [POST]
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

// @Summary Get Report List
// @Description get all report that will be summarized
// @Tags Sales Report
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=[]dto.GetReportList}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/sales-report [GET]
func (r *rest) GetReportList(ctx *gin.Context) {
	data, err := r.uc.SalesReport.ListReport(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, data, nil)
}

// @Summary Get Report Summary
// @Description get summary of specified report
// @Tags Sales Report
// @Security BearerAuth
// @Produce json
// @Param report_id path string true "report id"
// @Success 200 {object} entity.HTTPResp{data=dto.SummaryReport}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 404 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/sales-report/{report_id} [GET]
func (r *rest) GetSummary(ctx *gin.Context) {
	var param dto.ReportParam
	err := r.BindUri(ctx, &param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	res, err := r.uc.SalesReport.GetSummaryReport(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, res, nil)
}
