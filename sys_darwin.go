// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcpinfo

import (
	"time"
	"unsafe"

	"github.com/mikioh/tcpopt"
)

var options = [soMax]option{
	soInfo:   {ianaProtocolTCP, sysTCP_CONNECTION_INFO, parseInfo},
	soCCInfo: {ianaProtocolTCP, -1, nil},
	soCCAlgo: {ianaProtocolTCP, -1, nil},
}

// Marshal implements the Marshal method of tcpopt.Option interface.
func (i *Info) Marshal() ([]byte, error) { return (*[sizeofTCPConnInfo]byte)(unsafe.Pointer(i))[:], nil }

// A SysInfo represents platform-specific information.
type SysInfo struct {
	Flags        uint          `json:"flags"`     // flags
	SenderWindow uint          `json:"snd_wnd"`   // advertised sender window in bytes
	SenderInUse  uint          `json:"snd_inuse"` // bytes in send buffer including inflight data
	SRTT         time.Duration `json:"srtt"`      // smoothed round-trip time
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseInfo(b []byte) (tcpopt.Option, error) {
	sti := (*sysTCPConnInfo)(unsafe.Pointer(&b[0]))
	i := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPCI_OPT_WSCALE != 0 {
		i.Options = append(i.Options, WindowScale(sti.Snd_wscale))
		i.PeerOptions = append(i.PeerOptions, WindowScale(sti.Rcv_wscale))
	}
	if sti.Options&sysTCPCI_OPT_TIMESTAMPS != 0 {
		i.Options = append(i.Options, Timestamps(true))
		i.PeerOptions = append(i.PeerOptions, Timestamps(true))
	}
	i.SenderMSS = MaxSegSize(sti.Maxseg)
	i.ReceiverMSS = MaxSegSize(sti.Maxseg)
	i.RTT = time.Duration(sti.Rttcur) * time.Millisecond
	i.RTTVar = time.Duration(sti.Rttvar) * time.Millisecond
	i.RTO = time.Duration(sti.Rto) * time.Millisecond
	i.FlowControl = &FlowControl{
		ReceiverWindow: uint(sti.Rcv_wnd),
	}
	i.CongestionControl = &CongestionControl{
		SenderSSThreshold: uint(sti.Snd_ssthresh),
		SenderWindow:      uint(sti.Snd_cwnd),
	}
	i.Sys = &SysInfo{
		Flags:        uint(sti.Flags),
		SenderWindow: uint(sti.Snd_wnd),
		SenderInUse:  uint(sti.Snd_sbbytes),
		SRTT:         time.Duration(sti.Srtt) * time.Millisecond,
	}
	return i, nil
}

func parseCCAlgorithmInfo(name string, b []byte) (CCAlgorithmInfo, error) { return nil, errOpNoSupport }
