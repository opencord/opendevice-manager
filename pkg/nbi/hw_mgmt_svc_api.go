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
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/opencord/device-management-interface/go/dmi"
	"github.com/opencord/opendevice-manager/pkg/config"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"

	dev "github.com/opencord/opendevice-manager/pkg/models/device"
	hw "github.com/opencord/opendevice-manager/pkg/models/hwcomponents"

	empty "github.com/golang/protobuf/ptypes/empty"
)

// StartManagingDevice refers to the RPC method invoked for start Managing Device.
// Initializes context for a device and sets up required states
func (c *NativeHwManagementService) StartManagingDevice(req *dmi.ModifiableComponent, streamResp dmi.NativeHWManagementService_StartManagingDeviceServer) error {
	ctx := config.GetNewContextFromGlobalContxt("StartManagingDevice")

	logger.Infow(ctx, "StartManagingDevice-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "StartManagingDevice-on-grpc-server-completed", log.Fields{"req": req})

	if errorResp, ok := validateStartManagingDeviceReq(ctx, req); !ok {
		streamResp.Send(errorResp)
		return errors.New("failed at validateStartManagingDeviceReq")
	}

	devRec, err := dev.NewDeviceRecord(ctx, req)

	if err != nil {
		streamResp.Send(errRespStartManagingDevice(ctx, req, dmi.StartManagingDeviceResponse_INVALID_PARAMS, err))
		return err
	}

	adapter, err := connections.getConnection(ctx, devRec)

	if adapter == nil {
		streamResp.Send(errRespStartManagingDevice(ctx, req, dmi.StartManagingDeviceResponse_INVALID_PARAMS, err))
		return err
	}

	devRec.DBAddByName(ctx)

	err, connMade := adapter.StartManagingDevice(ctx, devRec, req, streamResp)

	if !connMade {
		devRec.DBDelRecord(ctx)
		adapter.Disconnect(ctx)
		connections.delConn(ctx, devRec.Name)
		streamResp.Send(errRespStartManagingDevice(ctx, req, dmi.StartManagingDeviceResponse_UNDEFINED_REASON, err))
	}

	return err
}

// StopManagingDevice - Stops management of a device and clean up any context and caches for that device
// This rpc can be called at any time, even before the StartManagingDevice operation
// has completed, and should be able to cleanup.
func (c *NativeHwManagementService) StopManagingDevice(ctx context.Context, req *dmi.StopManagingDeviceRequest) (*dmi.StopManagingDeviceResponse, error) {

	var devRec *dev.DeviceRecord
	resp := new(dmi.StopManagingDeviceResponse)
	var err error

	ctx = config.GetNewContextFromContxt(ctx, "StopManagingDevice")

	logger.Infow(ctx, "StopManagingDevice-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "StopManagingDevice-on-grpc-server-completed", log.Fields{"req": req})

	if recTmp, err := dev.DBGetByName(ctx, req.Name); recTmp == nil {
		return errRespStopManagingDevice(ctx, req, dmi.StopManagingDeviceResponse_UNKNOWN_DEVICE, err), err
	} else {
		devRec = recTmp
	}

	defer devRec.DBDelRecord(ctx)
	defer hw.DBDelAllHwComponents(ctx, devRec.Uuid)

	adapter, err := connections.getConnection(ctx, devRec)

	if adapter != nil {
		resp, err = adapter.StopManagingDevice(ctx, devRec, req)
		adapter.Disconnect(ctx)
		connections.delConn(ctx, devRec.Name)
	}

	if err != nil {
		logger.Errorw(ctx, "Errors-at-StopManagingDevice", log.Fields{"req": req, "error": err})
	}

	return resp, err
}

// GetManagedDevices - Returns an object containing a list of devices managed by this entity
func (c *NativeHwManagementService) GetManagedDevices(ctx context.Context, req *empty.Empty) (*dmi.ManagedDevicesResponse, error) {
	ctx = config.GetNewContextFromContxt(ctx, "GetManagedDevices")
	resp := new(dmi.ManagedDevicesResponse)

	logger.Infow(ctx, "GetManagedDevices-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "GetManagedDevices-on-grpc-server-completed", log.Fields{"req": req})

	listDevRecs, err := dev.DBGetAll(ctx)
	for _, devRec := range listDevRecs {
		modComp := new(dmi.ModifiableComponent)
		err2 := copier.Copy(&modComp, &devRec)
		if err2 != nil {
			logger.Errorw(ctx, "Copy-failed-at-GetManagedDevices", log.Fields{"req": req, "rec": devRec, "error": err})
		} else {
			modComp.Uri.Uri = devRec.Uri
			resp.Devices = append(resp.Devices, modComp)
		}
	}
	logger.Infow(ctx, "GetManagedDevices-completed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp, err
}

// GetPhysicalInventory - Place Holder for implementation in future
func (c *NativeHwManagementService) GetPhysicalInventory(req *dmi.PhysicalInventoryRequest, streamResp dmi.NativeHWManagementService_GetPhysicalInventoryServer) error {

	ctx := config.GetNewContextFromGlobalContxt("GetPhysicalInventory")

	logger.Infow(ctx, "GetPhysicalInventory-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "GetPhysicalInventory-on-grpc-server-completed", log.Fields{"req": req})

	var devRec *dev.DeviceRecord

	if recTmp, err := dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid); recTmp == nil {
		streamResp.Send(errRespGetPhysicatInventory(ctx, req, dmi.PhysicalInventoryResponse_UNKNOWN_DEVICE, err))
		return err
	} else {
		devRec = recTmp
	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		streamResp.Send(errRespGetPhysicatInventory(ctx, req, dmi.PhysicalInventoryResponse_DEVICE_UNREACHABLE, err))
		return err
	}

	err = adapter.GetPhysicalInventory(ctx, devRec, req, streamResp)

	if err != nil {
		logger.Errorw(ctx, "Errors-at-GetPhysicalInventory", log.Fields{"req": req, "error": err})
	}

	return err
}

// GetHWComponentInfo - refers to the RPC method invoked for get the details of a particular HW component
func (c *NativeHwManagementService) GetHWComponentInfo(req *dmi.HWComponentInfoGetRequest, streamResp dmi.NativeHWManagementService_GetHWComponentInfoServer) error {
	ctx := config.GetNewContextFromGlobalContxt("GetHWComponentInfo")

	logger.Infow(ctx, "GetHWComponentInfo-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "GetHWComponentInfo-on-grpc-server-completed", log.Fields{"req": req})

	devRec, err := dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid)
	if devRec == nil {
		streamResp.Send(errRespGetHWComponentInfo(ctx, req, dmi.HWComponentInfoGetResponse_UNKNOWN_DEVICE, err))
		return err
	}

	hwCompRec, err := hw.DBGetRecByUuid(ctx, req.DeviceUuid.Uuid, req.ComponentUuid.Uuid)
	if hwCompRec == nil {
		streamResp.Send(errRespGetHWComponentInfo(ctx, req, dmi.HWComponentInfoGetResponse_UNKNOWN_DEVICE, err))
		return err
	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		streamResp.Send(errRespGetHWComponentInfo(ctx, req, dmi.HWComponentInfoGetResponse_DEVICE_UNREACHABLE, err))
		return err
	}

	err = adapter.GetHWComponentInfo(ctx, devRec.Uuid, hwCompRec, req, streamResp)
	if err != nil {
		logger.Errorw(ctx, "Errors-at-GetHWComponentInfo", log.Fields{"req": req, "error": err})
	}
	return err
}

// SetHWComponentInfo is the nb api exposed for receiving SetHWComponentInfo from NEM
func (c *NativeHwManagementService) SetHWComponentInfo(ctx context.Context, req *dmi.HWComponentInfoSetRequest) (*dmi.HWComponentInfoSetResponse, error) {
	ctx = config.GetNewContextFromContxt(ctx, "SetHWComponentInfo")

	logger.Infow(ctx, "SetHWComponentInfo-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "SetHWComponentInfo-on-grpc-server-completed", log.Fields{"req": req})

	var devRec *dev.DeviceRecord
	var hwCompRec *hw.HwCompRecord
	var err error

	if devRec, err = dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid); devRec == nil {
		return errRespSetHWComponentInfo(ctx, req, dmi.HWComponentInfoSetResponse_UNKNOWN_DEVICE, err), err
	}

	if hwCompRec, err = hw.DBGetRecByUuid(ctx, req.DeviceUuid.Uuid, req.ComponentUuid.Uuid); hwCompRec == nil {
		return errRespSetHWComponentInfo(ctx, req, dmi.HWComponentInfoSetResponse_UNKNOWN_DEVICE, err), err
	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		return errRespSetHWComponentInfo(ctx, req, dmi.HWComponentInfoSetResponse_DEVICE_UNREACHABLE, err), err
	}

	resp, err := adapter.SetHWComponentInfo(ctx, devRec.Uuid, hwCompRec, req)

	if err != nil {
		logger.Errorw(ctx, "Errors-at-SetHWComponentInfo", log.Fields{"req": req, "error": err})
	}

	return resp, err
}

// SetLoggingEndpoint - refers to the RPC method invoked for set the location to which logs need to be shipped
func (c *NativeHwManagementService) SetLoggingEndpoint(ctx context.Context, req *dmi.SetLoggingEndpointRequest) (*dmi.SetRemoteEndpointResponse, error) {
	ctx = config.GetNewContextFromGlobalContxt("SetLoggingEndpoint")

	logger.Infow(ctx, "SetLoggingEndpoint-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "SetLoggingEndpoint-on-grpc-server-completed", log.Fields{"req": req})

	// Check device is there or not by Uuid
	devRec, err := dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid)
	if err != nil {
		return errRespSetLoggingEndpoint(ctx, req, dmi.SetRemoteEndpointResponse_UNKNOWN_DEVICE, err), err
	}

	adapter, err2 := connections.getConnection(ctx, devRec)
	if adapter == nil {
		return errRespSetLoggingEndpoint(ctx, req, dmi.SetRemoteEndpointResponse_DEVICE_UNREACHABLE, err2), err2
	}

	resp, err := adapter.SetLoggingEndpoint(ctx, devRec, req)
	if err != nil {
		logger.Errorw(ctx, "Errors-at-SetLoggingEndpoint", log.Fields{"req": req, "error": err})
	}

	return resp, err
}

// GetLoggingEndpoint - refers to the RPC method invoked for get the location to which logs need to be shipped
func (c *NativeHwManagementService) GetLoggingEndpoint(ctx context.Context, req *dmi.HardwareID) (*dmi.GetLoggingEndpointResponse, error) {
	ctx = config.GetNewContextFromGlobalContxt("GetLoggingEndpoint")

	logger.Infow(ctx, "GetLoggingEndpoint-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "GetLoggingEndpoint-on-grpc-server-completed", log.Fields{"req": req})

	// Check device is there or not by Uuid
	devRec, err := dev.DBGetByUuid(ctx, req.Uuid.Uuid)
	if err != nil {
		return errRespGetLoggingEndpoint(ctx, req, dmi.GetLoggingEndpointResponse_UNKNOWN_DEVICE, err), err
	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		return errRespGetLoggingEndpoint(ctx, req, dmi.GetLoggingEndpointResponse_DEVICE_UNREACHABLE, err), err
	}

	resp, err := adapter.GetLoggingEndpoint(ctx, devRec, req)
	if err != nil {
		logger.Errorw(ctx, "Errors-at-GetLoggingEndpoint", log.Fields{"req": req, "error": err})
	}
	return resp, err
}

// SetMsgBusEndpoint - Place Holder for implementation in future
func (c *NativeHwManagementService) SetMsgBusEndpoint(ctx context.Context, req *dmi.SetMsgBusEndpointRequest) (*dmi.SetRemoteEndpointResponse, error) {
	errMsg := "SetMsgBusEndpoint not yet implemented"
	fmt.Println(errMsg)
	return nil, errors.New(errMsg)
}

// GetMsgBusEndpoint - Place Holder for implementation in future
func (c *NativeHwManagementService) GetMsgBusEndpoint(ctx context.Context, req *empty.Empty) (*dmi.GetMsgBusEndpointResponse, error) {
	errMsg := "GetMsgBusEndpoint not yet implemented"
	fmt.Println(errMsg)
	return nil, errors.New(errMsg)
}

// GetLoggableEntities refers to the grpc northbound interface exposed for getting loggable entities from device
func (c *NativeHwManagementService) GetLoggableEntities(ctx context.Context, req *dmi.GetLoggableEntitiesRequest) (*dmi.GetLogLevelResponse, error) {

	ctx = config.GetNewContextFromContxt(ctx, "GetLoggableEntities")
	resp := new(dmi.GetLogLevelResponse)
	resp.DeviceUuid = req.DeviceUuid

	logger.Infow(ctx, "GetLoggableEntities-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "GetLoggableEntities-on-grpc-server-completed", log.Fields{"req": req})

	devRec, err := dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid)
	if err != nil {
		return errRespGetLoggableEntities(ctx, req, dmi.GetLogLevelResponse_UNKNOWN_DEVICE, err), err
	}

	if devRec.Logging.LoggableEntities != nil {
		resp.Status = dmi.Status_OK_STATUS
		resp.LogLevels, _ = devRec.GetLoggableEntitiesFromDevRec(ctx, nil)
		return resp, nil
	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		return errRespGetLoggableEntities(ctx, req, dmi.GetLogLevelResponse_DEVICE_UNREACHABLE, err), err
	}

	return adapter.GetLoggableEntities(ctx, devRec, req)
}

// SetLogLevel refers to the grpc northbound interface exposed for setting log level for entities
func (c *NativeHwManagementService) SetLogLevel(ctx context.Context, req *dmi.SetLogLevelRequest) (*dmi.SetLogLevelResponse, error) {

	ctx = config.GetNewContextFromContxt(ctx, "SetLogLevel")
	resp := new(dmi.SetLogLevelResponse)
	resp.DeviceUuid = req.DeviceUuid

	logger.Infow(ctx, "SetLogLevel-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "SetLogLevel-on-grpc-server-completed", log.Fields{"req": req})

	// validate request
	if ok, err := isValidSetLogLevel(ctx, req.Loglevels); ok {
		return errRespSetLogLevel(ctx, req, dmi.SetLogLevelResponse_UNKNOWN_LOG_ENTITY, err), err
	}

	devRec, err := dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid)
	if err != nil {
		return errRespSetLogLevel(ctx, req, dmi.SetLogLevelResponse_UNKNOWN_DEVICE, err), err
	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		return errRespSetLogLevel(ctx, req, dmi.SetLogLevelResponse_DEVICE_UNREACHABLE, err), err
	}

	return adapter.SetLogLevel(ctx, devRec, req)
}

// GetLogLevel refers to the grpc northbound interface exposed for getting log level of entities
func (c *NativeHwManagementService) GetLogLevel(ctx context.Context, req *dmi.GetLogLevelRequest) (*dmi.GetLogLevelResponse, error) {

	ctx = config.GetNewContextFromContxt(ctx, "GetLogLevel")
	resp := new(dmi.GetLogLevelResponse)
	resp.DeviceUuid = req.DeviceUuid

	logger.Infow(ctx, "GetLogLevel-invoked-on-grpc-server", log.Fields{"req": req})
	defer logger.Infow(ctx, "GetLogLevel-on-grpc-server-completed", log.Fields{"req": req})

	devRec, err := dev.DBGetByUuid(ctx, req.DeviceUuid.Uuid)
	if err != nil {
		return errRespGetLogLevel(ctx, req, dmi.GetLogLevelResponse_UNKNOWN_DEVICE, err), err
	}

	if devRec.Logging.LoggableEntities != nil {

		if output, ok := devRec.GetLoggableEntitiesFromDevRec(ctx, req.Entities); ok {
			resp.Status = dmi.Status_OK_STATUS
			resp.LogLevels = output
			return resp, nil
		}

	}

	adapter, err := connections.getConnection(ctx, devRec)
	if adapter == nil {
		return errRespGetLogLevel(ctx, req, dmi.GetLogLevelResponse_DEVICE_UNREACHABLE, err), err
	}

	return adapter.GetLogLevel(ctx, devRec, req)
}
