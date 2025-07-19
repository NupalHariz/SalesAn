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
		) RETURNING *
	`

	readSalesReport = `
		SELECT
			id,
			file_url,
			error_message,
			start_at,
			completed_at
		FROM
			sales_reports
	`

	updateSalesReport = `
		UPDATE
			sales_reports
	`
)