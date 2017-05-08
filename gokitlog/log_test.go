package gokitlog

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	levels "github.com/go-kit/kit/log/deprecated_levels"
	"github.com/go-kit/kit/log/level"

	"github.com/laohanlinux/assert"
)

func TestFirst(t *testing.T) {
	tmpLog := log.With(log.NewJSONLogger(os.Stdout), "caller", CallerNum, "level", level.InfoValue())
	tmpLog.Log("hello", "good")
}

func TestNopLogger(t *testing.T) {

	Error("time", time.Now().Unix())
	Debug("time", time.Now().Unix())
}

func TestGoKitLogger(t *testing.T) {
	logDir, err := os.Getwd()
	assert.Nil(t, err)
	opt := LogOption{
		LogDir:                logDir,
		SegmentationThreshold: 1,
		LogName:               "test",
		LogLevel:              "debug",
	}
	os.Remove(logDir + "/" + "test.log")
	runtime.GOMAXPROCS(runtime.NumCPU())
	lg, err := NewGoKitLogger(opt, FmtFormat)
	assert.Nil(t, err)
	assert.NotNil(t, lg)
	lg.Close()

	SetGlobalLog(opt)
	for i := 0; i < 1024; i++ {
		Info(i, i)
		Infof("i:%v", i)
	}
	Infof("%d", 1023)
	Debug("test", time.Now().Unix())
	Debug("test", time.Now().Unix())
	Debug("test", time.Now().Unix())
	Debugf("Hello Word, my name is lusy")
	Debug("array", []string{"a", "b"})
	GlobalLog().Log("log", "log")
	Crit("exit", true)
}

func TestGoKitLogWriter(t *testing.T) {
	logDir, err := os.Getwd()
	assert.Nil(t, err)
	opt := LogOption{
		LogDir:                logDir,
		SegmentationThreshold: 1,
		LogName:               "test",
		LogLevel:              "warn",
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	levelSets := strings.Split(opt.LogLevel, "|")
	lg, err := NewLogWriter(opt)
	defer lg.Close()
	logger := level.NewFilter(log.NewJSONLogger(lg), WrapLogLevel(levelSets)...)

	logger = log.With(logger, "caller", log.Caller(5))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// swap logger saftly
	logger = log.NewSyncLogger(logger)

	for i := 0; i < 10; i++ {
		go func(gid int) {
			levelLog := levels.New(logger)
			tmpLog := levelLog.Warn()
			if gid%2 == 0 {
				tmpLog = levelLog.Info()
			}
			tmpLog = log.With(tmpLog, "gorutine", gid)
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
