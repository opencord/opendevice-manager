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

// Package grpc holds utils for grpc Client implementation
package grpc

import (
	"context"
	"errors"
	"io"

	"github.com/opencord/device-management-interface/go/dmi"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"

	copy "github.com/jinzhu/copier"
	dev "github.com/opencord/opendevice-manager/pkg/models/device"
	hw "github.com/opencord/opendevice-manager/pkg/models/hwcomponents"
)

// StartManagingDevice is the adapter implementation for start managing device in grpc adapter layer
func (c *Client) StartManagingDevice(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.ModifiableComponent, streamResp dmi.NativeHWManagementService_StartManagingDeviceServer) (error, bool) {

	var connMade bool

	resp := new(dmi.StartManagingDeviceResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = dmi.StartManagingDeviceResponse_INTERNAL_ERROR

	logger.Info(ctx, "Invoked-StartManagingDevice-at-grpc-adapter", log.Fields{"req": req})

	stream, err := c.hwMgmtSvcClient.StartManagingDevice(ctx, req)

	if err != nil {
		logger.Error(ctx, "error-at-StartManagingDevice")
		resp.ReasonDetail = err.Error()
		streamResp.Send(resp)
		return errors.New("RPC Failed for StartManagingDevice"), connMade
	}

	for {

		respFromDev, err := stream.Recv()

		if err == io.EOF {
			logger.Info(ctx, "Exiting-StartManagingDevice-on-connection-break-from-grpc-server", log.Fields{"req": req, "error": err})
			break
		}

		if err != nil {
			logger.Error(ctx, "Failed-at-StartManagingDevice-while-receiving-server-response", log.Fields{"error": err, "req": req})
			return err, connMade
		}

		if respFromDev.Status == dmi.Status_OK_STATUS {
			connMade = true
			devRec.Uuid = respFromDev.DeviceUuid.Uuid
			// Store in DB
			devRec.DBAddByName(ctx)
			devRec.DBAddUuidLookup(ctx)
			logger.Infow(ctx, "received-success-response-from-dm-agent-for-StartManagingDevice-req", log.Fields{"req": req, "resp": respFromDev})
		} else if respFromDev.Status == dmi.Status_ERROR_STATUS || err != nil {
			logger.Errorw(ctx, "received-failed-response-from-dm-agent-for-StartManagingDevice-req", log.Fields{"req": req, "resp": respFromDev})
			if err == nil {
				err = errors.New(respFromDev.ReasonDetail)
			}
			streamResp.Send(respFromDev)
			return err, connMade
		}

		streamResp.Send(respFromDev)
	}

	return nil, connMade

}

// StopManagingDevice is the adapter implementation for stop managing device in grpc adapter layer
func (c *Client) StopManagingDevice(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.StopManagingDeviceRequest) (*dmi.StopManagingDeviceResponse, error) {

	logger.Info(ctx, "Invoked-StopManagingDevice-at-grpc-adapter", log.Fields{"req": req})

	return c.hwMgmtSvcClient.StopManagingDevice(ctx, req)

}

// SetLoggingEndpoint is the adapter implementation for set the location to which logs need to be shipped in grpc adapter layer
func (c *Client) SetLoggingEndpoint(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.SetLoggingEndpointRequest) (*dmi.SetRemoteEndpointResponse, error) {
	logger.Info(ctx, "Invoked-SetLoggingEndpoint", log.Fields{"req": req})

	resp, err := c.hwMgmtSvcClient.SetLoggingEndpoint(ctx, req)
	if err != nil {
		logger.Error(ctx, "error-at-SetLoggingEndpoint")
		return resp, err
	}
	err = errors.New(resp.ReasonDetail)

	if resp.Status == dmi.Status_OK_STATUS {
		devRec.Logging.EndPoint = req.LoggingEndpoint
		devRec.Logging.Protocol = req.LoggingProtocol
		// Store in DB
		devRec.DBAddByName(ctx)
		logger.Infow(ctx, "received-success-response-from-dm-agent-for-SetLoggingEndpoint-req", log.Fields{"req": req, "resp": resp})
	} else {
		logger.Errorw(ctx, "received-failed-response-from-dm-agent-for-SetLoggingEndpoint-req", log.Fields{"req": req, "resp": resp})
	}
	return resp, err
}

// GetLoggingEndpoint is the adapter implementation for get the location to which logs need to be shipped in grpc adapter layer
func (c *Client) GetLoggingEndpoint(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.HardwareID) (*dmi.GetLoggingEndpointResponse, error) {
	logger.Info(ctx, "Invoked-GetLoggingEndpoint", log.Fields{"req": req})

	resp, err := c.hwMgmtSvcClient.GetLoggingEndpoint(ctx, req)
	if err != nil {
		logger.Error(ctx, "error-at-GetLoggingEndpoint")
		return resp, err
	}
	err = errors.New(resp.ReasonDetail)

	if resp.Status == dmi.Status_OK_STATUS {
		devRec.Logging.EndPoint = resp.LoggingEndpoint
		devRec.Logging.Protocol = resp.LoggingProtocol
		// Store in DB
		devRec.DBAddByName(ctx)
		logger.Infow(ctx, "received-success-response-from-dm-agent-for-GetLoggingEndpoint-req", log.Fields{"req": req, "resp": resp})
	} else {
		logger.Errorw(ctx, "received-failed-response-from-dm-agent-for-GetLoggingEndpoint-req", log.Fields{"req": req, "resp": resp})
	}
	return resp, err
}

// GetPhysicalInventory is the adapter implementation for reading physical inventories in grpc adapter layer
func (c *Client) GetPhysicalInventory(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.PhysicalInventoryRequest, streamResp dmi.NativeHWManagementService_GetPhysicalInventoryServer) error {

	logger.Info(ctx, "Invoked-GetPhysicalInventory-at-grpc-adapter", log.Fields{"req": req})

	resp := new(dmi.PhysicalInventoryResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = dmi.PhysicalInventoryResponse_INTERNAL_ERROR

	stream, err := c.hwMgmtSvcClient.GetPhysicalInventory(ctx, req)

	if err != nil {
		logger.Error(ctx, "error-at-GetPhysicalInventory", log.Fields{"error": err})
		resp.ReasonDetail = err.Error()
		streamResp.Send(resp)
		return err
	}

	for {

		respFromDev, err := stream.Recv()

		if err == io.EOF {
			logger.Info(ctx, "Exiting-GetPhysicalInventory-on-connection-break-from-grpc-server", log.Fields{"req": req, "error": err})
			break
		}

		if err != nil {
			logger.Error(ctx, "Failed-at-GetPhysicalInventory-while-receiving-server-response", log.Fields{"error": err, "request": req})
			return err
		}

		if respFromDev.Status == dmi.Status_OK_STATUS {

			// Store in DB
			devRec.DBSaveHwInfo(ctx, respFromDev.Inventory)
			nameToUuidMap := make(map[string]string)
			hw.DBSaveHwCompsFromPhysicalInventory(ctx, devRec.Uuid, nameToUuidMap, respFromDev.Inventory.Root.Children)
			hw.DBAddNameToUuidlookup(ctx, devRec.Uuid, nameToUuidMap)
			logger.Infow(ctx, "received-success-response-from-dm-agent-for-GetPhysicalInventory-req", log.Fields{"req": req, "resp": respFromDev})

		} else if respFromDev.Status == dmi.Status_ERROR_STATUS || err != nil {

			logger.Errorw(ctx, "received-failed-response-from-dm-agent-for-GetPhysicalInventory-req", log.Fields{"req": req, "resp": respFromDev})
			if err == nil {
				err = errors.New(respFromDev.ReasonDetail)
			}
			streamResp.Send(respFromDev)
			return err

		}

		streamResp.Send(respFromDev)
	}

	return nil
}

// GetLoggableEntities is the adapter implementation for reading physical inventories in grpc adapter layer
func (c *Client) GetLoggableEntities(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.GetLoggableEntitiesRequest) (*dmi.GetLogLevelResponse, error) {

	logger.Info(ctx, "Invoked-GetLoggableEntities-at-grpc-adapter", log.Fields{"req": req})

	resp, err := c.hwMgmtSvcClient.GetLoggableEntities(ctx, req)

	if err != nil {
		logger.Error(ctx, "Failed-at-GetLoggableEntities-while-receiving-server-response", log.Fields{"error": err, "request": req})
		return resp, err
	}
	// update the db if get log response is success from device
	if resp.Status == dmi.Status_OK_STATUS {
		devRec.SaveLoggableEntities(ctx, resp.LogLevels)
		devRec.DBAddByName(ctx)
	}

	return resp, err
}

// SetLogLevel is the adapter implementation for reading physical inventories in grpc adapter layer
func (c *Client) SetLogLevel(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.SetLogLevelRequest) (*dmi.SetLogLevelResponse, error) {

	logger.Info(ctx, "Invoked-SetLogLevel-at-grpc-adapter", log.Fields{"req": req})

	resp, err := c.hwMgmtSvcClient.SetLogLevel(ctx, req)

	if err != nil {
		logger.Error(ctx, "Failed-at-SetLogLevel-while-receiving-server-response", log.Fields{"error": err, "request": req})
		return resp, err
	}
	// update the db if setting log response is success from device
	if resp.Status == dmi.Status_OK_STATUS {
		devRec.SaveLoggableEntities(ctx, req.Loglevels)
		devRec.DBAddByName(ctx)
	}

	return resp, err
}

// GetLogLevel is the adapter implementation for reading physical inventories in grpc adapter layer
func (c *Client) GetLogLevel(ctx context.Context, devRec *dev.DeviceRecord, req *dmi.GetLogLevelRequest) (*dmi.GetLogLevelResponse, error) {

	logger.Info(ctx, "Invoked-GetLogLevel-at-grpc-adapter", log.Fields{"req": req})

	resp, err := c.hwMgmtSvcClient.GetLogLevel(ctx, req)

	if err != nil {
		logger.Error(ctx, "Failed-at-GetLogLevel-while-receiving-server-response", log.Fields{"error": err, "request": req})
		return resp, err
	}

	// update the db if get log response is success from device
	if resp.Status == dmi.Status_OK_STATUS {
		devRec.SaveLoggableEntities(ctx, resp.LogLevels)
		devRec.DBAddByName(ctx)
	}

	return resp, err
}

// GetHWComponentInfo is the adapter implementation for get the details of a particular HW component
func (c *Client) GetHWComponentInfo(ctx context.Context, deviceUuid string, hwCompRec *hw.HwCompRecord, req *dmi.HWComponentInfoGetRequest, streamResp dmi.NativeHWManagementService_GetHWComponentInfoServer) error {

	logger.Info(ctx, "Invoked-GetHWComponentInfo-at-grpc-adapter", log.Fields{"req": req})

	resp := new(dmi.HWComponentInfoGetResponse)
	resp.Status = dmi.Status_ERROR_STATUS
	resp.Reason = dmi.HWComponentInfoGetResponse_INTERNAL_ERROR

	stream, err := c.hwMgmtSvcClient.GetHWComponentInfo(ctx, req)

	if err != nil {
		logger.Error(ctx, "error-at-GetHWComponentInfo", log.Fields{"error": err})
		resp.ReasonDetail = err.Error()
		streamResp.Send(resp)
		return err
	}

	for {

		respFromDev, err := stream.Recv()

		if err == io.EOF {
			logger.Info(ctx, "Exiting-GetHWComponentInfo-on-connection-break-from-grpc-server", log.Fields{"req": req, "error": err})
			break
		}

		if err != nil {
			logger.Error(ctx, "Failed-at-GetHWComponentInfo-while-receiving-server-response", log.Fields{"error": err, "request": req})
			return err
		}

		if respFromDev.Status == dmi.Status_OK_STATUS {
			if hwCompRec.State == nil {
				hwCompRec.State = new(dmi.ComponentState)
			}
			err = copy.Copy(&hwCompRec, &respFromDev.Component)
			if err != nil {
				logger.Error(ctx, "Failed-at-GetHWComponentInfo-copy-failed", log.Fields{"error": err, "resp": respFromDev.Component})
			}
			if respFromDev.Component.State != nil {
				hwCompRec.State.AdminState = respFromDev.Component.State.AdminState
			}
			// Store in DB
			hwCompRec.DBAddByUuid(ctx, deviceUuid)
			logger.Infow(ctx, "received-success-response-from-dm-agent-for-GetHWComponentInfo-req", log.Fields{"req": req, "resp": respFromDev})
		} else if respFromDev.Status == dmi.Status_ERROR_STATUS || err != nil {
			logger.Errorw(ctx, "received-failed-response-from-dm-agent-for-GetHWComponentInfo-req", log.Fields{"req": req, "resp": respFromDev})
			if err == nil {
				err = errors.New(respFromDev.ReasonDetail)
			}
			streamResp.Send(respFromDev)
			return err
		}

		streamResp.Send(respFromDev)
	}

	return nil
}

// SetHWComponentInfo method is the grpc adapter implementation for setting hw component info on device
func (c *Client) SetHWComponentInfo(ctx context.Context, deviceUuid string, hwCompRec *hw.HwCompRecord, req *dmi.HWComponentInfoSetRequest) (*dmi.HWComponentInfoSetResponse, error) {
	logger.Info(ctx, "Invoked-SetHWComponentInfo", log.Fields{"req": req})

	resp, err := c.hwMgmtSvcClient.SetHWComponentInfo(ctx, req)

	if err != nil {
		logger.Error(ctx, "error-at-SetHWComponentInfo", log.Fields{"req": req, "error": err})
		return resp, err
	}

	if resp.Status == dmi.Status_OK_STATUS {
		hwCompRec.State = new(dmi.ComponentState)
		err = copy.Copy(&hwCompRec, &req.Changes)
		hwCompRec.State.AdminState = req.Changes.AdminState
		// Store in DB
		hwCompRec.DBAddByUuid(ctx, deviceUuid)
		logger.Infow(ctx, "received-success-response-from-dm-agent-for-SetHWComponentInfo-req", log.Fields{"req": req, "resp": resp})
	} else {
		logger.Errorw(ctx, "received-failed-response-from-dm-agent-for-SetHWComponentInfo-req", log.Fields{"req": req, "resp": resp})
		err = errors.New(resp.ReasonDetail)
	}

	return resp, err
}
