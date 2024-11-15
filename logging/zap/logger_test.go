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

package zap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kanhai-syd/hailog/logging"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// testEncoderConfig encoder config for testing, copy from otelzap
func testEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// humanEncoderConfig copy from otelzap
func humanEncoderConfig() zapcore.EncoderConfig {
	cfg := testEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	return cfg
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(filepath.Dir(file), 0o744)
	}

	f, _ := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)

	return zapcore.AddSync(f)
}

// TestLogger test logger work with otelhertz
func TestLogger(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(WithZapOptions(zap.WithFatalHook(zapcore.WriteThenPanic)))
	defer logger.Sync()

	logging.SetLogger(logger)
	logging.SetOutput(buf)
	logging.SetLevel(logging.LevelDebug)

	type logMap map[string]string

	logTestSlice := []logMap{
		{
			"logMessage":       "this is a trace log",
			"formatLogMessage": "this is a trace log: %s",
			"logLevel":         "Trace",
			"zapLogLevel":      "debug",
		},
		{
			"logMessage":       "this is a debug log",
			"formatLogMessage": "this is a debug log: %s",
			"logLevel":         "Debug",
			"zapLogLevel":      "debug",
		},
		{
			"logMessage":       "this is a info log",
			"formatLogMessage": "this is a info log: %s",
			"logLevel":         "Info",
			"zapLogLevel":      "info",
		},
		{
			"logMessage":       "this is a notice log",
			"formatLogMessage": "this is a notice log: %s",
			"logLevel":         "Notice",
			"zapLogLevel":      "warn",
		},
		{
			"logMessage":       "this is a warn log",
			"formatLogMessage": "this is a warn log: %s",
			"logLevel":         "Warn",
			"zapLogLevel":      "warn",
		},
		{
			"logMessage":       "this is a error log",
			"formatLogMessage": "this is a error log: %s",
			"logLevel":         "Error",
			"zapLogLevel":      "error",
		},
		{
			"logMessage":       "this is a fatal log",
			"formatLogMessage": "this is a fatal log: %s",
			"logLevel":         "Fatal",
			"zapLogLevel":      "fatal",
		},
	}

	testHertzLogger := reflect.ValueOf(logger)

	for _, v := range logTestSlice {
		t.Run(v["logLevel"], func(t *testing.T) {
			if v["logLevel"] == "Fatal" {
				defer func() {
					assert.Equal(t, "this is a fatal log", recover())
				}()
			}
			logFunc := testHertzLogger.MethodByName(v["logLevel"])
			logFunc.Call([]reflect.Value{
				reflect.ValueOf(v["logMessage"]),
			})
			assert.Contains(t, buf.String(), v["logMessage"])
			assert.Contains(t, buf.String(), v["zapLogLevel"])

			buf.Reset()

			logfFunc := testHertzLogger.MethodByName(fmt.Sprintf("%sf", v["logLevel"]))
			logfFunc.Call([]reflect.Value{
				reflect.ValueOf(v["formatLogMessage"]),
				reflect.ValueOf(v["logLevel"]),
			})
			assert.Contains(t, buf.String(), fmt.Sprintf(v["formatLogMessage"], v["logLevel"]))
			assert.Contains(t, buf.String(), v["zapLogLevel"])

			buf.Reset()

			ctx := context.Background()
			ctxLogfFunc := testHertzLogger.MethodByName(fmt.Sprintf("Ctx%sf", v["logLevel"]))
			ctxLogfFunc.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(v["formatLogMessage"]),
				reflect.ValueOf(v["logLevel"]),
			})
			assert.Contains(t, buf.String(), fmt.Sprintf(v["formatLogMessage"], v["logLevel"]))
			assert.Contains(t, buf.String(), v["zapLogLevel"])

			buf.Reset()
		})
	}
}

// TestLogLevel test SetLevel
func TestLogLevel(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger()
	defer logger.Sync()

	// output to buffer
	logging.SetLogger(logger)
	logging.SetOutput(buf)

	logging.Debug("this is a debug log")
	assert.False(t, strings.Contains(buf.String(), "this is a debug log"))

	logger.SetLevel(logging.LevelDebug)

	logging.Debugf("this is a debug log %s", "msg")
	assert.True(t, strings.Contains(buf.String(), "this is a debug log"))

	logger.SetLevel(logging.LevelError)
	logging.Infof("this is a debug log %s", "msg")
	assert.False(t, strings.Contains(buf.String(), "this is a info log"))

	logging.Warnf("this is a warn log %s", "msg")
	assert.False(t, strings.Contains(buf.String(), "this is a warn log"))

	logging.Errorf("this is a error log %s", "msg")
	assert.True(t, strings.Contains(buf.String(), "this is a error log"))
}

func TestWithCoreEnc(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(WithCoreEnc(zapcore.NewConsoleEncoder(humanEncoderConfig())))
	defer logger.Sync()

	// output to buffer
	logging.SetLogger(logger)
	logging.SetOutput(buf)

	logging.Infof("this is a info log %s", "msg")
	assert.True(t, strings.Contains(buf.String(), "this is a info log"))
}

func TestWithCoreWs(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(WithCoreWs(zapcore.AddSync(buf)))
	defer logger.Sync()

	logging.SetLogger(logger)
	logging.SetOutput(buf)

	logging.Infof("this is a info log %s", "msg")
	assert.True(t, strings.Contains(buf.String(), "this is a info log"))
}

func TestWithCoreLevel(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(WithCoreLevel(zap.NewAtomicLevelAt(zapcore.WarnLevel)))
	defer logger.Sync()

	// output to buffer
	logging.SetLogger(logger)
	logging.SetOutput(buf)

	logging.Infof("this is a info log %s", "msg")
	assert.False(t, strings.Contains(buf.String(), "this is a info log"))

	logging.Warnf("this is a warn log %s", "msg")
	assert.True(t, strings.Contains(buf.String(), "this is a warn log"))
}

// TestCoreOption test zapcore config option
func TestCoreOption(t *testing.T) {
	buf := new(bytes.Buffer)

	dynamicLevel := zap.NewAtomicLevel()

	dynamicLevel.SetLevel(zap.InfoLevel)

	logger := NewLogger(
		WithCores([]CoreConfig{
			{
				Enc: zapcore.NewConsoleEncoder(humanEncoderConfig()),
				Ws:  zapcore.AddSync(os.Stdout),
				Lvl: dynamicLevel,
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("./all/log.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("./debug/log.log"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.DebugLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("./info/log.log"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.InfoLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("./warn/log.log"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zap.WarnLevel
				}),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer("./error/log.log"),
				Lvl: zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev >= zap.ErrorLevel
				}),
			},
		}...),
	)
	defer logger.Sync()

	logging.SetLogger(logger)
	logging.SetOutput(buf)

	logging.Debug("this is a debug log")
	// test log level
	assert.False(t, strings.Contains(buf.String(), "this is a debug log"))

	logging.Error("this is a warn log")
	// test log level
	assert.True(t, strings.Contains(buf.String(), "this is a warn log"))
	// test console encoder result
	assert.True(t, strings.Contains(buf.String(), "\tERROR\t"))

	logger.SetLevel(logging.LevelDebug)
	logging.Debug("this is a debug log")
	assert.True(t, strings.Contains(buf.String(), "this is a debug log"))
}

// TestCoreOption test zapcore config option
func TestZapOption(t *testing.T) {
	buf := new(bytes.Buffer)

	logger := NewLogger(
		WithZapOptions(zap.AddCaller()),
	)
	defer logger.Sync()

	logging.SetLogger(logger)
	logging.SetOutput(buf)

	logging.Debug("this is a debug log")
	assert.False(t, strings.Contains(buf.String(), "this is a debug log"))

	logging.Error("this is a warn log")
	// test caller in log result
	assert.True(t, strings.Contains(buf.String(), "caller"))
}

// TestWithExtraKeys test WithExtraKeys option
func TestWithExtraKeys(t *testing.T) {
	buf := new(bytes.Buffer)

	log := NewLogger(WithExtraKeys([]ExtraKey{"requestId"}))

	logging.SetLogger(log)
	logging.SetOutput(buf)

	ctx := context.WithValue(context.Background(), ExtraKey("requestId"), "123")

	logging.CtxInfof(ctx, "%s log", "extra")

	var logStructMap map[string]interface{}

	err := json.Unmarshal(buf.Bytes(), &logStructMap)

	assert.Nil(t, err)

	value, ok := logStructMap["requestId"]

	assert.True(t, ok)
	assert.Equal(t, value, "123")
}

func TestPutExtraKeys(t *testing.T) {
	logger := NewLogger(WithExtraKeys([]ExtraKey{"abc"}))

	assert.Contains(t, logger.GetExtraKeys(), ExtraKey("abc"))
	assert.NotContains(t, logger.GetExtraKeys(), ExtraKey("def"))

	logger.PutExtraKeys("def")
	assert.Contains(t, logger.GetExtraKeys(), ExtraKey("def"))
}

func TestExtraKeyAsStr(t *testing.T) {
	buf := new(bytes.Buffer)
	const v = "value"

	logger := NewLogger(WithExtraKeys([]ExtraKey{"abc"}))

	logging.SetLogger(logger)
	logging.SetOutput(buf)

	ctx1 := context.TODO()
	ctx1 = context.WithValue(ctx1, "key1", v) //nolint:staticcheck
	logging.CtxErrorf(ctx1, "%s", "error")

	assert.NotContains(t, buf.String(), v)

	buf.Reset()

	strLogger := NewLogger(WithExtraKeys([]ExtraKey{"abc"}), WithExtraKeyAsStr())

	logging.SetLogger(strLogger)
	logging.SetOutput(buf)

	ctx2 := context.TODO()
	ctx2 = context.WithValue(ctx2, "key2", v) //nolint:staticcheck

	logging.CtxErrorf(ctx2, "%s", v)

	assert.Contains(t, buf.String(), v)

	buf.Reset()
}

func BenchmarkNormal(b *testing.B) {
	buf := new(bytes.Buffer)
	log := NewLogger()

	logging.SetLogger(log)
	logging.SetOutput(buf)
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		logging.CtxInfof(ctx, "normal log")
	}
}

func BenchmarkWithExtraKeys(b *testing.B) {
	buf := new(bytes.Buffer)
	log := NewLogger(WithExtraKeys([]ExtraKey{"requestId"}))
	logging.SetLogger(log)
	logging.SetOutput(buf)

	ctx := context.WithValue(context.Background(), ExtraKey("requestId"), "123")
	for i := 0; i < b.N; i++ {
		logging.CtxInfof(ctx, "normal log")
	}
}
