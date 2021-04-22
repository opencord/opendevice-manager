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

// Package main holds functions for init
package main

import (
	"context"

	"github.com/opencord/opendevice-manager/pkg/config"
	"github.com/opencord/opendevice-manager/pkg/nbi"

	"github.com/opencord/opendevice-manager/pkg/db"
	"github.com/opencord/opendevice-manager/pkg/msgbus"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opencord/voltha-lib-go/v4/pkg/log"
)

// logger represents the log object
var logger log.CLogger

// init function for the package
func init() {
	logger = config.Initlog()
}

func printBanner() {
	fmt.Println("╔═╗╔═╗╔═╗╔╗╔  ╔╦╗╔═╗╦  ╦╦╔═╗╔═╗  ╔╦╗╔═╗╔╗╔╔═╗╔═╗╔═╗╦═╗")
	fmt.Println("║ ║╠═╝║╣ ║║║   ║║║╣ ╚╗╔╝║║  ║╣   ║║║╠═╣║║║╠═╣║ ╦║╣ ╠╦╝")
	fmt.Println("╚═╝╩  ╚═╝╝╚╝  ═╩╝╚═╝ ╚╝ ╩╚═╝╚═╝  ╩ ╩╩ ╩╝╚╝╩ ╩╚═╝╚═╝╩╚═")
}

func exitBanner() {
	fmt.Println("╔═╗╦  ╔═╗╔═╗╦╔╗╔╔═╗")
	fmt.Println("║  ║  ║ ║╚═╗║║║║║ ╦")
	fmt.Println("╚═╝╩═╝╚═╝╚═╝╩╝╚╝╚═╝ooo")
}

func waitForExit(ctx context.Context) int {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-signalChannel
	switch s {
	case syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT:
		logger.Infow(ctx, "closing-signal-received", log.Fields{"signal": s})
		return 0
	default:
		logger.Infow(ctx, "unexpected-signal-received", log.Fields{"signal": s})
		return 1
	}
}

func main() {
	printBanner()
	coreFlags := config.NewCoreFlags()
	coreFlags.ParseEnv()
	core := config.NewCoreConfig()
	ctx := config.GetNewContextFromGlobalContxt("main-service")
	startServices(ctx, coreFlags, core)
}

func startServices(ctx context.Context, coreFlags *config.CoreFlags, core *config.Core) {

	logger.Infow(ctx, "starting-core-services", log.Fields{"core-flags": coreFlags})

	defer close(core.Stopped)

	defer core.Cancel()

	// setup kv client
	logger.Debugw(ctx, "create-kv-client", log.Fields{"kvstore": config.KVStore})

	kvClient, err := db.NewKVClient(ctx, config.KVStore, coreFlags.DbEndPoint, coreFlags.DbTimeout)

	if err != nil {
		logger.Fatal(ctx, log.Fields{"err": err})
	}

	defer db.StopKVClient(log.WithSpanFromContext(context.Background(), ctx), kvClient)

	err = msgbus.InitMsgbusProducer(ctx)

	if err != nil {
		logger.Fatal(ctx, log.Fields{"err": err})
	}

	defer nbi.DeInitConnectMap(ctx)

	defer msgbus.Close(ctx)

	go nbi.StartGrpcServer(ctx)

	defer nbi.StopGrpcServer(ctx)

	waitForExit(core.Ctx)
	exitBanner()
}
