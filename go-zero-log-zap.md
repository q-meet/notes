# Zap 日志库实践

本文详细介绍了非常流行的 Uber 开源的 zap 日志库，同时介绍了如何搭配 Lumberjack、Rotatelogs 实现日志的切割和归档。

## 1、日志需求

我们重温一下一个好的日志记录器都需要能够提供下面哪些功能？

- 良好日志写入性能
- 支持不同的日志级别。并且可分离成多个日志文件
- 多输出 － 同时支持标准输出，文件等
- 能够打印基本信息，如调用文件 / 函数名和行号，日志时间等
- 可读性与结构化，Json格式或有分隔符，方便后续的日志采集、监控等
- 文件切割，可按小时、天进行日志拆分，或者按文件大小
- 日志书写友好，支持通过context自动log trace等
- 文件定时删除
- 开源性，与其他开源框架支持较好

## 2、Uber-go Zap

Zap 是非常快的、结构化的，分日志级别的 Go 日志库。

### 2.1 为什么选择 Uber-go zap

有关技术选型可以见：[www.yuque.com/jinsesihuan…](https://www.yuque.com/jinsesihuanian/gpwou5/aduc5c)
​

### 2.2 安装

运行下面的命令安装 zap

```bash
go get -u go.uber.org/zap
```

### 2.3 配置 Zap Logger

Zap 提供了两种类型的日志记录器—Sugared Logger和Logger。
在性能很好但不是很关键的上下文中，使用SugaredLogger。它比其他结构化日志记录包快 4-10 倍，并且支持结构化和 printf 风格的日志记录。

在每一微秒和每一次内存分配都很重要的上下文中，使用Logger。它甚至比SugaredLogger更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。

#### 2.3.1 Logger

通过调用zap.NewProduction()/zap.NewDevelopment()或者zap.Example()创建一个 Logger。
上面的每一个函数都将创建一个 logger。唯一的区别在于它将记录的信息不同。

例如 production logger 默认记录调用函数信息、日期和时间等。

通过 Logger 调用 Info/Error 等。
默认情况下日志都会打印到应用程序的 console 界面。

```go
var logger *zap.Logger

func main() {
	InitLogger()
  defer logger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}

func InitLogger() {
	logger, _ = zap.NewProduction()
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		logger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}
```

在上面的代码中，我们首先创建了一个 Logger，然后使用 Info/ Error 等 Logger 方法记录消息。

日志记录器方法的语法是这样的：

```go
func (log *Logger) MethodXXX(msg string, fields ...Field)
```

其中MethodXXX是一个可变参数函数，可以是 Info / Error/ Debug / Panic 等。每个方法都接受一个消息字符串和任意数量的zapcore.Field场参数。

每个zapcore.Field其实就是一组键值对参数。

我们执行上面的代码会得到如下输出结果：

```json
{"level":"error","ts":1572159218.912792,"caller":"zap_demo/temp.go:25","msg":"Error fetching url..","url":"www.baidu.com","error":"Get www.sogo.com: unsupported protocol scheme \"\"","stacktrace":"main.simpleHttpGet\n\t/Users/q1mi/zap_demo/temp.go:25\nmain.main\n\t/Users/q1mi/zap_demo/temp.go:14\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}
{"level":"info","ts":1572159219.1227388,"caller":"zap_demo/temp.go:30","msg":"Success..","statusCode":"200 OK","url":"http://www.baidu.com"}
```

#### 2.3.2 Sugared Logger

现在让我们使用 Sugared Logger 来实现相同的功能。

大部分的实现基本都相同。

惟一的区别是，我们通过调用主 logger 的. Sugar()方法来获取一个SugaredLogger。

然后使用SugaredLogger以printf格式记录语句

下面是修改过后使用SugaredLogger代替Logger的代码：

```go
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
}

func InitLogger() {
  logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}
```

当你执行上面的代码会得到如下输出：

```json
{"level":"error","ts":1572159149.923002,"caller":"logic/temp2.go:27","msg":"Error fetching URL www.baidu.com : Error = Get www.baidu.com: unsupported protocol scheme \"\"","stacktrace":"main.simpleHttpGet\n\t/Users/zap_demo/logic/temp2.go:27\nmain.main\n\t/Users/zap_demo/logic/temp2.go:14\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}
{"level":"info","ts":1572159150.192585,"caller":"logic/temp2.go:29","msg":"Success! statusCode = 200 OK for URL http://www.baidu.com"}
```

你应该注意到的了，到目前为止这两个 logger 都打印输出 JSON 结构格式。 在本博客的后面部分，我们将更详细地讨论 SugaredLogger，并了解如何进一步配置它。

### 2.4 定制 logger

#### 2.4.1 将日志写入文件而不是终端

我们要做的第一个更改是把日志写入文件，而不是打印到应用程序控制台。
​

我们将使用zap.New(…)方法来手动传递所有配置，而不是使用像zap.NewProduction()这样的预置方法来创建 logger。

```go
func New(core zapcore.Core, options ...Option) *Logger
```

zapcore.Core需要三个配置——Encoder，WriteSyncer，LogLevel。

1.Encoder: 编码器 (写入日志格式)。我们将使用开箱即用的NewJSONEncoder()，并使用预先设置的ProductionEncoderConfig()。

```go
zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
```

2.WriterSyncer ：指定日志写到哪里去。我们使用zapcore.AddSync()函数并且将打开的文件句柄传进去。

```go
  file, _ := os.Create("./test.log")
   writeSyncer := zapcore.AddSync(file)
```

3.Log Level：哪种级别的日志将被写入。
我们将修改上述部分中的 Logger 代码，并重写InitLogger()方法。其余方法—main() /SimpleHttpGet()保持不变。

```go
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}
```

当使用这些修改过的 logger 配置调用上述部分的main()函数时，以下输出将打印在文件——test.log中。

```json
{"level":"debug","ts":1572160754.994731,"msg":"Trying to hit GET request for www.baidu.com"}
{"level":"error","ts":1572160754.994982,"msg":"Error fetching URL www.sogo.com : Error = Get www.baidu.com: unsupported protocol scheme \"\""}
{"level":"debug","ts":1572160754.994996,"msg":"Trying to hit GET request for http://www.baidu.com"}
{"level":"info","ts":1572160757.3755069,"msg":"Success! statusCode = 200 OK for URL http://www.baidu.com"}
```

## 3、使用 Lumberjack 进行日志切割归档

_Zap 本身不支持切割归档日志文件，_为了添加日志切割归档功能，我们将使用第三方库 Lumberjack 来实现。

### 3.1 安装

执行下面的命令安装 Lumberjack

```bash
go get -u github.com/natefinch/lumberjack
```

### 3.2 zap logger中加入Lumberjack

要在 zap 中加入 Lumberjack 支持，我们需要修改WriteSyncer代码。我们将按照下面的代码修改getLogWriter()函数：

```go
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
```

Lumberjack Logger 采用以下属性作为输入:

Filename: 日志文件的位置
MaxSize：在进行切割之前，日志文件的最大大小（以 MB 为单位）
MaxBackups：保留旧文件的最大个数
MaxAges：保留旧文件的最大天数
Compress：是否压缩 / 归档旧文件

## 4、使用rotatelogs实现日志文件处理

### 4.1 当前问题

我们当前还存在着两个问题：

生成的日志并未包含日期后缀，Lumberjack在备份时会包含  

所有级别的日志都写在一个文件中了，需要根据日志级别写入不同文件

这里我写了一个基于zap + "github.com/lestrrat/go-file-rotatelogs" 工具的工具类，完成自动追加日期后缀，并将不同日志级别的日志写入不同的文件。效果如下

![Alt text](images/zap_tool.png)

```go
package zaplog

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Options struct {
	LogFileDir    string	//日志路径
	AppName       string	// Filename是要写入日志的文件前缀
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	MaxSize       int	// 一个文件多少Ｍ大于该数字开始切分文件
	MaxBackups    int	// MaxBackups是要保留的最大旧日志文件数
	MaxAge        int	// MaxAge是根据日期保留旧日志文件的最大天数
	zap.Config
}

var (
	logger                         *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

func init() {
	logger = &Logger{
		Opts: &Options{},
	}
}

type Logger struct {
	*zap.SugaredLogger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func initLogger(cf ...*Options) {
	logger.Lock()
	defer logger.Unlock()
	if logger.inited {
		logger.Info("[initLogger] logger Inited")
		return
	}
	if len(cf) > 0 {
		logger.Opts = cf[0]
	}
	logger.loadCfg()
	logger.init()
	logger.Info("[initLogger] zap plugin initializing completed")
	logger.inited = true
}


// GetLogger returns logger
func GetLogger() (ret *Logger) {
	return logger
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	mylogger, err := l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	l.SugaredLogger = mylogger.Sugar()
	defer l.SugaredLogger.Sync()
}

func (l *Logger) loadCfg() {
	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}
	// 默认输出到程序运行目录的logs子目录
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}
	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}
	if l.Opts.ErrorFileName == "" {
		l.Opts.ErrorFileName = "error.log"
	}
	if l.Opts.WarnFileName == "" {
		l.Opts.WarnFileName = "warn.log"
	}
	if l.Opts.InfoFileName == "" {
		l.Opts.InfoFileName = "info.log"
	}
	if l.Opts.DebugFileName == "" {
		l.Opts.DebugFileName = "debug.log"
	}
	if l.Opts.MaxSize == 0 {
		l.Opts.MaxSize = 100
	}
	if l.Opts.MaxBackups == 0 {
		l.Opts.MaxBackups = 30
	}
	if l.Opts.MaxAge == 0 {
		l.Opts.MaxAge = 30
	}
}

func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		//return zapcore.AddSync(&lumberjack.Logger{
		//	Filename:   logger.Opts.LogFileDir + sp + logger.Opts.AppName + "-" + fN,
		//	MaxSize:    logger.Opts.MaxSize,
		//	MaxBackups: logger.Opts.MaxBackups,
		//	MaxAge:     logger.Opts.MaxAge,
		//	Compress:   true,
		//	LocalTime:  true,
		//})
		// 每小时一个文件
		logf, _ := rotatelogs.New(l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN +".%Y_%m%d_%H",
			rotatelogs.WithLinkName(l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN),
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(time.Minute),
		)
		return zapcore.AddSync(logf)
	}
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	return
}

func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	//consoleEncoder := zapcore.NewConsoleEncoder(logger.zapConfig.EncoderConfig)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),
	}
	if l.Opts.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}
```

### 4.3 测试

```go
package zaplog

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestZapLog(t *testing.T) {
	data := &Options{
		LogFileDir: "/Users/didi/code/golang/learn-go/zap_log/v3/logs",
		AppName:    "logtool",
		MaxSize:    30,
		MaxBackups: 7,
		MaxAge:     7,
		Config:     zap.Config{},
	}
	data.Development = true
	initLogger(data)
	for i := 0; i < 2; i++ {
		time.Sleep(time.Second)
		logger.Debug(fmt.Sprint("debug log ", i), zap.Int("line", 999))
		logger.Info(fmt.Sprint("Info log ", i), zap.Any("level", "1231231231"))
		logger.Warn(fmt.Sprint("warn log ", i), zap.String("level", `{"a":"4","b":"5"}`))
		logger.Error(fmt.Sprint("err log ", i), zap.String("level", `{"a":"7","b":"8"}`))
	}
}
```

## 我的参考实现 go-zero

```go
package zapx

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"time"
)

var (
	//logger                         *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer // IO输出
	//debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	//errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type ZapWriter struct {
	logger *zap.Logger
}

func Level() zapcore.Level {
	return zapcore.DebugLevel
}
func InitLogger() (logx.Writer, error) {
	getLogWriter()
	encoder := getEncoder()

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-Level() > -1
	})

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, errWS, errPriority),
		zapcore.NewCore(encoder, warnWS, warnPriority),
		zapcore.NewCore(encoder, infoWS, infoPriority),
		zapcore.NewCore(encoder, debugWS, debugPriority),
	)

	logger := zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑

	return &ZapWriter{
		logger: logger,
	}, nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.StacktraceKey = "Stacktrace"
	//encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString(t.Format("2006-01-02 15:04:05"))
	//}
	return zapcore.NewJSONEncoder(encoderConfig)
	//return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() {
	f := func(fN string) zapcore.WriteSyncer {
		logf, _ := rotatelogs.New(
			"logs"+sp+fN+".%Y_%m%d.log",
			//"logs"+sp+fN+".%Y_%m%d_%H.log",
			//rotatelogs.WithLinkName("logs"+sp+fN),
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(time.Minute),
		)
		return zapcore.AddSync(logf)
	}
	var ErrorFileName = "Error"
	var WarnFileName = "Warn"
	var InfoFileName = "Info"
	var DebugFileName = "Debug"
	errWS = f(ErrorFileName)
	warnWS = f(WarnFileName)
	infoWS = f(InfoFileName)
	debugWS = f(DebugFileName)
	/*
		Filename: 日志文件的位置
		MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups：保留旧文件的最大个数
		MaxAges：保留旧文件的最大天数
		Compress：是否压缩/归档旧文件
	*/
	/*
		   lumberJackLogger := &lumberjack.Logger{
				   Filename:   "./test.log",
				   MaxSize:    10,
				   MaxBackups: 5,
				   MaxAge:     30,
				   Compress:   false,
				}
			  return zapcore.AddSync(lumberJackLogger)
	*/
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *ZapWriter) Close() error {
	return w.logger.Sync()
}

func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toZapDeFields(fields...)...)
}

func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapDeFields(fields...)...)
}

func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapDeFields(fields...)...)
}

func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapDeFields(fields...)...)
}

func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapDeFields(fields...)...)
}

func toZapDeFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
```

使用 在启动服务之前赋值

```go
	writer, err := zapx.InitLogger()
	logx.Must(err)
	logx.SetWriter(writer)
```

参考地址：[github](https://github.com/q-meet/go-zero-mall/blob/main/rpc-common/log/zapx/demerger.go)
