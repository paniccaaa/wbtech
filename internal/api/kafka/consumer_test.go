package kafka

import (
	"context"
	"log/slog"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/paniccaaa/wbtech/internal/api/kafka/mocks"
	"github.com/paniccaaa/wbtech/internal/app"
)

func TestConsumer_processEvent(t *testing.T) {
	type fields struct {
		client   Poller
		cfgKafka app.Kafka
		handler  MessageHandler
		log      *slog.Logger
	}
	type args struct {
		ctx context.Context
		ev  kafka.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid kafka.Message",
			fields: fields{
				handler: mocks.NewMessageHandler(t),
			},
			args: args{
				ctx: context.Background(),
				ev: &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: ptrToString("test-topic")},
					Value:          []byte("test-message"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consumer{
				client:   tt.fields.client,
				cfgKafka: tt.fields.cfgKafka,
				handler:  tt.fields.handler,
				log:      tt.fields.log,
			}

			if tt.args.ev != nil {
				if message, ok := tt.args.ev.(*kafka.Message); ok {
					tt.fields.handler.(*mocks.MessageHandler).
						On("ProcessKafkaMessage", tt.args.ctx, *message.TopicPartition.Topic, message.Value).
						Return(nil).Once()
				}
			}

			if err := c.processEvent(tt.args.ctx, tt.args.ev); (err != nil) != tt.wantErr {
				t.Errorf("Consumer.processEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.ev != nil {
				if _, ok := tt.args.ev.(*kafka.Message); ok {
					tt.fields.handler.(*mocks.MessageHandler).AssertExpectations(t)
				}
			}
		})
	}
}

func ptrToString(s string) *string {
	return &s
}
