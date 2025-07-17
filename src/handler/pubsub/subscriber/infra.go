package subscriber

import "github.com/NupalHariz/SalesAn/src/business/entity"

type setUpQueue struct {
	QueueName    string
	RoutingKey   string
	ExchangeName string
}

var setUpQueues = []setUpQueue{
	{
		QueueName:    entity.QueueSalesReport,
		RoutingKey:   entity.KeySalesReport,
		ExchangeName: entity.ExchangeSalesReport,
	},
	{
		QueueName:    entity.QueueHi,
		RoutingKey:   entity.KeyHi,
		ExchangeName: entity.ExchangeSalesReport,
	},
}
