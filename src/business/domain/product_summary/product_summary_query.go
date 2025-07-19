package productsummary

const (
	insertProductSummary = `
		INSERT INTO
			product_summaries(
				report_id,
				product_name,
				quantity,
				revenue
		) VALUES(
			:report_id,
			:product_name,
			:quantity,
			:revenue	 
		)
	`

	readProductSummary = `
		SELECT
			id,
			report_id,
			product_name,
			quantity,
			revenue
		FROM
		 	product_summaries
	`
)
