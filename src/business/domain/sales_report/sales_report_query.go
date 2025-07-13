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
)
