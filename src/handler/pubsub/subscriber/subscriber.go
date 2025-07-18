package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/NupalHariz/SalesAn/src/business/entity"
	"github.com/NupalHariz/SalesAn/src/business/usecase"
	"github.com/NupalHariz/SalesAn/src/utils/config"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/reyhanmichiels/go-pkg/v2/appcontext"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
)

type callFunc func(ctx context.Context, data entity.PubSubMessage) error

type Interface interface {
	InitSubscription()
}

type subsriber struct {
	pbsb         *amqp091.Connection
	cfg          config.RabbitMQConfig
	log          log.Interface
	json         parser.JSONInterface
	uc           usecase.Usecases
	callsFuncMap map[string]callFunc
}

type InitParam struct {
	Cfg  config.RabbitMQConfig
	Log  log.Interface
	Json parser.JSONInterface
	UC   usecase.Usecases
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

	subsriber := &subsriber{
		pbsb: conn,
		log:  param.Log,
		json: param.Json,
		uc:   param.UC,
		cfg:  param.Cfg,
	}

	subsriber.assignEvent()

	notifyClose := subsriber.pbsb.NotifyClose(make(chan *amqp091.Error))
	go func() {
		err := <-notifyClose
		if err != nil {
			conn.Close()
		}
	}()

	param.Log.Info(context.Background(), fmt.Sprintf("success to create subscriber rabbit mq connection, conn: %s", connAddr))

	return subsriber
}

func (s *subsriber) InitSubscription() {
	go func() {
		for _, q := range setUpQueues {
			err := s.subscribe(context.Background(), q)
			if err != nil {
				s.log.Error(context.Background(), fmt.Sprintf("error occured: %s", err))
			}
		}
	}()
}

func (s *subsriber) assignEvent() {
	s.callsFuncMap = map[string]callFunc{
		fmt.Sprintf("%s:%s", entity.QueueSalesReport, entity.KeySalesReport): s.uc.SalesReport.SummarizeReport,
	}
}

func (s *subsriber) addFieldsToContext(ctx context.Context, eventName string) context.Context {
	ctx = appcontext.SetRequestId(ctx, uuid.NewString())
	ctx = appcontext.SetUserAgent(ctx, eventName)
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())

	return ctx
}

func (s *subsriber) subscribe(ctx context.Context, param setUpQueue) error {
	eventName := fmt.Sprintf("%s:%s", param.QueueName, param.RoutingKey)
	ctx = s.addFieldsToContext(ctx, eventName)

	ch, err := s.pbsb.Channel()
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("failed to open channel, err: %s", err))
	}
	notifyClose := ch.NotifyClose(make(chan *amqp091.Error))

	err = ch.ExchangeDeclare(
		param.ExchangeName,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("failed to declare exchange %s, err: %v", param.ExchangeName, err))
		return err
	}

	q, err := ch.QueueDeclare(
		param.QueueName,
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("failed to declare queue %s, err: %v", param.QueueName, err))
		return err
	}

	err = ch.QueueBind(
		q.Name,
		param.RoutingKey,
		param.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		s.log.Error(ctx, fmt.Sprintf("failed to create bind queue %s, err: %v", q.Name, err))
		return err
	}

	deliveries, err := ch.Consume(
		q.Name,
		"",
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)

	go func() {
		for delivery := range deliveries {
			var data entity.PubSubMessage
			err = s.json.Unmarshal(delivery.Body, &data)
			if err != nil {
				s.log.Error(ctx, fmt.Sprintf("error unmarshall: %v", err))
				continue
			}

			eventKey := fmt.Sprintf("%s:%s", param.QueueName, param.RoutingKey)
			if handler, ok := s.callsFuncMap[eventKey]; ok {
				if err := handler(ctx, data); err != nil {
					s.log.Error(ctx, fmt.Sprintf("handler error: %v", err))
				} else {
					startTime := appcontext.GetRequestStartTime(ctx)
					s.log.Info(ctx, fmt.Sprintf("successfully handled event %s in %v", eventName, time.Since(startTime)))
				}
			} else {
				s.log.Warn(ctx, fmt.Sprintf("no handler for %s", eventKey))
			}

		}
	}()

	go func() {
		err := <-notifyClose
		if err != nil {
			ch.Close()
		}
	}()

	return nil
}
