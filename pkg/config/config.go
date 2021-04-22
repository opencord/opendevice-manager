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

// Package config Common Logger initialization
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Constants defined for environment variables
const (
	envLogLevel            = "LOG_LEVEL"
	envMsgbusEndPoint      = "MSGBUS_END_POINT"
	envMsgbusRetryInterval = "MSGBUS_RETRY_INTERVAL"
	envDbEndPoint          = "DB_END_POINT"
	envDbTimeout           = "DB_TIMEOUT"
	envGrpcEndPoint        = "GRPC_END_POINT"
	envGrpcRetryInterval   = "GRPC_RETRY_INTERVAL"
	envGrpcBackoffMaxDelay = "GRPC_BACKOFF_MAX_DELAY"
	envGrpcMaxRetryCount   = "GRPC_MAX_RETRY_COUNT"
	envSecureConnection    = "SECURE_GRPC"
)

// Default values defined if not provided through environment variables
const (
	defaultLogLevel            = LogLevelDebug
	defaultMsgbusEndPoint      = "127.0.0.1:9092"
	defaultMsgbusRetryInterval = 10 * time.Second
	defaultDbEndPoint          = "127.0.0.1:2379"
	defaultDbTimeout           = 5 * time.Second
	defaultGrpcEndPoint        = "0.0.0.0:9292"
	defaultGrpcRetryInterval   = 10 * time.Second
	defaultGrpcBackoffMaxDelay = 30 * time.Second
	defaultGrpcMaxRetryCount   = 5
	defaultGrpcHostName        = "DMI"
	defaultSecureConnection    = false
)

// Constants defined for certiifcates
const (
	pathRootCaCrt = "/etc/Root_CA.crt"
	pathServerCrt = "/etc/grpc_server.crt"
	pathServerKey = "/etc/grpc_server.key"
)

// DB versions
const (
	DBVer1 = "v1"
)

// Constants defined for Db
const (
	KVStore  = "etcd"
	DBPrefix = "/OpenDevMgr/"
	CurDBVer = DBVer1
)

// Constants defined for Msgbus Topic to receive messages
const (
	OpenDevMgrEventsTopic   = "dm.events"
	OpenDevMgrMetricsTopic  = "dm.metrics"
	OpenDevMgrEventsMsgType = "OPEN-DEV-MGR-EVENTS-MSG"
)

// Log Level Constants
const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

type correlationIdType int8

// Constants defined for context management
const (
	RequestIdKey correlationIdType = iota
	SessionIdKey
)

// ContextId constant used to print context id in logs
const (
	ContextId = "context-id"
)

var coreFlags *CoreFlags

// NewCoreFlags creates a new CoreFlag object for you
func NewCoreFlags() *CoreFlags {
	if coreFlags != nil {
		return coreFlags
	}
	coreFlags = new(CoreFlags)
	coreFlags.LogLevel = defaultLogLevel
	coreFlags.MsgbusEndPoint = defaultMsgbusEndPoint
	coreFlags.MsgbusRetryInterval = defaultMsgbusRetryInterval
	coreFlags.DbEndPoint = defaultDbEndPoint
	coreFlags.DbTimeout = defaultDbTimeout
	coreFlags.SecureConnection = defaultSecureConnection
	coreFlags.GrpcFlags.GrpcEndPoint = defaultGrpcEndPoint
	coreFlags.GrpcFlags.GrpcRetryInterval = defaultGrpcRetryInterval
	coreFlags.GrpcFlags.GrpcBackoffMaxDelay = defaultGrpcBackoffMaxDelay
	coreFlags.GrpcFlags.GrpcMaxRetryCount = defaultGrpcMaxRetryCount
	coreFlags.GrpcFlags.GrpcHostName = defaultGrpcHostName
	coreFlags.CertsPath.RootCaCrt = pathRootCaCrt
	coreFlags.CertsPath.ServerCrt = pathServerCrt
	coreFlags.CertsPath.ServerKey = pathServerKey
	return coreFlags
}

// GetCoreFlags returns the entire config values
func GetCoreFlags() *CoreFlags {
	if coreFlags != nil {
		return coreFlags
	}
	return nil
}

// GetGrpcFlags returns the grpc config values
func GetGrpcFlags() *GrpcFlags {
	if coreFlags != nil {
		return &coreFlags.GrpcFlags
	}
	return nil
}

// GrpcFlags is struct defined to store grpc parameters
type GrpcFlags struct {
	GrpcEndPoint        string
	GrpcRetryInterval   time.Duration
	GrpcBackoffMaxDelay time.Duration
	GrpcMaxRetryCount   int
	GrpcHostName        string
}

// CertsPath is struct defined to store certificates path
type CertsPath struct {
	RootCaCrt string
	ServerCrt string
	ServerKey string
}

// CoreFlags is a structure defined to store all configurations
type CoreFlags struct {
	LogLevel            string
	MsgbusEndPoint      string
	MsgbusRetryInterval time.Duration
	DbEndPoint          string
	DbTimeout           time.Duration
	SecureConnection    bool
	GrpcFlags
	CertsPath
}

// ParseEnv method retrieves environment variables passed and replaces with
// corresponding default variables stored in the CoreFlags object
func (cf *CoreFlags) ParseEnv() {

	if env := os.Getenv(envLogLevel); env != "" {
		cf.LogLevel = env
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envLogLevel, cf.LogLevel)

	if env := os.Getenv(envMsgbusEndPoint); env != "" {
		cf.MsgbusEndPoint = env
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envMsgbusEndPoint, cf.MsgbusEndPoint)

	if env := os.Getenv(envMsgbusRetryInterval); env != "" {
		interval, err := strconv.Atoi(env)
		if err == nil {
			cf.MsgbusRetryInterval = time.Duration(interval) * time.Second
		} else {
			fmt.Printf("Invalid value '%s' passed for '%s'. Taking the default value.\n", env, envMsgbusRetryInterval)
		}
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envMsgbusRetryInterval, cf.MsgbusRetryInterval)

	if env := os.Getenv(envDbEndPoint); env != "" {
		cf.DbEndPoint = env
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envDbEndPoint, cf.DbEndPoint)

	if env := os.Getenv(envDbTimeout); env != "" {
		interval, err := strconv.Atoi(env)
		if err == nil {
			cf.DbTimeout = time.Duration(interval) * time.Second
		} else {
			fmt.Printf("Invalid value '%s' passed for '%s'. Taking the default value.\n", env, envDbTimeout)
		}
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envDbTimeout, cf.DbTimeout)

	if env := os.Getenv(envGrpcEndPoint); env != "" {
		cf.GrpcFlags.GrpcEndPoint = env
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envGrpcEndPoint, cf.GrpcFlags.GrpcEndPoint)

	if env := os.Getenv(envGrpcRetryInterval); env != "" {
		interval, err := strconv.Atoi(env)
		if err == nil {
			cf.GrpcFlags.GrpcRetryInterval = time.Duration(interval) * time.Second
		} else {
			fmt.Printf("Invalid value '%s' passed for '%s'. Taking the default value.\n", env, envGrpcRetryInterval)
		}
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envGrpcRetryInterval, cf.GrpcFlags.GrpcRetryInterval)

	if env := os.Getenv(envGrpcBackoffMaxDelay); env != "" {
		interval, err := strconv.Atoi(env)
		if err == nil {
			cf.GrpcFlags.GrpcBackoffMaxDelay = time.Duration(interval) * time.Second
		} else {
			fmt.Printf("Invalid value '%s' passed for '%s'. Taking the default value.\n", env, envGrpcBackoffMaxDelay)
		}
	}
	fmt.Printf("Environment variable '%s' setting to : %s\n", envGrpcBackoffMaxDelay, cf.GrpcFlags.GrpcBackoffMaxDelay)

	if env := os.Getenv(envGrpcMaxRetryCount); env != "" {
		maxRetry, err := strconv.Atoi(env)
		if err == nil {
			cf.GrpcFlags.GrpcMaxRetryCount = maxRetry
		} else {
			fmt.Printf("Invalid value '%s' passed for '%s'. Taking the default value.\n", env, envGrpcMaxRetryCount)
		}
	}
	fmt.Printf("Environment variable '%s' setting to : %v\n", envGrpcMaxRetryCount, cf.GrpcFlags.GrpcMaxRetryCount)

	if env := os.Getenv(envSecureConnection); env != "" {
		secureCon, err := strconv.ParseBool(env)
		if err == nil {
			cf.SecureConnection = secureCon
		} else {
			fmt.Printf("Invalid value '%s' passed for '%s'. Taking the default value.\n", env, envSecureConnection)
		}
	}
	fmt.Printf("Environment variable '%s' setting to : %v\n", envSecureConnection, cf.SecureConnection)

}
