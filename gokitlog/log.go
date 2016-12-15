package gokitlog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	level "github.com/go-kit/kit/log/experimental_level"
	"github.com/go-kit/kit/log/levels"
)

const (
	logNameFormat = "2006-01-02_15:04"
	CallerNum     = 5
)

func init() {
	tmpLog := log.NewJSONLogger(os.Stdout)
	tmpLog = log.NewContext(tmpLog).With("caller", log.DefaultCaller)
	tmpLog = log.NewContext(tmpLog).With("ts", log.DefaultTimestampUTC)
	tmpLog = level.New(tmpLog, level.Config{Allowed: []string{"warn", "info", "debug", "error", "crit"}})
	tmpLog = log.NewSyncLogger(tmpLog)
	levelslog := levels.New(tmpLog)

	lg = &GoKitLogger{
		Logger:   tmpLog,
		Levels:   &levelslog,
		ioWriter: nil,
		sync:     true,
	}
}

func NewGoKitLogger(opt LogOption) (*GoKitLogger, error) {
	ioWriter, err := NewLogWriter(opt)
	if err != nil {
		return nil, err
	}
	levelsSets := strings.Split(opt.LogLevel, "|")
	tmpLog := log.NewJSONLogger(ioWriter)
	tmpLog = log.NewContext(tmpLog).With("caller", log.DefaultCaller)
	tmpLog = log.NewContext(tmpLog).With("ts", log.DefaultTimestampUTC)
	tmpLog = level.New(tmpLog, level.Config{Allowed: levelsSets})
	if opt.Sync {
		tmpLog = log.NewSyncLogger(tmpLog)
	}
	levellog := levels.New(tmpLog)
	return &GoKitLogger{Levels: &levellog, Logger: tmpLog}, nil
}

func GlobalLog() *GoKitLogger {
	return lg
}

func SetGlobalLog(opt LogOption) {
	//close old logger io writer
	Close()
	tmpLog, err := NewGoKitLogger(opt)
	if err != nil {
		panic(err)
	}
	lg = tmpLog
}

func SetGlobalLogWithLog(logger log.Logger, levelConf ...level.Config) {
	ioWriter := lg.ioWriter
	defer func() {
		if ioWriter != nil {
			ioWriter.Close()
		}
	}()
	lg.Logger = log.NewContext(logger).With("caller", log.Caller(CallerNum))
	var levelLog levels.Levels
	if len(levelConf) > 0 {
		levelLog = levels.New(level.New(lg.Logger, levelConf[0]), nil)
	} else {
		levelLog = levels.New(lg.Logger, nil)
	}
	lg.Levels = &levelLog
}

var lg *GoKitLogger

type GoKitLogger struct {
	log.Logger
	*levels.Levels
	ioWriter *LogWriter
	sync     bool
}

func (gklog *GoKitLogger) Close() error {
	if gklog.ioWriter != nil {
		return gklog.ioWriter.Close()
	}
	return nil
}

func Debug(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Debug()).With("caller",
		log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Debugf(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Debug()).With("caller",
		log.Caller(CallerNum))
	logPrintf(tmpLog, args)
}

func Info(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Info()).With("caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Infof(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Info()).With("caller", log.Caller(CallerNum))
	logPrintf(tmpLog, args)
}

func Warn(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Warn()).With("caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Warnf(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Warn()).With("caller", log.Caller(CallerNum))
	logPrintf(tmpLog, args)
}

func Error(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Error()).With("caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Errorf(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Error()).With("caller", log.Caller(CallerNum))
	logPrintf(tmpLog, args)
}

func Crit(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Crit()).With("caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
	os.Exit(1)
}

func Critf(args ...interface{}) {
	tmpLog := log.NewContext(lg.Levels.Crit()).With("caller", log.Caller(CallerNum))
	logPrintf(tmpLog, args)
	os.Exit(1)
}

func Log(args ...interface{}) {
	tmpLog := log.NewContext(lg.Logger).With("caller", log.Caller(CallerNum))
	logPrint(tmpLog, args)
}

func Close() error {
	if lg.ioWriter != nil {
		return lg.ioWriter.Close()
	}
	return nil
}

type LogOption struct {
	// unit in minutes
	SegmentationThreshold int    `toml:"threshold"`
	LogDir                string `toml:"log_dir"`
	LogName               string `toml:"log_name"`
	LogLevel              string `toml:"log_level"`
	Sync                  bool   `toml:"sync"`
}

type LogWriter struct {
	oldTime               time.Time
	segmentationThreshold float64
	logDir                string
	logName               string
	*os.File
}

func NewLogWriter(opt LogOption) (*LogWriter, error) {
	logWriter := &LogWriter{
		oldTime:               time.Now(),
		segmentationThreshold: float64(opt.SegmentationThreshold),
		logDir:                opt.LogDir,
		logName:               opt.LogName,
	}

	fp, err := os.OpenFile(fmt.Sprintf("%s/%s.log", opt.LogDir, opt.LogName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	logWriter.File = fp
	return logWriter, nil
}

// TODO
// use bufio buffer
func (lw *LogWriter) Write(p []byte) (n int, err error) {
	if time.Since(lw.oldTime).Minutes() > lw.segmentationThreshold {
		if err = lw.renameLogFile(); err != nil {
			return -1, err
		}

		lw.File, err = os.OpenFile(fmt.Sprintf("%s/%s.log", lw.logDir, lw.logName),
			os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return -1, err
		}

	}
	return lw.File.Write(p)
}

func (lw *LogWriter) Close() error {
	// return lw.renameLogFile()
	return lw.File.Close()
}

func (lw *LogWriter) renameLogFile() error {
	// split log file
	stat, err := lw.File.Stat()
	if err != nil {
		return err
	}
	srcFileName := fmt.Sprintf("%s/%s", lw.logDir, stat.Name())
	if err = lw.File.Close(); err != nil {
		return err
	}
	dstFileName := fmt.Sprintf("%s/%s_%s.log", lw.logDir, lw.logName,
		lw.oldTime.Format(logNameFormat))
	fmt.Println(dstFileName, srcFileName)
	os.Rename(srcFileName, dstFileName)
	lw.oldTime = time.Now()
	return nil
}

func logPrint(logger log.Logger, args interface{}) {
	tmpArgs, _ := args.([]interface{})
	if tmpArgs == nil || len(tmpArgs) == 0 {
		logger.Log()
		return
	}
	logger.Log(tmpArgs...)
}

func logPrintf(logger log.Logger, args interface{}) {
	tmpArgs, _ := args.([]interface{})
	if tmpArgs == nil || len(tmpArgs) == 0 {
		logger.Log("msg")
		return
	}
	if len(tmpArgs) == 1 {
		logger.Log("msg", fmt.Sprintf("%s", tmpArgs[0]))
		return
	}

	msgContent := fmt.Sprintf(fmt.Sprintf("%v", tmpArgs[0]), tmpArgs[1:]...)
	logger.Log("msg", msgContent)
}
