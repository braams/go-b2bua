// Copyright (c) 2003-2005 Maxim Sobolev. All rights reserved.
// Copyright (c) 2006-2015 Sippy Software, Inc. All rights reserved.
// Copyright (c) 2015 Andrii Pylypenko. All rights reserved.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation and/or
// other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package sippy_log

import (
    "fmt"
    "runtime"
    "os"
    "strings"
    "sync"
    "time"
)

type ErrorLogger interface {
    ErrorAndTraceback(interface{})
    Error(...interface{})
    Debug(...interface{})
    Errorf(string, ...interface{})
    Debugf(string, ...interface{})
}

type errorLogger struct {
    lock    sync.Mutex
}

func NewErrorLogger() *errorLogger {
    return &errorLogger{}
}

func (self *errorLogger) ErrorAndTraceback(err interface{}) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.Error(err)
    buf := make([]byte, 16384)
    n := runtime.Stack(buf, false)
    s := string(buf[:n])
    for _, l := range strings.Split(s, "\n") {
        self.Error(l)
    }
}

func (self *errorLogger) Debug(params...interface{}) {
    self.write("DEBUG:", params...)
}

func (self *errorLogger) Debugf(format string, params...interface{}) {
    self.write("DEBUG:", fmt.Sprintf(format, params...))
}

func (self *errorLogger) Error(params...interface{}) {
    self.write("ERROR:", params...)
}

func (self *errorLogger) Errorf(format string, params...interface{}) {
    self.write("ERROR:", fmt.Sprintf(format, params...))
}

func (*errorLogger) Reopen() {
}

func (*errorLogger) write(prefix string, params ...interface{}) {
    t := time.Now().UTC()
    buf := []interface{}{ fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d+00", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()), " ", prefix }
    for _, it := range params {
        buf = append(buf, " ", it)
    }
    buf = append(buf, "\n")
    fmt.Fprint(os.Stderr, buf...)
}
