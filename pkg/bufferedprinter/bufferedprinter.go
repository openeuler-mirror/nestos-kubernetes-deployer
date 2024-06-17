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

package bufferedprinter

import (
	"bufio"
	"bytes"
	"strings"
	"sync"
)

type print func(args ...interface{})

type bufferedPrinter struct {
	buf   bytes.Buffer
	print print

	sync.Mutex
	once sync.Once
}

func New(print print) *bufferedPrinter {
	return &bufferedPrinter{
		print: print,
	}
}

func (bp *bufferedPrinter) Write(p []byte) (int, error) {
	bp.Lock()
	defer bp.Unlock()

	n, err := bp.buf.Write(p)
	if err != nil {
		return n, err
	}

	scanner := bufio.NewScanner(&bp.buf)
	for scanner.Scan() {
		bp.print(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return n, err
	}

	return n, nil
}

func (bp *bufferedPrinter) Close() error {
	bp.once.Do(func() {
		bp.Lock()
		defer bp.Unlock()

		line := bp.buf.String()
		if len(line) > 0 {
			bp.print(line)
		}
	})

	return nil
}

func TrimLastNewline(args ...interface{}) []interface{} {
	if len(args) > 0 {
		if lastArg, ok := args[len(args)-1].(string); ok {
			args[len(args)-1] = strings.TrimRight(lastArg, "\n")
		}
	}
	return args
}
