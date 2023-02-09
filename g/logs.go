package g

import (
	"encoding/json"
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

//var SvcContext *svc.ServiceContext

var logger *zap.Logger
var Log *zap.Logger
var Logs *zap.SugaredLogger
var l *zap.SugaredLogger

// Any 任意类型
type Any = interface{}

type ZapLogConf struct {
	Level       string `json:",default=info"`
	LogPath     string `json:",default=./logs/mylog.log"`
	FileEncoder string `json:",default=json"`
}

// InitLog 传入 nil, 则使用缺省的配置
func InitLog(conf *ZapLogConf) {
	conf = InitDefaultConfig(conf)
	var logPath string
	var logLevel zapcore.Level = 0
	logPath = conf.LogPath
	logLevel.UnmarshalText([]byte(conf.Level))
	configConsole := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		CallerKey:     "file",
		StacktraceKey: "stack",
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）带颜色
		EncodeCaller:  zapcore.ShortCallerEncoder,       //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}

	configFile := zapcore.EncoderConfig{
		MessageKey:   "msg",                       //结构化（json）输出：msg的key
		LevelKey:     "level",                     //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      "ts",                        //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    "file",                      //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		StacktraceKey: "stack",
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	//自定义日志级别：自定义Info级别
	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.WarnLevel && lvl >= logLevel
	//})

	//自定义日志级别：自定义Warn级别
	//warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.WarnLevel && lvl >= logLevel
	//})

	// 获取io.Writer的实现
	infoWriter := getWriter(logPath)
	//warnWriter := getWriter(errPath)

	// 实现多个输出
	// NewConsoleEncoder 是非结构化输出
	// NewJSONEncoder 是结构化输出

	var fileEncoder zapcore.Encoder
	if conf.FileEncoder == "console" {
		fileEncoder = zapcore.NewConsoleEncoder(configFile)
	} else {
		fileEncoder = zapcore.NewJSONEncoder(configFile)
	}
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(infoWriter), logLevel),
		//zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel), //将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
		//zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(warnWriter), warnLevel),//warn及以上写入errPath
		zapcore.NewCore(zapcore.NewConsoleEncoder(configConsole), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), logLevel), //同时将日志输出到控制台，
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	Log = logger
	Logs = logger.Sugar()
	l = logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	Logs.Debugf("zap log config init :\n%s", JsonPrettify(conf))
}
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 filename.YYmmddHH

	var hook io.Writer
	var err error
	if runtime.GOOS == "windows" {
		hook, err = rotatelogs.New(
			filename+".%Y%m%d%H",
			// windows 下无法创建软链接
			//rotatelogs.WithLinkName(filename),
			rotatelogs.WithMaxAge(time.Hour*24*30),
			rotatelogs.WithRotationTime(time.Hour*24),
		)
	} else {
		// filename是指向最新日志的链接
		hook, err = rotatelogs.New(
			filename+".%Y%m%d%H",
			rotatelogs.WithLinkName(filename),
			rotatelogs.WithMaxAge(time.Hour*24*30),
			rotatelogs.WithRotationTime(time.Hour*24),
		)
	}

	if err != nil {
		log.Println("日志启动异常")
		panic(err)
	}
	return hook
}

func InitDefaultConfig(conf *ZapLogConf) *ZapLogConf {
	if conf == nil {
		conf = &ZapLogConf{}
	}
	if conf.LogPath == "" {
		conf.LogPath = "./logs/mylog.log"
	}
	if conf.Level == "" {
		conf.Level = "info"
	}
	return conf
}

// =====

func Debug(args ...interface{}) {
	l.Info(args...)
}

func Info(args ...interface{}) {
	l.Info(args...)
}

func Warn(args ...interface{}) {
	l.Warn(args...)
}

func Error(args ...interface{}) {
	l.Error(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	l.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	l.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	l.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	l.Errorf(template, args...)
}

func JsonPrettify(m Any) string {
	resp, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		Logs.Errorf("转换json字符串失败，%+v", m)
	}
	return string(resp)
}
