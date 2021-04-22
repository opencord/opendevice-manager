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
	"context"
	"encoding/json"
	"errors"
	"fmt"

	copy "github.com/jinzhu/copier"
	dmi "github.com/opencord/device-management-interface/go/dmi"

	"github.com/opencord/opendevice-manager/pkg/db"
	log "github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// DBAddByName inserts Device Info record to DB with Name as Key
func (rec *HwCompRecord) DBAddByUuid(ctx context.Context, deviceUuid string) error {
	if rec.Uuid == "" || deviceUuid == "" {
		logger.Errorw(ctx, "DBAddByUuid-failed-missing-uuid", log.Fields{"rec": rec, "dev-uuid": deviceUuid})
		return errors.New("missing uuid")
	}
	key := fmt.Sprintf(DbPathUuidToRecord, deviceUuid, rec.Uuid)
	b, _ := json.Marshal(rec)
	entry := string(b)
	err := db.Put(ctx, key, entry)
	cache.store(deviceUuid, rec)
	logger.Infow(ctx, "Inserting-hw-comp-info-to-Db-in-DBAddByUuid-method", log.Fields{"rec": rec, "error": err})
	return err
}

// DBAddByName inserts Device Info record to DB with Name as Key
func DBAddNameToUuidlookup(ctx context.Context, deviceUuid string, nameToUuidMap map[string]string) error {
	if deviceUuid == "" || len(nameToUuidMap) == 0 {
		logger.Errorw(ctx, "DBAddNameToUuidlookup-failed-missing-uuid-or-map", log.Fields{"map": nameToUuidMap, "dev-uuid": deviceUuid})
		return errors.New("missing info")
	}
	key := fmt.Sprintf(DbPathNameToUuid, deviceUuid)
	b, _ := json.Marshal(nameToUuidMap)
	entry := string(b)
	err := db.Put(ctx, key, entry)
	logger.Infow(ctx, "DBAddNameToUuidlookup-method-complete", log.Fields{"map": nameToUuidMap, "error": err})
	return err
}

// DBSaveHwCompsFromPhysicalInventory iterates through each children and store hwcomponents in db
func DBSaveHwCompsFromPhysicalInventory(ctx context.Context, deviceUuid string, nameToUuidMap map[string]string, children []*dmi.Component) {
	if len(children) == 0 {
		return
	}
	for _, child := range children {
		hwRec := new(HwCompRecord)
		if err := copy.Copy(&hwRec, &child); hwRec.Name == "" {
			logger.Errorw(ctx, "Failed-at-copying-hw-comp-from-inventory-list", log.Fields{"error": err, "child": child, "hw-comp": hwRec})
			continue
		}
		if child.Uri != nil {
			hwRec.Uri = child.Uri.Uri
		}
		if child.Uuid != nil {
			hwRec.Uuid = child.Uuid.Uuid
		}
		for _, grandChild := range child.Children {
			hwRec.Children = append(hwRec.Children, grandChild.Uuid.Uuid)
		}
		hwRec.DBAddByUuid(ctx, deviceUuid)
		nameToUuidMap[hwRec.Name] = hwRec.Uuid
		DBSaveHwCompsFromPhysicalInventory(ctx, deviceUuid, nameToUuidMap, child.Children)
		logger.Infow(ctx, "Successful-at-creating-object-for-new-hw-info", log.Fields{"new-object": child})
	}
}

// DBDelRecord deletes all entries for Device Info
func DBDelAllHwComponents(ctx context.Context, deviceUuid string) error {

	var err error

	if deviceUuid == "" {
		logger.Errorw(ctx, "deleting-all-hw-components-failed", log.Fields{"reason": "missing-device-uuid"})
		return errors.New("missing device uuid")
	}

	key := fmt.Sprintf(DbPrefix, deviceUuid)
	err = db.DelAll(ctx, key)
	cache.delDevice(deviceUuid)
	logger.Infow(ctx, "deleting-all-hw-components-completed", log.Fields{"key": key})

	return err
}

// DBGetRecByUuid func reads hw comp record by uuid
func DBGetRecByUuid(ctx context.Context, deviceUuid, hwCompUuid string) (*HwCompRecord, error) {
	if deviceUuid == "" || hwCompUuid == "" {
		logger.Errorw(ctx, "DBGetHwCompRec-failed-missing-uuid", log.Fields{"device-uuid": deviceUuid, "hw-comp-uuid": hwCompUuid})
		return nil, errors.New("uuid field is empty")
	}

	if rec := cache.get(deviceUuid, hwCompUuid); rec != nil {
		return rec, nil
	}

	key := fmt.Sprintf(DbPathUuidToRecord, deviceUuid, hwCompUuid)
	entry, err := db.Read(ctx, key)
	if err != nil {
		logger.Errorw(ctx, "DBGetRecByUuid-failed-read-db", log.Fields{"error": err, "key": key})
		return nil, err
	}

	rec := new(HwCompRecord)
	if err = json.Unmarshal([]byte(entry), rec); err != nil {
		logger.Errorw(ctx, "Failed-to-unmarshal-at-DBGetRecByUuid", log.Fields{"reason": err, "entry": entry})
		return nil, err
	}

	cache.store(deviceUuid, rec)

	logger.Debugw(ctx, "DBGetHwCompRec-completed", log.Fields{"device-uuid": deviceUuid, "hw-comp-uuid": hwCompUuid, "rec": rec})
	return rec, nil
}

// DBGetRecByName func reads hw comp record by name
func DBGetRecByName(ctx context.Context, deviceUuid, hwName string) (*HwCompRecord, error) {
	if deviceUuid == "" || hwName == "" {
		logger.Errorw(ctx, "DBGetRecByName-failed-missing-uuid", log.Fields{"device-uuid": deviceUuid, "hw-comp-name": hwName})
		return nil, errors.New("name field is empty")
	}
	key := fmt.Sprintf(DbPathNameToUuid, deviceUuid)
	entry, err := db.Read(ctx, key)
	if err != nil {
		logger.Errorw(ctx, "DBGetRecByName-failed-read-db", log.Fields{"error": err, "key": key})
		return nil, err
	}

	nameToUuidMap := make(map[string]string)
	if err = json.Unmarshal([]byte(entry), &nameToUuidMap); nameToUuidMap == nil || err != nil {
		logger.Errorw(ctx, "Failed-to-unmarshal-at-DBGetRecByName", log.Fields{"reason": err, "entry": entry})
		return nil, err
	}

	if hwUuid, ok := nameToUuidMap[hwName]; ok {
		rec, err2 := DBGetRecByUuid(ctx, deviceUuid, hwUuid)
		logger.Debugw(ctx, "DBGetRecByName-completed", log.Fields{"device-uuid": deviceUuid, "hw-comp-name": hwName, "rec": rec, "error": err2})
	}

	return nil, errors.New("name not found")
}
