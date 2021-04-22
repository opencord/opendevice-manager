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

// Package hwcomponents stores methods and functions related to hardware
package hwcomponents

import (
	"sync"

	config "github.com/opencord/opendevice-manager/pkg/config"
	v1 "github.com/opencord/opendevice-manager/pkg/models/hwcomponents/v1"
	log "github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// Constants defined are the DB Path meant for storing hw component info records
const (
	DbPrefix = config.DBPrefix + config.CurDBVer + "/HwCompRec/%s"
	// Key : /OpenDevMgr/v1/HwCompRec/{Device-Uuid}/Components
	// Val : Map => {"hw-comp-name-1":"hw-comp-uuid-1", "hw-comp-name-2":"hw-comp-uuid-2"}
	DbPathNameToUuid = DbPrefix + "/Components"
	// Key : /OpenDevMgr/v1/HwCompRec/{Device-Uuid}/Uuid/{Hw-Comp-Uuid}
	// Val : HwCompRecord{}
	DbPathUuidToRecord = DbPrefix + "/Uuid/%s"
)

// compCache stores component information in buffer
type compCache struct {
	uuidToRec map[string]map[string]*HwCompRecord // nameToRecord maintains cache for mapping from name to main record
	mutex     sync.Mutex
}

var cache *compCache

// logger represents the log object
var logger log.CLogger

// initCache initialises device cache
func initCache() {
	cache = new(compCache)
	cache.uuidToRec = make(map[string]map[string]*HwCompRecord)
	cache.mutex = sync.Mutex{}
}

// init function for the package
func init() {
	logger = config.Initlog()
	initCache()
}

type HwCompRecord v1.HwCompRecordV1_0

func (*compCache) store(devUuid string, rec *HwCompRecord) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var uuidToRecMap map[string]*HwCompRecord

	if val, ok := cache.uuidToRec[devUuid]; !ok {
		uuidToRecMap = make(map[string]*HwCompRecord)
	} else {
		uuidToRecMap = val
	}

	uuidToRecMap[rec.Uuid] = rec
	cache.uuidToRec[devUuid] = uuidToRecMap
}

func (*compCache) get(devUuid, compUuid string) *HwCompRecord {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var uuidToRecMap map[string]*HwCompRecord

	if val, ok := cache.uuidToRec[devUuid]; !ok {
		return nil
	} else {
		uuidToRecMap = val
	}

	if rec, ok := uuidToRecMap[compUuid]; ok {
		return rec
	}

	return nil
}

func (*compCache) delDevice(devUuid string) *HwCompRecord {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.uuidToRec, devUuid)
	return nil
}
