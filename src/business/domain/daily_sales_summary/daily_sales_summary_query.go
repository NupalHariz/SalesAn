package dailysalessummary

const (
	insertDailySalesSummary = `
		INSERT INTO
			daily_sales_summaries(
				report_id,
				date,
				total_transaction,
				total_revenue
			) VALUES(
			:report_id,
			:date,
			:total_transaction,
			:total_revenue	 
		)
	`

	readDailySalesSummary = `
		SELECT
			id,
			report_id,
			date,
			total_transaction,
			total_revenue
		FROM
		 	daily_sales_summaries
	`
)
