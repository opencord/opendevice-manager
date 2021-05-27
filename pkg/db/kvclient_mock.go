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
	"strings"
	"sync"
	"time"

	"github.com/opencord/voltha-lib-go/v4/pkg/db/kvstore"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

type mockKVClient struct {
}

var testKvPairCache *sync.Map

// MockKVClient function mimics the kvclient
func MockKVClient() {
	kvClient = new(mockKVClient)
	testKvPairCache = new(sync.Map)
}

// ClearCache function clears the kvclient cache
func ClearCache() {
	testKvPairCache = new(sync.Map)
}

// List function implemented for KVClient.
func (kvclient *mockKVClient) List(ctx context.Context, prefix string) (map[string]*kvstore.KVPair, error) {
	kvPairMap := make(map[string]*kvstore.KVPair)
	testKvPairCache.Range(func(key, value interface{}) bool {
		if strings.Contains(key.(string), prefix) {
			kvPair := new(kvstore.KVPair)
			kvPair.Key = key.(string)
			kvPair.Value = value.([]byte)
			kvPairMap[kvPair.Key] = kvPair
		}
		return true
	})

	if len(kvPairMap) != 0 {
		logger.Debugw(ctx, "List of MockKVClient called", log.Fields{"kvPairMap": kvPairMap})
		return kvPairMap, nil
	}

	return nil, errors.New("key didn't find")
}

// Get mock function implementation for KVClient
func (kvclient *mockKVClient) Get(ctx context.Context, key string) (*kvstore.KVPair, error) {
	logger.Debugw(ctx, "Warning Warning Warning: Get of MockKVClient called", log.Fields{"key": key})

	if val, ok := testKvPairCache.Load(key); ok {
		kvPair := new(kvstore.KVPair)
		kvPair.Key = key
		kvPair.Value = val
		return kvPair, nil
	}

	return nil, errors.New("key didn't find")
}

// Put mock function implementation for KVClient
func (kvclient *mockKVClient) Put(ctx context.Context, key string, value interface{}) error {
	if key != "" {
		value = []byte(value.(string))
		testKvPairCache.Store(key, value)
		return nil
	}
	return errors.New("key didn't find")
}

// Delete mock function implementation for KVClient
func (kvclient *mockKVClient) Delete(ctx context.Context, key string) error {
	logger.Infow(ctx, "Error Error Error Key:", log.Fields{})
	if key != "" {
		testKvPairCache.Delete(key)
		return nil
	}
	return errors.New("key didn't find")
}

// DeleteWithPrefix mock function implementation for KVClient
func (kvclient *mockKVClient) DeleteWithPrefix(ctx context.Context, prefix string) error {
	testKvPairCache.Range(func(key, value interface{}) bool {
		if strings.Contains(key.(string), prefix) {
			testKvPairCache.Delete(key)
		}
		return true
	})
	return nil
}

// Reserve mock function implementation for KVClient
func (kvclient *mockKVClient) Reserve(ctx context.Context, key string, value interface{}, ttl time.Duration) (interface{}, error) {
	return nil, errors.New("key didn't find")
}

// ReleaseReservation mock function implementation for KVClient
func (kvclient *mockKVClient) ReleaseReservation(ctx context.Context, key string) error {
	return nil
}

// ReleaseAllReservations mock function implementation for KVClient
func (kvclient *mockKVClient) ReleaseAllReservations(ctx context.Context) error {
	return nil
}

// RenewReservation mock function implementation for KVClient
func (kvclient *mockKVClient) RenewReservation(ctx context.Context, key string) error {
	return nil
}

// Watch mock function implementation for KVClient
func (kvclient *mockKVClient) Watch(ctx context.Context, key string, withPrefix bool) chan *kvstore.Event {
	return nil
}

// AcquireLock mock function implementation for KVClient
func (kvclient *mockKVClient) AcquireLock(ctx context.Context, lockName string, timeout time.Duration) error {
	return nil
}

// ReleaseLock mock function implementation for KVClient
func (kvclient *mockKVClient) ReleaseLock(lockName string) error {
	return nil
}

// IsConnectionUp mock function implementation for KVClient
func (kvclient *mockKVClient) IsConnectionUp(ctx context.Context) bool { // timeout in second
	return true
}

// CloseWatch mock function implementation for KVClient
func (kvclient *mockKVClient) CloseWatch(ctx context.Context, key string, ch chan *kvstore.Event) {
}

// Close mock function implementation for KVClient
func (kvclient *mockKVClient) Close(ctx context.Context) {
}
