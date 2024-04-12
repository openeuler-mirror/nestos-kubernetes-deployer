/*
Copyright 2024 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	logFilePrefix = "logs/nkd-"
	logFileSuffix = ".log"
	maxLogSize    = 10 // 每个日志文件的最大大小（MB）
	maxBackups    = 10 // 保留旧日志文件的最大个数
	maxAgeInDays  = 30 // 保留旧日志文件的最大天数
)

var LogLevel string

type loggerHook struct {
	file      io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func NewloggerHook(file io.Writer, level logrus.Level, formatter logrus.Formatter) *loggerHook {
	return &loggerHook{
		file:      file,
		level:     level,
		formatter: formatter,
	}
}

// Levels 返回允许记录的日志级别列表
func (h *loggerHook) Levels() []logrus.Level {
	var levels []logrus.Level
	for _, level := range logrus.AllLevels {
		if level <= h.level {
			levels = append(levels, level)
		}
	}

	return levels
}

// Fire 实现钩子的 Fire 方法，用于记录日志
func (h *loggerHook) Fire(entry *logrus.Entry) error {
	orig := entry.Message
	defer func() { entry.Message = orig }()

	msgs := strings.Split(orig, "\n")
	for _, msg := range msgs {
		entry.Message = msg
		line, err := h.formatter.Format(entry)
		if err != nil {
			return err
		}
		if _, err := h.file.Write(line); err != nil {
			return err
		}
	}

	return nil
}

// 设置日志文件的基本配置，包括创建日志目录，打开日志文件、设置日志格式等
func SetuploggerHook(baseDir string) func() {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create base directory for logs"))
	}

	logfilePath := filepath.Join(baseDir, generateLogFileName())
	logfile := &lumberjack.Logger{
		Filename:   logfilePath,
		MaxSize:    maxLogSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeInDays,
		Compress:   true,
		LocalTime:  true,
	}

	orgHooks := logrus.LevelHooks{}
	for k, v := range logrus.StandardLogger().Hooks {
		orgHooks[k] = v
	}
	logrus.AddHook(NewloggerHook(logfile, logrus.TraceLevel, &logrus.TextFormatter{
		DisableColors:          true,
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: false,
	}))
	logrus.SetLevel(logrus.TraceLevel)

	return func() {
		logfile.Close()
		logrus.StandardLogger().ReplaceHooks(orgHooks)
	}
}

func generateLogFileName() string {
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")
	return fmt.Sprintf("%s%s%s", logFilePrefix, dateString, logFileSuffix)
}
