module github.com/pemako/gopkg

go 1.22.1

replace (
	github.com/pemako/gopkg/config => ./config
	github.com/pemako/gopkg/ctxlog => ./ctxlog
	github.com/pemako/gopkg/envload => ./envload
	github.com/pemako/gopkg/logger => ./logger
	github.com/pemako/gopkg/lumberjack => ./lumberjack
	github.com/pemako/gopkg/rotatelogs => ./rotatelogs
	github.com/pemako/gopkg/strftime => ./strftime
)
