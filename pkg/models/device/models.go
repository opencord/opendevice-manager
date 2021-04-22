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

// Package device stores methods and functions related to device
package device

import (
	"context"
	"sync"

	"github.com/opencord/device-management-interface/go/dmi"

	"github.com/jinzhu/copier"
	"github.com/opencord/opendevice-manager/pkg/config"

	v1 "github.com/opencord/opendevice-manager/pkg/models/device/v1"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// Constants defined are the DB Path meant for storing device info records
const (
	DbPathUuidToName   = config.DBPrefix + config.CurDBVer + "/DevRec/DevUuid/%s"
	DbPathNameToRecord = config.DBPrefix + config.CurDBVer + "/DevRec/DevName/%s"
)

// deviceCache stores device informations in buffer
type deviceCache struct {
	nameToRec  *sync.Map // nameToRecord maintains cache for mapping from name to main record
	uuidToName *sync.Map // uuidToName maintains cache for mapping from uuid to name
}

var cache *deviceCache

// logger represents the log object
var logger log.CLogger

// initCache initialises device cache
func initCache() {
	cache = new(deviceCache)
	cache.nameToRec = new(sync.Map)
	cache.uuidToName = new(sync.Map)
}

// init function for the package
func init() {
	logger = config.Initlog()
	initCache()
}

// ClearCacheEntry clearsentry from device cache
func ClearCacheEntry(ctx context.Context, name, uuid string) {
	if name != "" {
		logger.Debugw(ctx, "Clearing-name-key-from-device-cache", log.Fields{"name": name})
		cache.nameToRec.Delete(name)
	}
	if uuid != "" {
		logger.Debugw(ctx, "Clearing-uuid-key-from-device-cache", log.Fields{"uuid": uuid})
		cache.uuidToName.Delete(name)
	}
}

// DeviceRecord refers to the structure defined for storing OLT info
type DeviceRecord v1.DeviceRecordV1_0

// NewDeviceRecord return record for aliased ModifiableComponent
func NewDeviceRecord(ctx context.Context, req *dmi.ModifiableComponent) (*DeviceRecord, error) {
	rec := new(DeviceRecord)
	err := copier.Copy(&rec, &req)
	if err != nil {
		logger.Errorw(ctx, "Failed-at-creating-object-for-new-device-info", log.Fields{"error": err, "req": req})
		return nil, err
	}
	rec.Uri = req.Uri.Uri
	rec.State = new(dmi.ComponentState)
	rec.State.AdminState = req.AdminState
	logger.Infow(ctx, "Successful-at-creating-object-for-new-device-info", log.Fields{"new-object": rec})
	return rec, nil
}
