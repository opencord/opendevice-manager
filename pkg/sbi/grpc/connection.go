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

// Package grpc holds utils for grpc client implementation
package grpc

import (
	"context"

	"github.com/opencord/device-management-interface/go/dmi"
	"github.com/opencord/opendevice-manager/pkg/config"
	"github.com/opencord/voltha-lib-go/v4/pkg/log"
	"google.golang.org/grpc/credentials"

	g "google.golang.org/grpc"
)

// logger represents the log object
var logger log.CLogger

// init function for the package
func init() {
	logger = config.Initlog()
}

// Client holds the parameters for grpc
type Client struct {
	uri             string
	conn            *g.ClientConn
	hwMgmtSvcClient dmi.NativeHWManagementServiceClient
}

// NewClient returns a new Grpc Client
func NewClient(uri string) *Client {
	c := new(Client)
	c.uri = uri
	return c
}

func (c *Client) getDialOpts(ctx context.Context) []g.DialOption {

	coreFlags := config.NewCoreFlags()
	var opts []g.DialOption

	if coreFlags.SecureConnection {
		logger.Info(ctx, "Trying-to-establish-secure-connection")

		creds, err := credentials.NewClientTLSFromFile(coreFlags.CertsPath.RootCaCrt, "")
		if err != nil {
			logger.Fatalf(ctx, "could-not-process-the-credentials", log.Fields{"err": err})
		}

		err = creds.OverrideServerName(coreFlags.GrpcHostName)
		if err != nil {
			logger.Fatalf(ctx, "Overriding-server-name-failed-at-getDialOpts()", log.Fields{"err": err})
		}

		opts = append(opts, g.WithTransportCredentials(creds))
	} else {
		logger.Info(ctx, "Trying-to-establish-insecure-connection")
		opts = append(opts, g.WithInsecure())
	}

	opts = append(opts, g.WithTimeout(coreFlags.GrpcRetryInterval))
	backoffConfig := g.BackoffConfig{MaxDelay: coreFlags.GrpcBackoffMaxDelay}
	opts = append(opts, g.WithBackoffConfig(backoffConfig))

	return opts
}

// Connect will establish a connection
func (c *Client) Connect(ctx context.Context) error {
	logger.Info(ctx, "Invoked-connectGrpcServer")
	// log.Info("Invoked-connectGrpcServer", log.Opts{"peer-ID": peerID})
	opts := c.getDialOpts(ctx)
	// Establishing the server connection
	conn, err := g.Dial(c.uri, opts...)
	if err != nil {
		logger.Error(ctx, "Grpc-client-connection-failed", log.Fields{"error": err})
		return err
	}
	c.conn = conn
	logger.Info(ctx, "Connection-established", log.Fields{"conn": conn})
	// Constructing a client object
	c.hwMgmtSvcClient = dmi.NewNativeHWManagementServiceClient(conn)
	return nil
}

// Disconnect will remove the connection
func (c *Client) Disconnect(ctx context.Context) error {
	logger.Infow(ctx, "Invoked-Disconnect", log.Fields{"client": c})
	return c.conn.Close()
}
