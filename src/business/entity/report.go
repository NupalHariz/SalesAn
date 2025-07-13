package entity

type Report struct {
	InvoiceId     string `validate:"required"`
	Date          string `validate:"required,len=10"`
	CustomerName  string `validate:"required"`
	Item          string `validate:"required"`
	Quantity      int64  `validate:"required,numeric,gt=0"`
	UnitPrice     int64  `validate:"required,numeric,gt=0"`
	Total         int64  `validate:"required,numeric,gt=0"`
	Status        string `validate:"required"`
	PaymentMethod string `validate:"required"`
}
