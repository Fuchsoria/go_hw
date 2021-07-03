package scheduler

import (
	"context"
	"fmt"
	"time"

	simpleconsumer "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/amqp/consumer"
	simpleproducer "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/amqp/producer"
	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/streadway/amqp"
)

type Scheduler struct {
	recheckDelay int64
	logger       *logger.Logger
	storage      *sqlstorage.Storage
	ticker       *time.Ticker
	amqpURI      string
	amqpName     string
}

func NewScheduler(logger *logger.Logger, storage *sqlstorage.Storage, recheckDelay int64, amqpURI string, amqpName string) *Scheduler {
	return &Scheduler{
		recheckDelay: recheckDelay,
		logger:       logger,
		storage:      storage,
		amqpURI:      amqpURI,
		amqpName:     amqpName,
	}
}

func (s *Scheduler) CleanOldEvents() {
	current := time.Now()
	yearAgo := current.AddDate(-1, 0, 0)

	err := s.storage.RemoveOldEvents(yearAgo.Unix())
	if err != nil {
		s.logger.Error(fmt.Errorf("cannot clean old events, %w", err).Error())
	}
}

func (s *Scheduler) SendNotificationsToEvents(producer *simpleproducer.Producer) {
	events, err := s.storage.GetEventsForNotification()
	if err != nil {
		s.logger.Error(fmt.Errorf("cannot get events for notification, %w", err).Error())
	}

	for _, event := range events {
		message := simpleconsumer.AMQPMessage{
			EventID:    event.ID,
			EventTitle: event.Title,
			EventDate:  event.Date,
			User:       event.OwnerID,
		}

		err := producer.Publish(message)
		if err != nil {
			s.logger.Error(fmt.Errorf("cannot publish message, %w", err).Error())
		}
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.ticker = time.NewTicker(time.Duration(s.recheckDelay) * time.Second)
	defer s.ticker.Stop()

	conn, err := amqp.Dial(s.amqpURI)
	if err != nil {
		return fmt.Errorf("cannot connect to amqp, %w", err)
	}

	producer := simpleproducer.New(s.amqpName, conn)
	err = producer.Connect()
	if err != nil {
		return fmt.Errorf("cannot connect to amqp producer, %w", err)
	}

	go func(producer *simpleproducer.Producer) {
		for range s.ticker.C {
			s.SendNotificationsToEvents(producer)
			s.CleanOldEvents()
		}
	}(producer)

	<-ctx.Done()

	return nil
}
