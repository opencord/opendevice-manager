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

// Package contains version v1 of Device Info
package v1

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/opencord/device-management-interface/go/dmi"
)

type DeviceRecordV1_0 struct {
	Uuid         string               `json:"uuid,omitempty"`
	Name         string               `json:"name,omitempty"`
	Make         string               `json:"make,omitempty"`
	Class        dmi.ComponentType    `json:"class,omitempty"`
	Parent       *dmi.Component       `json:"parent,omitempty"`
	ParentRelPos int32                `json:"parent_rel_pos,omitempty"`
	Alias        string               `json:"alias,omitempty"`
	AssetId      string               `json:"asset_id,omitempty"`
	Uri          string               `json:"uri,omitempty"`
	HardwareRev  string               `json:"hardware_rev,omitempty"`
	FirmwareRev  string               `json:"firmware_rev,omitempty"`
	SoftwareRev  string               `json:"software_rev,omitempty"`
	SerialNum    string               `json:"serial_num,omitempty"`
	ModelName    string               `json:"model_name,omitempty"`
	MfgName      string               `json:"mfg_name,omitempty"`
	MfgDate      *timestamp.Timestamp `json:"mfg_date,omitempty"`
	State        *dmi.ComponentState  `json:"state,omitempty"`
	Inventories  map[string]string    `json:"inventories,omitempty"`
	Children     []string             `json:"children,omitempty"` // Children stores uuid of all direct child
	Logging      LoggingInfo          `json:"logging,omitempty"`
	LastChange   *timestamp.Timestamp `json:"last_change,omitempty"`
	LastBooted   *timestamp.Timestamp `json:"last_booted,omitempty"` // Timestamp at which the hardware last booted
}

type LoggingInfo struct {
	EndPoint         string                  `json:"end_point,omitempty"`
	Protocol         string                  `json:"protocol,omitempty"`
	LogLevel         dmi.LogLevel            `json:"log_level,omitempty"`
	LoggableEntities map[string]dmi.LogLevel `json:"loggable_entities,omitempty"`
}
