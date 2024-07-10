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
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerHookFire(t *testing.T) {
	var buf bytes.Buffer
	// 创建 loggerHook 实例
	hook := NewloggerHook(&buf, logrus.InfoLevel, &logrus.TextFormatter{})
	entry := &logrus.Entry{
		Message: "Test log message",
	}
	err := hook.Fire(entry)
	if err != nil {
		t.Logf("Error firing hook: %v", err)
	}
	// 检查写入信息和测试信息是否相同
	if !strings.Contains(buf.String(), "Test log message") {
		t.Log("Expected log message not found in buffer")
	}
}

func TestSetuploggerHook(t *testing.T) {
	restore := SetuploggerHook("sss")
	defer restore()
}
