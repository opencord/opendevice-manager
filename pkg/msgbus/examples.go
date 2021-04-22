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

	"github.com/opencord/device-management-interface/go/dmi"
)

// SendSampleMetric sends a sample metric
func SendSampleMetric(ctx context.Context) {
	m := new(dmi.Metric)
	m.MetricId = dmi.MetricNames_METRIC_CPU_TEMP
	m.MetricMetadata = new(dmi.MetricMetaData)
	m.MetricMetadata.ComponentName = "CPU-COMPONENT"
	m.MetricMetadata.ComponentUuid = new(dmi.Uuid)
	m.MetricMetadata.ComponentUuid.Uuid = "uuid-123"
	m.MetricMetadata.DeviceUuid = new(dmi.Uuid)
	m.MetricMetadata.DeviceUuid.Uuid = "dev-uuid-123"
	go SendMetric(ctx, m)
}

// SendSampleEvent sends a sample event
func SendSampleEvent(ctx context.Context) {
	e := new(dmi.Event)
	e.EventId = dmi.EventIds_EVENT_FAN_FAILURE
	e.EventMetadata = new(dmi.EventMetaData)
	e.EventMetadata.ComponentName = "CPU-COMPONENT"
	e.EventMetadata.ComponentUuid = new(dmi.Uuid)
	e.EventMetadata.ComponentUuid.Uuid = "uuid-123"
	go SendEvent(ctx, e)
}
