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
	"context"
	"fmt"

	"github.com/google/uuid"
)

var ctxGlobal context.Context

// Core represent read,write core attributes
type Core struct {
	Ctx     context.Context
	Cancel  context.CancelFunc
	Stopped chan struct{}
}

// NewCoreConfig will give the new core configuration
func NewCoreConfig() *Core {
	ctx, cancelCtx := context.WithCancel(context.Background())
	core := &Core{Ctx: ctx, Cancel: cancelCtx, Stopped: make(chan struct{})}
	ctxGlobal = ctx
	return core
}

// GetContext returns global context
func GetContext() context.Context {
	return ctxGlobal
}

func retConstructedContxt(ctx context.Context, msg string) context.Context {
	id, _ := uuid.NewRandom()
	var reqID string
	if msg != "" {
		reqID = fmt.Sprintf("%s-%v", msg, id)
	} else {
		reqID = fmt.Sprintf("%v", id)
	}
	ctx2 := context.WithValue(ctx, RequestIdKey, reqID)
	return ctx2
}

// GetNewContextFromGlobalContxt returns context from global context
func GetNewContextFromGlobalContxt(msg string) context.Context {
	return retConstructedContxt(ctxGlobal, msg)
}

// GetNewContextFromContxt returns new context from passed context
func GetNewContextFromContxt(ctx context.Context, msg string) context.Context {
	return retConstructedContxt(ctx, msg)
}
