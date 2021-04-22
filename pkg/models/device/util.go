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

	"github.com/opencord/device-management-interface/go/dmi"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// GetLoggableEntitiesFromDevRec represnets the fetch the log level with entity from device record
func (rec *DeviceRecord) GetLoggableEntitiesFromDevRec(ctx context.Context, entities []string) ([]*dmi.EntitiesLogLevel, bool) {

	var traceLevel, debugLevel, infoLevel, warnLevel, errorLevel []string

	if len(entities) > 0 {
		// getLogLevel request for given entities
		for _, entity := range entities {
			if logLevel, ok := rec.Logging.LoggableEntities[entity]; ok {
				switch logLevel {
				case dmi.LogLevel_TRACE:
					traceLevel = append(traceLevel, entity)
				case dmi.LogLevel_DEBUG:
					debugLevel = append(debugLevel, entity)
				case dmi.LogLevel_INFO:
					infoLevel = append(infoLevel, entity)
				case dmi.LogLevel_WARN:
					warnLevel = append(warnLevel, entity)
				case dmi.LogLevel_ERROR:
					errorLevel = append(errorLevel, entity)
				}
			} else {
				logger.Warnw(ctx, "entity-was-not-found-in-device-record", log.Fields{"device-name": rec.Name, "entity": entity})
				return nil, false
			}

		}
	} else if len(rec.Logging.LoggableEntities) == 0 {
		// if LoggableEntities length is zero means loglevel is applicable for entire hardware
		logger.Debug(ctx, "all-entities-have-common-loglevel", log.Fields{"device-name": rec.Name, "log-level": rec.Logging.LogLevel})
		return []*dmi.EntitiesLogLevel{{LogLevel: rec.Logging.LogLevel}}, true
	} else {
		// get globle log level or get loggble entities will  invoke here
		logger.Debug(ctx, "all-entities-have-diffrent-loglevel", log.Fields{"device-name": rec.Name})
		for entity, logLevel := range rec.Logging.LoggableEntities {
			switch logLevel {
			case dmi.LogLevel_TRACE:
				traceLevel = append(traceLevel, entity)
			case dmi.LogLevel_DEBUG:
				debugLevel = append(debugLevel, entity)
			case dmi.LogLevel_INFO:
				infoLevel = append(infoLevel, entity)
			case dmi.LogLevel_WARN:
				warnLevel = append(warnLevel, entity)
			case dmi.LogLevel_ERROR:
				errorLevel = append(errorLevel, entity)
			}
		}
	}

	entitiesLogLevel := []*dmi.EntitiesLogLevel{
		{LogLevel: dmi.LogLevel_TRACE, Entities: traceLevel},
		{LogLevel: dmi.LogLevel_DEBUG, Entities: debugLevel},
		{LogLevel: dmi.LogLevel_INFO, Entities: infoLevel},
		{LogLevel: dmi.LogLevel_WARN, Entities: warnLevel},
		{LogLevel: dmi.LogLevel_ERROR, Entities: errorLevel},
	}
	logger.Debug(ctx, "entities-with-log-level", log.Fields{"entities": entitiesLogLevel})

	return entitiesLogLevel, true
}

// SaveLoggableEntities func is is used to save the log level with entity in device record
func (rec *DeviceRecord) SaveLoggableEntities(ctx context.Context, listEntities []*dmi.EntitiesLogLevel) {

	if rec.Logging.LoggableEntities == nil {

		logger.Debug(ctx, "allocating-memory-for-loggable-entitie", log.Fields{"device-name": rec.Name})
		rec.Logging.LoggableEntities = make(map[string]dmi.LogLevel)
	}

	if len(listEntities) == 1 && listEntities[0].Entities == nil {

		logger.Debug(ctx, "set-global-log-level", log.Fields{"device-name": rec.Name})
		rec.Logging.LogLevel = listEntities[0].LogLevel

	} else {

		logger.Debug(ctx, "setting-entity-log-level", log.Fields{"device-name": rec.Name, "list-of-entities": listEntities})
		for _, entities := range listEntities {
			logLevel := entities.LogLevel
			for _, entity := range entities.Entities {
				rec.Logging.LoggableEntities[entity] = logLevel
			}
		}
	}
}
