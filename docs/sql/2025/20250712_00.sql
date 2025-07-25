DROP TABLE IF EXISTS sales_reports;

CREATE TABLE IF NOT EXISTS sales_reports(
    id SERIAL PRIMARY KEY,
    file_url VARCHAR(255) NOT NULL,
    start_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message VARCHAR(255)
);

DROP TABLE IF EXISTS sales_summaries;

CREATE TABLE IF NOT EXISTS sales_summaries(
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL,
    total_transaction INT NOT NULL,
    success INT NOT NULL, 
    failed INT NOT NULL,
    total_revenue BIGINT NOT NULL,
    most_payment_method VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS product_summaries;

CREATE TABLE IF NOT EXISTS product_summaries(
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    revenue BIGINT NOT NULL
);

DROP TABLE IF EXISTS daily_sales_summaries;

CREATE TABLE IF NOT EXISTS daily_sales_summaries(
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL,
    date DATE NOT NULL,
    total_transaction INT NOT NULL,
    total_revenue BIGINT NOT NULL
);