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
	"net"
	"strconv"
	"strings"

	"github.com/opencord/device-management-interface/go/dmi"
	dev "github.com/opencord/opendevice-manager/pkg/models/device"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

/* validateUri() verifies if the ip and port are valid and already registered then return the truth value of the desired state specified by the following 2 switches,
   wantRegistered: 'true' if the fact of an ip is registered is the desired state
   includePort: 'true' further checks if <ip>:<port#> does exist in the devicemap in case an ip is found registered
*/
func validateUri(ctx context.Context, uri string) (ok bool, err error) {
	ok = false
	if !strings.Contains(uri, ":") {
		logger.Errorw(ctx, "Invalid-uri", log.Fields{"uri-received": uri, "expected-uri": "ip:port"})
		err = errors.New("incorrect IP address format (<ip>:<port #>)")
		return
	}
	splits := strings.Split(uri, ":")
	ip, port := splits[0], splits[1]
	if net.ParseIP(ip) == nil {
		// also check to see if it's a valid hostname
		if _, err2 := net.LookupIP(ip); err2 != nil {
			logger.Errorw(ctx, "Invalid-ip", log.Fields{"uri-received": uri, "ip": ip})
			err = errors.New("invalid IP address " + ip)
			return
		}
	}
	if _, err2 := strconv.Atoi(port); err2 != nil {
		logger.Errorw(ctx, "Invalid-port", log.Fields{"uri-received": uri, "port": port})
		err = errors.New("Port number " + port + " needs to be an integer")
		return
	}
	ok = true
	return
}

// validateStartManagingDeviceReq validates the 'StartManagingDevice' request is proper or not
func validateStartManagingDeviceReq(ctx context.Context, req *dmi.ModifiableComponent) (resp *dmi.StartManagingDeviceResponse, ok bool) {
	resp = new(dmi.StartManagingDeviceResponse)
	resp.DeviceUuid = new(dmi.Uuid)
	resp.Status = dmi.Status_ERROR_STATUS
	if ok1, err := validateUri(ctx, req.Uri.Uri); !ok1 {
		logger.Errorw(ctx, "validation-failed-for-StartManagingDevice-request", log.Fields{"error": err, "req": req})
		resp.Reason = dmi.StartManagingDeviceResponse_INVALID_PARAMS
		resp.ReasonDetail = err.Error()
		return
	}

	if rec, _ := dev.DBGetByName(ctx, req.Name); rec != nil {
		logger.Errorw(ctx, "validation-failed-for-StartManagingDevice-request-record-already-exists", log.Fields{"req": req, "rec": rec})
		resp.Reason = dmi.StartManagingDeviceResponse_DEVICE_ALREADY_MANAGED
		resp.ReasonDetail = "device already exists and managed with uuid " + rec.Uuid + " and uri " + rec.Uri
		return
	}

	ok = true

	return
}

// isValidSetLogLevel check is valid set loglevel request
func isValidSetLogLevel(ctx context.Context, listEntities []*dmi.EntitiesLogLevel) (bool, error) {

	if len(listEntities) == 0 {
		// atleast one entities is required for set loglevel
		logger.Errorw(ctx, "found-empty-entities", log.Fields{"entities": listEntities})
		return false, errors.New("found empty entities")
	}

	if len(listEntities) > 1 {
		// if set Entities more than 1, atleast 1 entity in nested struct
		for _, entities := range listEntities {
			if len(entities.Entities) == 0 {
				logger.Errorw(ctx, "entities-has-empty-entries", log.Fields{"entities": entities})
				return false, errors.New("set-empty-entries-not-allowed")
			}
		}
	}

	logger.Debug(ctx, "valid-set-log-request")
	return true, nil
}
