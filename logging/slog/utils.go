// Copyright 2023 CloudWeGo Authors.
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

package slog

import (
	"log/slog"

	"github.com/kanhai-syd/hailog/logging"
)

// Adapt log level to otelslog level
func tranSLevel(level logging.Level) (lvl slog.Level) {
	switch level {
	case logging.LevelTrace:
		lvl = LevelTrace
	case logging.LevelDebug:
		lvl = slog.LevelDebug
	case logging.LevelInfo:
		lvl = slog.LevelInfo
	case logging.LevelWarn:
		lvl = slog.LevelWarn
	case logging.LevelNotice:
		lvl = LevelNotice
	case logging.LevelError:
		lvl = slog.LevelError
	case logging.LevelFatal:
		lvl = LevelFatal
	default:
		lvl = slog.LevelWarn
	}
	return
}
