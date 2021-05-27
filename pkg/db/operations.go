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

// Package db holds utils for datastore implementation
package db

import (
	"context"
	"errors"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// Read function reads key value pair from db/kvstore
func Read(ctx context.Context, key string) (string, error) {
	if kvClient != nil {
		logger.Debugw(ctx, "Reading-key-value-pair-from-kv-store", log.Fields{"key": key})
		kvPair, err := kvClient.Get(ctx, key)
		if err != nil {
			return "", err
		}
		if kvPair == nil {
			return "", errors.New("key not found")
		}
		return string(kvPair.Value.([]byte)), nil

	}
	logger.Errorw(ctx, "Reading-key-value-pair-in kv-store-failed-because-kvstore-not-initialised", log.Fields{"key": key})
	return "", errors.New("kvstore not initialised")
}

// ReadAll function reads all key value pair from db/kvstore
func ReadAll(ctx context.Context, keyPrefix string) (map[string]string, error) {
	keyValues := make(map[string]string)
	if kvClient != nil {
		logger.Debugw(ctx, "Reading-all-key-value-pairs-from-kv-store", log.Fields{"key-prefix": keyPrefix})
		kvPairs, err := kvClient.List(ctx, keyPrefix)
		if err != nil {
			return keyValues, err
		}
		if kvPairs == nil {
			return keyValues, errors.New("key not found")
		}

		for key, kvPair := range kvPairs {
			keyValues[key] = string(kvPair.Value.([]byte))
		}
		return keyValues, nil
	}
	logger.Errorw(ctx, "Reading-all-key-value-pair-in-kv-store-failed-because-kvstore-not-initialised", log.Fields{"key-prefix": keyPrefix})
	return keyValues, errors.New("kvstore not initialised")
}

// Del function deletes key value pair from db/kvstore
func Del(ctx context.Context, key string) error {
	if kvClient != nil {
		logger.Debugw(ctx, "Deleting-key-value-pair-from-kv-store", log.Fields{"key": key})
		return kvClient.Delete(ctx, key)
	}
	logger.Errorw(ctx, "Deleting-key-value-pair-in-kv-store-failed-because-kvstore-not-initialised", log.Fields{"key": key})
	return errors.New("kvstore not initialised")
}

// DelAll function deletes all key value pair from db/kvstore with provided key prefix
func DelAll(ctx context.Context, keyPrefix string) error {
	if kvClient != nil {
		logger.Debugw(ctx, "Deleting-all-key-value-pair-from-kv-store-with-prefix", log.Fields{"key-prefix": keyPrefix})
		return kvClient.DeleteWithPrefix(ctx, keyPrefix)
	}
	logger.Errorw(ctx, "Deleting-all-key-value-pair-in-kv-store-with-prefix-failed-because-kvstore-not-initialised", log.Fields{"key-prefix": keyPrefix})
	return errors.New("kvstore not initialised")
}

// Put function stores key value pair in db/kvstore
func Put(ctx context.Context, key string, val string) error {
	if kvClient != nil {
		logger.Debugw(ctx, "Storing-key-value-pair-in-kv-store", log.Fields{"key": key, "value": val})
		return kvClient.Put(ctx, key, val)
	}
	logger.Errorw(ctx, "Storing-key-value-pair-in-kv-store-failed-because-kvstore-not-initialised", log.Fields{"key": key, "value": val})
	return errors.New("kvstore not initialised")
}
