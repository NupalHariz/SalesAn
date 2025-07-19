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
)
