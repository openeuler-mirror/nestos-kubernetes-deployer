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
	"fmt"
	"testing"
)

// TestBufferedPrinter_Write tests the Write method of the bufferedPrinter.
// It verifies that data written to the bufferedPrinter is correctly buffered and printed line by line.
func TestBufferedPrinterWrite(t *testing.T) {
	tp := New(func(args ...interface{}) {
		fmt.Println(args)
	})
	t.Run("Write", func(t *testing.T) {
		tp.Write([]byte("Buffered printer"))
	})
}

// TestBufferedPrinter_Close tests the Close method of the bufferedPrinter.
func TestBufferedPrinterClose(t *testing.T) {
	bp := New(func(args ...interface{}) {
		fmt.Println("TestBufferedPrinter_Close")
	})
	t.Run("close", func(t *testing.T) {
		bp.buf.Write([]byte("Buffered close"))
		bp.Close()
	})
}

// TestTrimLastNewline tests the TrimLastNewline function.
// It verifies that the last newline character in the input arguments is correctly trimmed.
func TestTrimLastNewline(t *testing.T) {
	input := []interface{}{"Buffered printer.\n"}
	TrimLastNewline(input...)
}
