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
package sippy_sdp

import (
    "strings"

    "sippy/net"
)

type SdpMedia struct {
    stype       string
    port        string
    transport   string
    formats     []string
}

func ParseSdpMedia(body string) *SdpMedia {
    if body == "" {
        return nil
    }
    params := strings.Fields(body)
    if len(params) < 3 {
        return nil
    }
    return &SdpMedia{
        stype       : params[0],
        port        : params[1],
        transport   : params[2],
        formats     : params[3:],
    }
}

func (self *SdpMedia) String() string {
    rval := self.stype + " " + self.port + " " + self.transport
    for _, format := range self.formats {
        rval += " " + format
    }
    return rval
}

func (self *SdpMedia) LocalStr(hostport *sippy_net.HostPort) string {
    return self.String()
}

func (self *SdpMedia) GetCopy() *SdpMedia {
    if self == nil {
        return nil
    }
    formats := make([]string, len(self.formats))
    copy(formats, self.formats)
    return &SdpMedia{
        stype       : self.stype,
        port        : self.port,
        transport   : self.transport,
        formats     : formats,
    }
}

func (self *SdpMedia) GetTransport() string {
    return self.transport
}

func (self *SdpMedia) GetPort() string {
    return self.port
}

func (self *SdpMedia) SetPort(port string) {
    self.port = port
}

func (self *SdpMedia) HasFormat(format string) bool {
    for _, f := range self.formats {
        if f == format {
            return true
        }
    }
    return false
}

func (self *SdpMedia) GetFormats() []string {
    return self.formats
}

// WARNING! Use this function only if know what you do!
// Otherwise consider using the sdpMediaDescription.SetFormats() instead.
func (self *SdpMedia) SetFormats(formats []string) {
    self.formats = formats
}
