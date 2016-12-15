package gokitlog

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	level "github.com/go-kit/kit/log/experimental_level"
	"github.com/go-kit/kit/log/levels"
	"github.com/laohanlinux/assert"
)

func TestGoKitLogger(t *testing.T) {
	logDir, err := os.Getwd()
	assert.Nil(t, err)
	opt := LogOption{
		LogDir:                logDir,
		SegmentationThreshold: 1,
		LogName:               "test",
		LogLevel:              "info|warn|error|debug|crit",
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	lg, err := NewGoKitLogger(opt)
	assert.Nil(t, err)
	assert.NotNil(t, lg)
	lg.Crit().Log("testElement", time.Now().String())
	lg.Close()

	SetGlobalLog(opt)
	for i := 0; i < 1024; i++ {
		Info(i, i)
		Infof("i:%d", i)
	}
	Info()
	Infof()
	Infof("%d")
	Infof("%d", "dddd")
	Infof("%d", 1023)
	Crit("exit", true)
}

func TestGoKitLogWriter(t *testing.T) {
	logDir, err := os.Getwd()
	assert.Nil(t, err)
	opt := LogOption{
		LogDir:                logDir,
		SegmentationThreshold: 1,
		LogName:               "test",
		LogLevel:              "warn|error",
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	levelSets := strings.Split(opt.LogLevel, "|")
	lg, err := NewLogWriter(opt)
	defer lg.Close()
	logger := level.New(log.NewJSONLogger(lg), level.Config{Allowed: levelSets})
	logger = log.NewContext(logger).With("caller", log.Caller(5))
	logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)

	// swap logger saftly
	logger = log.NewSyncLogger(logger)

	for i := 0; i < 10; i++ {
		go func(gid int) {
			levelLog := levels.New(logger)
			tmpLog := levelLog.Warn()
			if gid%2 == 0 {
				tmpLog = levelLog.Info()
			}
			tmpLog = log.NewContext(tmpLog).With("gorutine", gid)
			for {
				tmpLog.Log("msg", time.Now().Unix())
				//time.Sleep(time.Second)
			}
		}(i)
	}

	time.Sleep(time.Second * 170)

}

func TestLogWriter(t *testing.T) {
	logDir, err := os.Getwd()
	assert.Nil(t, err)
	opt := LogOption{
		LogDir:                logDir,
		SegmentationThreshold: 1,
		LogName:               "test",
	}

	lg, err := NewLogWriter(opt)
	defer lg.Close()
	assert.Nil(t, err)
	for i := 0; i < 62; i++ {
		lg.Write([]byte(strconv.Itoa(i)))
		time.Sleep(time.Second)
	}

}
