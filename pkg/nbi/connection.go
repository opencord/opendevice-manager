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

// Package nbi holds rpc server apis implemented
package nbi

import (
	"context"
	"sync"

	dev "github.com/opencord/opendevice-manager/pkg/models/device"
	"github.com/opencord/opendevice-manager/pkg/sbi"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// connectMap store all device connections established
type connectMap struct {
	nameToAdapter map[string]sbi.Adapter // key is name and value is adapter
	mutex         *sync.RWMutex          // mutex is used to lock when accessing
}

var connections *connectMap

// initConnectMap initialises map for storing connections
func initConnectMap() {
	connections = new(connectMap)
	connections.nameToAdapter = make(map[string]sbi.Adapter)
	connections.mutex = &sync.RWMutex{}
}

// DeInitConnectMap clears all stored connections
func DeInitConnectMap(ctx context.Context) {
	connections.nameToAdapter = nil
	connections.mutex = nil
	connections = nil
}

// getConnection retrieves connection object from map using name
func (conn *connectMap) getConnection(ctx context.Context, devRec *dev.DeviceRecord) (sbi.Adapter, error) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if val, ok := conn.nameToAdapter[devRec.Name]; ok {
		return val, nil
	}

	// Get the right adapter
	adapter := sbi.GetHwMgmtSvcClient(devRec)
	if err := adapter.Connect(ctx); err != nil {
		return nil, err
	}

	conn.nameToAdapter[devRec.Name] = adapter
	logger.Infow(ctx, "getConnection-completed", log.Fields{"name": devRec.Name, "adapter": adapter})

	return adapter, nil
}

// // storeConn stores connection object in map using uuid and name
// func (conn *connectMap) storeConnWithName(ctx context.Context, name string, adapter sbi.Adapter) {
// 	conn.mutex.Lock()
// 	defer conn.mutex.Unlock()
// 	conn.nameToAdapter[name] = adapter
// 	logger.Infow(ctx, "storeConnWithName-completed", log.Fields{"name": name, "adapter": adapter})
// }

// delConn deletes connection object from map using uuid and name
func (conn *connectMap) delConn(ctx context.Context, name string) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	if name != "" {
		delete(conn.nameToAdapter, name)
	}
	logger.Infow(ctx, "delConn-completed", log.Fields{"name": name})
}
