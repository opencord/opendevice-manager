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

// Package sbi holds interfaces for adapter operations
package sbi

import (
	"context"

	dmi "github.com/opencord/device-management-interface/go/dmi"
	dev "github.com/opencord/opendevice-manager/pkg/models/device"
	hw "github.com/opencord/opendevice-manager/pkg/models/hwcomponents"
	grpc "github.com/opencord/opendevice-manager/pkg/sbi/grpc"
)

// GetHwMgmtSvcClient returns the adapter
func GetHwMgmtSvcClient(devRec *dev.DeviceRecord) Adapter {
	switch devRec.Make {
	case "ROLT":
		return grpc.NewClient(devRec.Uri)
	}
	return grpc.NewClient(devRec.Uri)
}

// Adapter interface contains all methods for rpc calls
type Adapter interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	AdapterHwMgmtSvc
}

// AdapterHwMgmtSvc refers to the interface used for defining RPCs for HW management Service
type AdapterHwMgmtSvc interface {
	StartManagingDevice(context.Context, *dev.DeviceRecord, *dmi.ModifiableComponent, dmi.NativeHWManagementService_StartManagingDeviceServer) (error, bool)
	StopManagingDevice(context.Context, *dev.DeviceRecord, *dmi.StopManagingDeviceRequest) (*dmi.StopManagingDeviceResponse, error)
	GetPhysicalInventory(context.Context, *dev.DeviceRecord, *dmi.PhysicalInventoryRequest, dmi.NativeHWManagementService_GetPhysicalInventoryServer) error
	SetHWComponentInfo(context.Context, string, *hw.HwCompRecord, *dmi.HWComponentInfoSetRequest) (*dmi.HWComponentInfoSetResponse, error)
	GetHWComponentInfo(context.Context, string, *hw.HwCompRecord, *dmi.HWComponentInfoGetRequest, dmi.NativeHWManagementService_GetHWComponentInfoServer) error
	SetLoggingEndpoint(context.Context, *dev.DeviceRecord, *dmi.SetLoggingEndpointRequest) (*dmi.SetRemoteEndpointResponse, error)
	GetLoggingEndpoint(context.Context, *dev.DeviceRecord, *dmi.HardwareID) (*dmi.GetLoggingEndpointResponse, error)
	GetLoggableEntities(context.Context, *dev.DeviceRecord, *dmi.GetLoggableEntitiesRequest) (*dmi.GetLogLevelResponse, error)
	SetLogLevel(context.Context, *dev.DeviceRecord, *dmi.SetLogLevelRequest) (*dmi.SetLogLevelResponse, error)
	GetLogLevel(context.Context, *dev.DeviceRecord, *dmi.GetLogLevelRequest) (*dmi.GetLogLevelResponse, error)
}
