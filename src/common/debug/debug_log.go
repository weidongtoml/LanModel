// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"fmt"
	"log"
	"runtime"
)

var DLog DebugLog

func Print(msg string) {
	DLog.Print(msg)
}

func Printf(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	Print(str)
}

type DebugLog struct {
}

func (d *DebugLog) Print(msg string) {
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		file = "Unknown"
		line = 0
	}
	log.Printf("%s %d: %s\n", file, line, msg)
}
