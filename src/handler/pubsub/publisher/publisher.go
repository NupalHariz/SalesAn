package publisher

import (
	"context"
	"fmt"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/NupalHariz/SalesAn/src/utils/config"
	"github.com/rabbitmq/amqp091-go"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
)

type Interface interface {
	Publish(ctx context.Context, exchaneName string, routingKey string, data interface{}) error
}

type publisher struct {
	pbsb *amqp091.Connection
	log  log.Interface
	json parser.JSONInterface
}

type InitParam struct {
	Cfg  config.RabbitMQConfig
	Log  log.Interface
	Json parser.JSONInterface
}

func Init(param InitParam) Interface {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		param.Cfg.Username,
		param.Cfg.Password,
		param.Cfg.Host,
		param.Cfg.Port,
	)

	conn, err := amqp091.Dial(connAddr)
	if err != nil {
		param.Log.Fatal(context.Background(), fmt.Sprintf("failed to create connection, err: %s", err))
	}

	notifyClose := conn.NotifyClose(make(chan *amqp091.Error))
	go func() {
		err := <-notifyClose
		if err != nil {
			conn.Close()
		}
	}()

	param.Log.Info(context.Background(), fmt.Sprintf("success to create publisher rabbitmq connection, conn: %s", connAddr))

	return &publisher{
		pbsb: conn,
		log:  param.Log,
		json: param.Json,
	}
}

func (p *publisher) Publish(ctx context.Context, exchangeName string, routingKey string, data interface{}) error {
	ch, err := p.pbsb.Channel()
	if err != nil {
		p.log.Error(context.Background(), fmt.Sprintf("failed to open channel, err: %s", err))
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		p.log.Error(context.Background(), fmt.Sprintf("failed to declare exchange %s, err: %s", exchangeName, err))
		return err
	}

	jsonData, err := p.json.Marshal(data)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("error unmarshal: %v", err))
		return err
	}

	pubsubMsg := entity.PubSubMessage{
		Event:   fmt.Sprintf("%s:%s", exchangeName, routingKey),
		Payload: jsonData,
	}

	payload, err := p.json.Marshal(pubsubMsg)
	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("error unmarshal: %v", err))
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			Body: payload,
		},
	)

	if err != nil {
		p.log.Error(ctx, fmt.Sprintf("failed to publish message %v, err : %s", pubsubMsg, err))
		return errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	p.log.Info(ctx, fmt.Sprintf("success to publish message with body: %v", payload))

	return nil
}
