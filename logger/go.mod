module github.com/pemako/gopkg/logger

go 1.22.1

require (
	github.com/pemako/gopkg/lumberjack v0.1.1
	github.com/pemako/gopkg/rotatelogs v0.1.1
	go.uber.org/zap v1.27.0
)


replace (
  github.com/pemako/gopkg/lumberjack => ../lumberjack
  github.com/pemako/gopkg/rotatelogs => ../rotatelogs
  github.com/pemako/gopkg/strftime => ../strftime
)
