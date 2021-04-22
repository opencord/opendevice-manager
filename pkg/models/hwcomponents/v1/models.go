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

// Package v1 stores models for hw components
package v1

import (
	dmi "github.com/opencord/device-management-interface/go/dmi"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

type HwCompRecordV1_0 struct {
	Name         string                     `json:"name,omitempty"`
	Class        dmi.ComponentType          `json:"class,omitempty"`
	Description  string                     `json:"description,omitempty"`
	Parent       string                     `json:"parent,omitempty"`
	ParentRelPos int32                      `json:"parent_rel_pos,omitempty"`
	Children     []string                   `json:"children,omitempty"` // Children stores uuid of all direct child
	SerialNum    string                     `json:"serial_num,omitempty"`
	MfgName      string                     `json:"mfg_name,omitempty"`
	ModelName    string                     `json:"model_name,omitempty"`
	Alias        string                     `json:"alias,omitempty"`
	AssetId      string                     `json:"asset_id,omitempty"`
	IsFru        bool                       `json:"is_fru,omitempty"`
	MfgDate      *timestamp.Timestamp       `json:"mfg_date,omitempty"`
	Uri          string                     `json:"uri,omitempty"`
	Uuid         string                     `json:"uuid,omitempty"`
	State        *dmi.ComponentState        `json:"state,omitempty"`
	SensorData   []*dmi.ComponentSensorData `json:"sensor_data,omitempty"`
	Specific     string                     `json:"specific,omitempty"`
}
