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

	"github.com/opencord/device-management-interface/go/dmi"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// errRespStartManagingDevice represents the func to construct error response for rpc StartManagingDevice
func errRespStartManagingDevice(ctx context.Context, req *dmi.ModifiableComponent, reason dmi.StartManagingDeviceResponse_Reason, err error) *dmi.StartManagingDeviceResponse {
	resp := new(dmi.StartManagingDeviceResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "StartManagingDevice-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespStopManagingDevice represents the func to construct error response for rpc StopManagingDevice
func errRespStopManagingDevice(ctx context.Context, req *dmi.StopManagingDeviceRequest, reason dmi.StopManagingDeviceResponse_Reason, err error) *dmi.StopManagingDeviceResponse {
	resp := new(dmi.StopManagingDeviceResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "StopManagingDevice-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespGetPhysicatInventory represents the func to construct error response for rpc GetPhysicatInventory
func errRespGetPhysicatInventory(ctx context.Context, req *dmi.PhysicalInventoryRequest, reason dmi.PhysicalInventoryResponse_Reason, err error) *dmi.PhysicalInventoryResponse {
	resp := new(dmi.PhysicalInventoryResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "GetPhysicatInventory-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespGetHWComponentInfo represents the func to construct error response for rpc GetHWComponentInfo
func errRespGetHWComponentInfo(ctx context.Context, req *dmi.HWComponentInfoGetRequest, reason dmi.HWComponentInfoGetResponse_Reason, err error) *dmi.HWComponentInfoGetResponse {
	resp := new(dmi.HWComponentInfoGetResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "GetHWComponentInfo-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespSetHWComponentInfo represents the func to construct error response for rpc SetHWComponentInfo
func errRespSetHWComponentInfo(ctx context.Context, req *dmi.HWComponentInfoSetRequest, reason dmi.HWComponentInfoSetResponse_Reason, err error) *dmi.HWComponentInfoSetResponse {
	resp := new(dmi.HWComponentInfoSetResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "SetHWComponentInfo-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespSetLoggingEndpoint represents the func to construct error response for rpc SetLoggingEndpoint
func errRespSetLoggingEndpoint(ctx context.Context, req *dmi.SetLoggingEndpointRequest, reason dmi.SetRemoteEndpointResponse_Reason, err error) *dmi.SetRemoteEndpointResponse {
	resp := new(dmi.SetRemoteEndpointResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "SetLoggingEndpoint-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespGetLoggingEndpoint represents the func to construct error response for rpc GetLoggingEndpoint
func errRespGetLoggingEndpoint(ctx context.Context, req *dmi.HardwareID, reason dmi.GetLoggingEndpointResponse_Reason, err error) *dmi.GetLoggingEndpointResponse {
	resp := new(dmi.GetLoggingEndpointResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "GetLoggingEndpoint-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespGetLoggableEntities represents the func to construct error response for rpc GetLoggableEntities
func errRespGetLoggableEntities(ctx context.Context, req *dmi.GetLoggableEntitiesRequest, reason dmi.GetLogLevelResponse_Reason, err error) *dmi.GetLogLevelResponse {
	resp := new(dmi.GetLogLevelResponse)
	resp.DeviceUuid = req.DeviceUuid
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "GetLoggableEntities-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespSetLogLevel represents the func to construct error response for rpc SetLogLevel
func errRespSetLogLevel(ctx context.Context, req *dmi.SetLogLevelRequest, reason dmi.SetLogLevelResponse_Reason, err error) *dmi.SetLogLevelResponse {
	resp := new(dmi.SetLogLevelResponse)
	resp.DeviceUuid = req.DeviceUuid
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "SetLogLevel-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}

// errRespGetLogLevel represents the func to construct error response for rpc GetLogLevel
func errRespGetLogLevel(ctx context.Context, req *dmi.GetLogLevelRequest, reason dmi.GetLogLevelResponse_Reason, err error) *dmi.GetLogLevelResponse {
	resp := new(dmi.GetLogLevelResponse)
	resp.DeviceUuid = req.DeviceUuid
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = reason
	resp.ReasonDetail = err.Error()
	logger.Errorw(ctx, "GetLogLevel-on-grpc-server-failed", log.Fields{"req": req, "resp": resp, "error": err})
	return resp
}
