module github.com/kanhai-syd/hailog/logging/logrus

go 1.21

replace github.com/kanhai-syd/hailog => ../..

require (
	github.com/kanhai-syd/hailog v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
