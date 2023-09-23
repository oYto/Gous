package config

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

// tabFormatter tab 数据格式化
type logFormatter struct {
	log.TextFormatter
}

func (c *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	// 将调用栈中的文件名和行号格式化为更可读的形式，以便在日志或错误信息中使用。
	prettyCaller := func(frame *runtime.Frame) string { // 传入调用栈的信息
		_, fileName := filepath.Split(frame.File)         // 分割路径
		return fmt.Sprintf("%s:%d", fileName, frame.Line) // 返回调用栈文件名和行号
	}

	var b *bytes.Buffer
	// 缓存日志消息
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	// 将日志时间和日志级别添加到缓冲区 b 中，使用 fmt.Sprintf 进行格式化。
	b.WriteString(fmt.Sprintf("[%s] %s", entry.Time.Format(c.TimestampFormat), // 输出日志时间
		strings.ToUpper(entry.Level.String())))
	// entry.HasCaller() 检查日志条目是否包含调用者信息。
	if entry.HasCaller() {
		b.WriteString(fmt.Sprintf("[%s]", prettyCaller(entry.Caller))) // 输出日志所在文件，行数位置
	}
	//  b.WriteString 将日志内容添加到缓冲区 b 中，使用 fmt.Sprintf 进行格式化。
	b.WriteString(fmt.Sprintf(" %s\n", entry.Message)) // 输出日志内容
	// b.Bytes() 将缓冲区 b 转换为字节数组返回
	return b.Bytes(), nil
}

func setGinLog(out io.Writer) {
	gin.DefaultWriter = out
	gin.DefaultErrorWriter = out
}
