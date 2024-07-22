module github.com/pemako/gopkg/logger

go 1.22.1

require (
	github.com/pemako/gopkg/lumberjack v0.1.0
	github.com/pemako/gopkg/rotatelogs v0.1.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/pemako/gopkg/strftime v0.1.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace github.com/pemako/gopkg/lumberjack => ../lumberjack

replace github.com/pemako/gopkg/rotatelogs => ../rotatelogs

replace github.com/pemako/gopkg/strftime => ../strftime
