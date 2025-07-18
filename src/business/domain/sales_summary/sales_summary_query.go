package salessumarry

const (
	insertSalesSummary = `
		INSERT INTO
			sales_summaries(
				report_id,
				total_transaction,
				success,
				failed,
				total_revenue,
				most_payment_method
			)
		VALUES(
			:report_id,
			:total_transaction,
			:success,
			:failed,
			:total_revenue,
			:most_payment_method
		)
	`
)
