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
package sippy

import (
    "sippy/conf"
    "sippy/headers"
    "sippy/net"
    "sippy/types"
)

type statefulProxy struct {
    sip_tm      sippy_types.SipTransactionManager
    destination *sippy_net.HostPort
    config      sippy_conf.Config
}

func NewStatefulProxy(sip_tm sippy_types.SipTransactionManager, destination *sippy_net.HostPort, config sippy_conf.Config) *statefulProxy {
    return &statefulProxy{
        sip_tm      : sip_tm,
        destination : destination,
        config      : config,
    }
}

func (self *statefulProxy) RecvRequest(req sippy_types.SipRequest, t sippy_types.ServerTransaction) *sippy_types.Ua_context {
    via0 := sippy_header.NewSipVia(self.config)
    via0_body, _ := via0.GetBody()
    via0_body.GenBranch()
    req.InsertFirstVia(via0)
    req.SetTarget(self.destination)
    //print req
    self.sip_tm.BeginNewClientTransaction(req, self, nil, nil, nil, nil)
    return &sippy_types.Ua_context{}
}

func (self *statefulProxy) RecvResponse(resp sippy_types.SipResponse, t sippy_types.ClientTransaction) {
    resp.RemoveFirstVia()
    self.sip_tm.SendResponse(resp, /*lock*/true, nil)
}
