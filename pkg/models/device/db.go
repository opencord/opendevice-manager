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

// Package modifiablecomponent stores ModifiableComponent methods and functions
package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opencord/device-management-interface/go/dmi"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"

	"github.com/opencord/opendevice-manager/pkg/db"

	copy "github.com/jinzhu/copier"
)

// DBGetByName func reads device record by name
func DBGetByName(ctx context.Context, name string) (*DeviceRecord, error) {
	if name == "" {
		logger.Errorw(ctx, "DBGetByName-failed-missing-device-name", log.Fields{})
		return nil, errors.New("name field is empty")
	}

	logger.Debugw(ctx, "DBGetByName-invoked", log.Fields{"name": name})
	defer logger.Debugw(ctx, "DBGetByName-completed", log.Fields{"name": name})

	if val, ok := cache.nameToRec.Load(name); ok {
		return val.(*DeviceRecord), nil
	}

	key := fmt.Sprintf(DbPathNameToRecord, name)
	entry, err := db.Read(ctx, key)
	if err != nil {
		logger.Errorw(ctx, "DBGetByName-failed-read-db", log.Fields{"error": err, "key": key})
		return nil, err
	}

	rec := new(DeviceRecord)
	if err = json.Unmarshal([]byte(entry), rec); err != nil {
		logger.Errorw(ctx, "Failed-to-unmarshal-at-DBGetByName", log.Fields{"reason": err, "entry": entry})
		return nil, err
	}

	cache.nameToRec.Store(name, rec)

	return rec, nil
}

// DBGetByUuid func reads device record by Uuid
func DBGetByUuid(ctx context.Context, uuid string) (*DeviceRecord, error) {

	if uuid == "" {
		logger.Errorw(ctx, "DBGetByUuid-failed-missing-device-uuid", log.Fields{})
		return nil, errors.New("uuid field is empty")
	}

	logger.Debugw(ctx, "DBGetByUuid-invoked", log.Fields{"uuid": uuid})
	defer logger.Debugw(ctx, "DBGetByUuid-completed", log.Fields{"uuid": uuid})

	var name string
	var err error

	if val, ok := cache.uuidToName.Load(uuid); ok {
		name = val.(string)
	} else {

		key := fmt.Sprintf(DbPathUuidToName, uuid)
		name, err = db.Read(ctx, key)
		if err != nil {
			logger.Errorw(ctx, "DBGetByUuid-failed-read-db", log.Fields{"error": err, "key": key})
			return nil, err
		}
	}

	cache.uuidToName.Store(uuid, name)

	return DBGetByName(ctx, name)
}

// DBGetAll func returns all device records
func DBGetAll(ctx context.Context) ([]*DeviceRecord, error) {
	key := fmt.Sprintf(DbPathNameToRecord, "")
	kvPairs, err := db.ReadAll(ctx, key)
	if err != nil {
		logger.Errorw(ctx, "DBGetAll-failed-read-db", log.Fields{"error": err, "key": key})
		return nil, err
	}

	var listDevs []*DeviceRecord

	for _, entry := range kvPairs {
		rec := new(DeviceRecord)
		if err = json.Unmarshal([]byte(entry), rec); err != nil {
			logger.Errorw(ctx, "Failed-to-unmarshal-at-DBGetByName", log.Fields{"reason": err, "entry": entry})
		} else {
			listDevs = append(listDevs, rec)
		}
	}

	logger.Infow(ctx, "DBGetAll-success", log.Fields{"entry": listDevs})

	return listDevs, nil
}

// DBAddByName inserts Device Info record to DB with Name as Key
func (rec *DeviceRecord) DBAddByName(ctx context.Context) error {
	if rec.Name == "" {
		logger.Errorw(ctx, "DBAddByName-failed-missing-device-name", log.Fields{"rec": rec})
		return errors.New("missing name")
	}
	key := fmt.Sprintf(DbPathNameToRecord, rec.Name)
	b, _ := json.Marshal(rec)
	entry := string(b)
	err := db.Put(ctx, key, entry)
	cache.nameToRec.Store(rec.Name, rec)
	logger.Infow(ctx, "Inserting-device-info-to-Db-in-DBAddByName-method", log.Fields{"rec": rec, "error": err})
	return err
}

// DBAddUuidLookup creates a lookup of name from uuid
func (rec *DeviceRecord) DBAddUuidLookup(ctx context.Context) error {
	if rec.Uuid == "" || rec.Name == "" {
		logger.Errorw(ctx, "DBAddUuidLookup-failed-missing-device-name-or-uuid", log.Fields{"rec": rec})
		return errors.New("missing name")
	}
	key := fmt.Sprintf(DbPathUuidToName, rec.Uuid)
	err := db.Put(ctx, key, rec.Name)
	cache.uuidToName.Store(rec.Uuid, rec.Name)
	logger.Infow(ctx, "DBAddUuidLookup-success", log.Fields{"rec": rec, "error": err})
	return err
}

// DBDelRecord deletes all entries for Device Info
func (rec *DeviceRecord) DBDelRecord(ctx context.Context) error {

	var err error

	if rec.Name != "" {
		key := fmt.Sprintf(DbPathNameToRecord, rec.Name)
		logger.Infow(ctx, "deleting-device-info-record-using-name", log.Fields{"name": rec.Name, "key": key})
		err = db.Del(ctx, key)
		cache.nameToRec.Delete(rec.Name)
	}

	if err == nil && rec.Uuid != "" {
		key := fmt.Sprintf(DbPathUuidToName, rec.Uuid)
		logger.Infow(ctx, "deleting-device-info-record-using-uuid", log.Fields{"uuid": rec.Uuid, "key": key})
		err = db.Del(ctx, key)
		cache.uuidToName.Delete(rec.Uuid)
	}

	return err
}

// DBSaveHwInfo stores hardware copies info from response and stores in db
func (rec *DeviceRecord) DBSaveHwInfo(ctx context.Context, hw *dmi.Hardware) error {
	defer logger.Infow(ctx, "saving-hw-info-to-device-record-completed", log.Fields{"rec": rec})
	rec.LastBooted = hw.LastBooted
	rec.LastChange = hw.LastChange
	name := rec.Name
	uuid := rec.Uuid
	if err := copy.Copy(&rec, &hw.Root); err != nil {
		logger.Errorw(ctx, "copy-failed-at-DBSaveHwInfo", log.Fields{"rec": rec, "error": err, "hw": hw})
		return err
	}
	rec.Children = []string{}
	for _, child := range hw.Root.Children {
		rec.Children = append(rec.Children, child.Uuid.Uuid)
	}
	rec.Name = name
	rec.Uuid = uuid
	return rec.DBAddByName(ctx)
}
