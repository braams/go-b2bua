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
    "sippy/time"
    "sippy/types"
)

type UasStateTrying struct {
    *uaStateGeneric
}

func NewUasStateTrying(ua sippy_types.UA, config sippy_conf.Config) *UasStateTrying {
    return &UasStateTrying{
        uaStateGeneric : newUaStateGeneric(ua, config),
    }
}

func (self *UasStateTrying) OnActivation() {
}

func (self *UasStateTrying) String() string {
    return "Trying(UAS)"
}

func (self *UasStateTrying) RecvEvent(_event sippy_types.CCEvent) (sippy_types.UaState, error) {
    eh := _event.GetExtraHeaders()
    switch event := _event.(type) {
    case *CCEventRing:
        code, reason, body := event.scode, event.scode_reason, event.body
        if code == 0 {
            code, reason, body = 180, "Ringing", nil
        } else {
            if code == 100 {
                return nil, nil
            }
            if body != nil && self.ua.HasOnLocalSdpChange() && body.NeedsUpdate() {
                self.ua.OnLocalSdpChange(body, event, func (sippy_types.MsgBody) { self.ua.RecvEvent(event) })
                return nil, nil
            }
        }
        self.ua.SetLSDP(body)
        self.ua.SendUasResponse(nil, code, reason, body, nil, false, eh...)
        if self.ua.HasNoProgressTimer() {
            self.ua.CancelNoProgressTimer()
            if self.ua.GetExMtime() != nil {
                self.ua.StartExpireTimer(self.ua.GetExMtime())
            }
        }
        if self.ua.GetP1xxTs() == nil {
            self.ua.SetP1xxTs(event.GetRtime())
        }
        return NewUasStateRinging(self.ua, event.GetRtime(), event.GetOrigin(), code, self.config), nil
    case *CCEventPreConnect:
        code, reason, body := event.scode, event.scode_reason, event.body
        if body != nil && self.ua.HasOnLocalSdpChange() && body.NeedsUpdate() {
            self.ua.OnLocalSdpChange(body, event, func (sippy_types.MsgBody) { self.ua.RecvEvent(event) })
            return nil, nil
        }
        self.ua.SetLSDP(body)
        self.ua.CancelNoProgressTimer()
        self.ua.SendUasResponse(nil, code, reason, body, self.ua.GetLContacts(), /*ack_wait*/ true, eh...)
        return NewUaStateConnected(self.ua, nil, "", self.config), nil
    case *CCEventConnect:
        code, reason, body := event.scode, event.scode_reason, event.body
        if body != nil && self.ua.HasOnLocalSdpChange() && body.NeedsUpdate() {
            self.ua.OnLocalSdpChange(body, event, func (sippy_types.MsgBody) { self.ua.RecvEvent(event) })
            return nil, nil
        }
        self.ua.SetLSDP(body)
        self.ua.SendUasResponse(nil, code, reason, body, self.ua.GetLContacts(), false, eh...)
        self.ua.CancelExpireTimer()
        self.ua.CancelNoProgressTimer()
        self.ua.StartCreditTimer(event.GetRtime())
        self.ua.SetConnectTs(event.GetRtime())
        return NewUaStateConnected(self.ua, event.GetRtime(), event.GetOrigin(), self.config), nil
    case *CCEventRedirect:
        self.ua.SendUasResponse(nil, event.scode, event.scode_reason, event.body, event.GetContacts(), false, eh...)
        self.ua.CancelExpireTimer()
        self.ua.CancelNoProgressTimer()
        self.ua.SetDisconnectTs(event.GetRtime())
        return NewUaStateFailed(self.ua, event.GetRtime(), event.GetOrigin(), event.scode, self.config), nil
    case *CCEventFail:
        code, reason := event.scode, event.scode_reason
        if code == 0 {
            code, reason = 500, "Failed"
        }
        self.ua.SendUasResponse(nil, code, reason, nil, nil, false, eh...)
        self.ua.CancelExpireTimer()
        self.ua.CancelNoProgressTimer()
        self.ua.SetDisconnectTs(event.GetRtime())
        return NewUaStateFailed(self.ua, event.GetRtime(), event.GetOrigin(), code, self.config), nil
    case *CCEventDisconnect:
        self.ua.SendUasResponse(nil, 500, "Disconnected", nil, nil, false, eh...)
        self.ua.CancelExpireTimer()
        self.ua.CancelNoProgressTimer()
        self.ua.SetDisconnectTs(event.GetRtime())
        return NewUaStateDisconnected(self.ua, event.GetRtime(), event.GetOrigin(), self.ua.GetLastScode(), nil, self.config), nil
    }
    //return nil, fmt.Errorf("uas-trying: wrong event %s in the Trying state", _event.String())
    return nil, nil
}

func (self *UasStateTrying) Cancel(rtime *sippy_time.MonoTime, req sippy_types.SipRequest) {
    event := NewCCEventDisconnect(nil, rtime, self.ua.GetOrigin())
    if req != nil {
        event.SetReason(req.GetReason())
    }
    self.ua.SetDisconnectTs(rtime)
    self.ua.ChangeState(NewUaStateDisconnected(self.ua, rtime, self.ua.GetOrigin(), 0, req, self.config))
    self.ua.EmitEvent(event)
}
