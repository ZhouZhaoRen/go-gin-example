package logging
import (
	"fmt"
	"go-gin-example/pkg/file"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix = ""
	DefaultCallerDepth = 2

	logger *log.Logger
	logPrefix = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota // 从这里一直加1
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	var err error
	filePath := getLogFilePath() // runtime/logs
	fileName:=getLogFileName()  // log20201123.log
	//F = openLogFile(filePath)
	F,err=file.MustOpen(fileName,filePath)
	if err!=nil{
		log.Fatalf("logging.Setup err: %v",err)
	}
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
// 日志分为两部分，第一部分是日志前缀，比如说产生日志的位置，时间等，第二部分是我们想记录的日志的具体消息
// 给log日志添加前缀  [INFO][auth.go:43]2020/10/26 21:16:34 [请求参数错误]
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth) // 获取到文件信息和所在的行
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
