/*
 * Copyright 2020-present Open Networking Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package msgbus holds messagebus related util functions
package msgbus

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/opencord/device-management-interface/go/dmi"

	"github.com/Shopify/sarama"
	"github.com/opencord/opendevice-manager/pkg/config"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

var kafkaProducer sarama.SyncProducer

// logger represents the log object
var logger log.CLogger

// init function for the package
func init() {
	logger = config.Initlog()
}

// InitMsgbusProducer initialises producer for kafka msgbus
func InitMsgbusProducer(ctx context.Context) error {
	cf := config.NewCoreFlags()
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 6
	saramaConfig.Producer.Retry.Backoff = time.Millisecond * 30
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	// The level of acknowledgement reliability needed from the broker.
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{cf.MsgbusEndPoint}
	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)

	if err != nil {
		logger.Errorw(ctx, "Failed-creating-kafka-producer", log.Fields{"error": err, "sarama-config": saramaConfig})
		return err
	}

	kafkaProducer = producer
	logger.Infow(ctx, "creating-kafka-producer-successful", log.Fields{"sarama-config": saramaConfig})
	return nil
}

// SendEvent sends events over kafka bus
func SendEvent(ctx context.Context, event *dmi.Event) error {
	e, err := json.Marshal(event)
	if err != nil {
		logger.Errorw(ctx, "marshal-event-failed", log.Fields{"event": event})
		return err
	}
	logger.Infow(ctx, "sending-event", log.Fields{"event": event})
	return sendMsg(ctx, string(e), config.OpenDevMgrEventsTopic, event.EventId.String())
}

// SendMetric sends metrics over kafka bus
func SendMetric(ctx context.Context, metric *dmi.Metric) error {
	e, err := json.Marshal(metric)
	if err != nil {
		logger.Errorw(ctx, "marshal-metrics-failed", log.Fields{"metrics": metric})
		return err
	}
	logger.Infow(ctx, "sending-metric", log.Fields{"metrics": metric})
	return sendMsg(ctx, string(e), config.OpenDevMgrMetricsTopic, metric.MetricId.String())
}

// SendMsg function will help to publish the message to msgbus/kafka
func sendMsg(ctx context.Context, msg, topic, key string) error {
	if kafkaProducer != nil {
		logger.Debugw(ctx, "sending-message", log.Fields{"msg": msg})
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.StringEncoder(msg),
		}

		partition, offset, err := kafkaProducer.SendMessage(msg)
		logger.Debugw(ctx, "kafka-msg-sent-info", log.Fields{"msg": msg, "partition": partition, "offset": offset, "error": err})
		return err
	}
	logger.Errorw(ctx, "kafka-producer-not-found", log.Fields{"msg": msg, "topic": topic, "key": key})
	return errors.New("kafka producer not found")
}

// Close close the msgbus connection
func Close(ctx context.Context) {
	if kafkaProducer != nil {
		reason := "pod exited"
		logger.Warnw(ctx, "Exiting-msg-bus", log.Fields{"reason": reason})
		kafkaProducer.Close()
	}
}
