package mq

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/matt-abi/abi-lib/dynamic"
	"github.com/matt-abi/abi-lib/json"
	kafka "github.com/segmentio/kafka-go"
)

type kafkaConfig struct {
	Addr      string `json:"addr"`
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	MinBytes  int    `json:"min-bytes"`
	MaxBytes  int    `json:"max-bytes"`
}

type kafkaDriver struct {
	config kafkaConfig
	conn   *kafka.Conn
}

func newKafkaDriver(driver string, config interface{}) (Driver, error) {
	v := &kafkaDriver{}
	dynamic.SetValue(&v.config, config)

	conn, err := kafka.Dial("tcp", v.config.Addr)

	if err != nil {
		return nil, err
	}

	v.conn = conn
	return v, nil
}

func (k *kafkaDriver) Send(topic string, name string, data interface{}) error {
	b, _ := json.Marshal(map[string]interface{}{"name": name, "data": data})
	_, err := k.conn.WriteMessages(&kafka.Message{Value: b, Topic: topic})
	return err
}

func (k *kafkaDriver) On(queue string, fn func(name string, data interface{}) bool) context.CancelFunc {

	ctx, cancel := context.WithCancel(context.Background())

	vs := strings.Split(queue, "/")

	groupId := ""
	topic := vs[0]

	if len(vs) > 1 {
		groupId = vs[1]
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.config.Addr},
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: k.config.MinBytes, // 10KB
		MaxBytes: k.config.MaxBytes, // 10MB
	})

	done := func() error {

		m, err := r.FetchMessage(ctx)

		if err != nil {
			return err
		}

		var data interface{} = nil

		err = json.Unmarshal(m.Value, &data)

		if err != nil {
			log.Println("[kafka]", "[err:1]", err)
			return nil
		} else {
			name := dynamic.StringValue(dynamic.Get(data, "name"), "")
			if name != "" {
				err = fn(name, dynamic.Get(data, "data"))
				if err != nil {
					return err
				}
				err = r.CommitMessages(ctx, m)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	go func() {

		tk := time.NewTicker(30 * time.Second)

		defer tk.Stop()

		running := true

		for running {

			err := done()

			if err != nil {
				log.Println("[kafka]", "[err:1]", err, "Wait 30 seconds and try again")
				select {
				case <-ctx.Done():
					running = false
				case <-tk.C:
				}
			} else {
				select {
				case <-ctx.Done():
					running = false
				default:
				}
			}

		}

	}()

	return cancel
}

func (k *kafkaDriver) Recycle() {
	if k.conn != nil {
		k.conn.Close()
		k.conn = nil
	}
}

func init() {
	Reg("kafka", newKafkaDriver)
}
