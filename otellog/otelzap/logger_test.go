// Copyright 2022 CloudWeGo Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otelzap

import (
	"bytes"
	"context"
	"testing"

	logging "github.com/kanhai-syd/hailog/logging"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func stdoutProvider(ctx context.Context) func() {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

// TestLogger test logger work with opentelemetry
func TestLogger(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)

	shutdown := stdoutProvider(ctx)
	defer shutdown()

	logger := NewLogger(
		WithTraceErrorSpanLevel(zap.WarnLevel),
		WithRecordStackTraceInSpan(true),
	)
	defer logger.Sync()

	logging.SetLogger(logger)
	logging.SetOutput(buf)
	logging.SetLevel(logging.LevelDebug)

	logging.Info("log from origin otelzap")
	assert.Contains(t, buf.String(), "log from origin otelzap")
	buf.Reset()

	tracer := otel.Tracer("test otel std logger")

	ctx, span := tracer.Start(ctx, "root")

	logging.CtxInfof(ctx, "hello %s", "world")
	assert.Contains(t, buf.String(), "trace_id")
	assert.Contains(t, buf.String(), "span_id")
	assert.Contains(t, buf.String(), "trace_flags")
	buf.Reset()

	span.End()

	ctx, child1 := tracer.Start(ctx, "child1")

	logging.CtxTracef(ctx, "trace %s", "this is a trace log")
	logging.CtxDebugf(ctx, "debug %s", "this is a debug log")
	logging.CtxInfof(ctx, "info %s", "this is a info log")

	child1.End()
	assert.Equal(t, codes.Unset, child1.(sdktrace.ReadOnlySpan).Status().Code)

	ctx, child2 := tracer.Start(ctx, "child2")
	logging.CtxNoticef(ctx, "notice %s", "this is a notice log")
	logging.CtxWarnf(ctx, "warn %s", "this is a warn log")
	logging.CtxErrorf(ctx, "error %s", "this is a error log")

	child2.End()
	assert.Equal(t, codes.Error, child2.(sdktrace.ReadOnlySpan).Status().Code)

	_, errSpan := tracer.Start(ctx, "error")

	logging.Info("no trace context")

	errSpan.End()
}

// TestLogLevel test SetLevel
func TestLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(
		WithTraceErrorSpanLevel(zap.WarnLevel),
		WithRecordStackTraceInSpan(true),
	)
	defer logger.Sync()

	// output to buffer
	logger.SetOutput(buf)

	logging.SetLogger(logger)
	logging.Debug("this is a debug log")
	assert.NotContains(t, buf.String(), "this is a debug log")

	logger.SetLevel(logging.LevelDebug)

	logging.Debugf("this is a debug log %s", "msg")
	assert.Contains(t, buf.String(), "this is a debug log")
}

func TestZapConfigLogger(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(
		WithCoreEnc(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())),
		WithCoreWs(zapcore.AddSync(&bytes.Buffer{})),
		WithCoreLevel(zap.NewAtomicLevelAt(zap.DebugLevel)),
		WithCustomFields("key1", "value1", "key2", "value2"),
		WithZapOptions(zap.AddCaller()),
	)
	defer logger.Sync()

	logger.SetOutput(buf)

	assert.NotNil(t, logger)

	logging.SetLogger(logger)
	logging.Info("this is a info log")
	
	assert.Contains(t, buf.String(), "this is a info log")
	assert.NotContains(t, buf.String(), "Ignored key without a value.")
}
