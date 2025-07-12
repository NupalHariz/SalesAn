DROP TYPE IF EXISTS sales_report_status;

CREATE TYPE sales_report_status AS ENUM('WAITING', 'PROCESSING', 'DONE', 'FAILED');

DROP TABLE IF EXISTS sales_reports;

CREATE TABLE IF NOT EXISTS sales_reports(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    status sales_report_status DEFAULT('WAITING'),
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
    total_revenue INT NOT NULL,
    most_payment_method VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS product_summaries;

CREATE TABLE IF NOT EXISTS product_summaries(
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    revenue INT NOT NULL
);

DROP TABLE IF EXISTS daily_sales_summaries;

CREATE TABLE IF NOT EXISTS daily_sales_summaries(
    id SERIAL PRIMARY KEY,
    report_id INT NOT NULL,
    date DATE NOT NULL,
    total_transaction INT NOT NULL,
    total_revenue INT NOT NULL
);