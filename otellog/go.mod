module github.com/kanhai-syd/hailog/telemetry

go 1.22

toolchain go1.22.2

replace github.com/kanhai-syd/hailog/logging/logrus => ../log/logging/logrus

replace github.com/kanhai-syd/hailog/logging/slog => ../log/logging/slog

replace github.com/kanhai-syd/hailog/logging/zap => ../log/logging/zap

replace github.com/kanhai-syd/hailog/logging/zerolog => ../log/logging/zerolog

require (
	github.com/kanhai-syd/hailog v0.0.0-20241109102407-978409ba316a
	github.com/kanhai-syd/hailog/logging/logrus v0.0.0-00010101000000-000000000000
	github.com/kanhai-syd/hailog/logging/slog v0.0.0-00010101000000-000000000000
	github.com/kanhai-syd/hailog/logging/zap v0.0.0-00010101000000-000000000000
	github.com/kanhai-syd/hailog/logging/zerolog v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.33.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel v1.32.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.32.0
	go.opentelemetry.io/otel/sdk v1.32.0
	go.opentelemetry.io/otel/trace v1.32.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
