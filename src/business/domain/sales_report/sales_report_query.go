package salesreport

const (
	insertSalesReport = `
		INSERT INTO
			sales_reports(
			user_id,
			file_url
			)
		VALUES(
			:user_id,
			:file_url
		)
	`

	readSalesReportList = `
		SELECT
			id,
			file_url,
			error_message,
			start_at,
			completed_at
		FROM
			sales_reports
	`
)