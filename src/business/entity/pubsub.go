package entity

import "encoding/json"

const (
	ExchangeSalesReport = "summary-report"
)

const (
	QueueSalesReport = "salesan.summary-report"
	QueueHi          = "salesan.hi"
)

const (
	KeySalesReport = "report.summarize"
	KeyHi          = "hi.name"
)

type PubSubMessage struct {
	Event   string
	Payload json.RawMessage
}
