package transport

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/events"
	sharedcqrs "github.com/pawverse/pawcare-core/pkg/watermill/cqrs"
	"github.com/pawverse/pawcare-core/pkg/watermill/log"
	"github.com/pawverse/pawcare-core/pkg/watermill/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewEventBus(viper viper.Viper, logger *zap.Logger) (*cqrs.EventBus, error) {
	publisherConfig := kafka.PublisherConfig{
		Brokers:   viper.GetStringSlice(common.KafkaBrokersKey),
		Marshaler: router.NewDefaultPartitionKeyMarshaler(),
	}

	publisher, err := kafka.NewPublisher(publisherConfig, log.NewLogger(logger))
	if err != nil {
		return nil, err
	}

	eventBusConfig := cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return "profiles." + params.EventName, nil
		},

		OnPublish: func(params cqrs.OnEventSendParams) error {
			middleware.SetCorrelationID(watermill.NewUUID(), params.Message)

			if event, ok := params.Event.(events.IEvent); ok {
				key := event.Key()
				params.Message.SetContext(context.WithValue(params.Message.Context(), "key", key))
				params.Message.Metadata["key"] = key
			}

			return nil
		},
		Marshaler: cqrs.JSONMarshaler{
			GenerateName: sharedcqrs.StructNameDotLower,
		},
		Logger: log.NewLogger(logger.With(zap.String("component", "eventBus"))),
	}

	eventBus, err := cqrs.NewEventBusWithConfig(publisher, eventBusConfig)
	if err != nil {
		return nil, err
	}

	return eventBus, nil
}
