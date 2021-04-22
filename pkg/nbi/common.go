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
	"net"
	"os"

	"github.com/opencord/opendevice-manager/pkg/config"

	"github.com/opencord/device-management-interface/go/dmi"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// logger represents the log object
var logger log.CLogger

// grpcServer refers to object which holds grpc serevr
var grpcServer *g.Server

// init function for the package
func init() {
	logger = config.Initlog()
	initConnectMap()
}

func registerServers(grpcServer *g.Server) {
	hwMgmtSvc := new(NativeHwManagementService)
	dmi.RegisterNativeHWManagementServiceServer(grpcServer, hwMgmtSvc)
	reflection.Register(grpcServer)
}

// StartGrpcServer starts the grpc server for listening to NEM requests
func StartGrpcServer(ctx context.Context) {
	coreFlags := config.GetCoreFlags()
	lis, err := net.Listen("tcp", coreFlags.GrpcEndPoint)
	if err != nil {
		logger.Error(ctx, "Failed-to-listen-on-Grpc-Port", log.Fields{"grpc-flags": coreFlags.GrpcFlags, "error": err})
		os.Exit(1)
	}

	if coreFlags.SecureConnection {
		creds, err := credentials.NewServerTLSFromFile(coreFlags.ServerCrt, coreFlags.ServerKey)
		if err != nil {
			logger.Error(ctx, "could-not-process-the-credentials", log.Fields{"error": err})
		}
		grpcServer = g.NewServer(g.Creds(creds))
	} else {
		grpcServer = g.NewServer()
	}

	registerServers(grpcServer)

	logger.Infow(ctx, "Grpc-server-starting", log.Fields{"grpc-server-info": grpcServer, "grpc-env-info": coreFlags.GrpcFlags, "is-secure-conn": coreFlags.SecureConnection})

	// Starting the server
	if err := grpcServer.Serve(lis); err != nil {
		logger.Errorw(ctx, "Failed-to-start-Grpc-Server", log.Fields{"server": coreFlags.GrpcFlags, "error": err})
		os.Exit(1)
	}

	logger.Infow(ctx, "grpc-server-stopped-successfully", log.Fields{"grpc-server-info": grpcServer, "grpc-env-info": coreFlags.GrpcFlags})

}

// StopGrpcServer tear down the gRPC connection from opendevice manager to NEM
func StopGrpcServer(ctx context.Context) {
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	logger.Infow(ctx, "grpc-server-teardown-success", log.Fields{"grpc-server-info": grpcServer})
}
