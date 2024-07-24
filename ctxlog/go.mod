module github.com/pemako/gopkg/ctxlog

go 1.22.1

require (
	github.com/pemako/gopkg/lumberjack v0.1.2
	github.com/pemako/gopkg/rotatelogs v0.1.2
	go.uber.org/zap v1.27.0
)

require (
	github.com/pemako/gopkg/strftime v0.1.2 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace (
	github.com/pemako/gopkg/lumberjack => ../lumberjack
	github.com/pemako/gopkg/rotatelogs => ../rotatelogs
	github.com/pemako/gopkg/strftime => ../strftime
)
