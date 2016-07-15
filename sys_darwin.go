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
func (i *Info) Marshal() ([]byte, error) {
	return (*[sizeofTCPConnectionInfo]byte)(unsafe.Pointer(i))[:], nil
}

type SysFlags uint

func (f SysFlags) String() string {
	s := ""
	for i, name := range []string{
		"loss recovery",
		"reordering detected",
	} {
		if f&(1<<uint(i)) != 0 {
			if s != "" {
				s += "|"
			}
			s += name
		}
	}
	if s == "" {
		s = "0"
	}
	return s
}

// A SysInfo represents platform-specific information.
type SysInfo struct {
	Flags                   SysFlags      `json:"flags"`           // flags
	SenderWindowBytes       uint          `json:"snd_wnd_bytes"`   // advertised sender window in bytes
	SenderInUseBytes        uint          `json:"snd_inuse_bytes"` // # of bytes in send buffer including inflight data
	SRTT                    time.Duration `json:"srtt"`            // smoothed round-trip time
	SegsSent                uint64        `json:"segs_sent"`       // # of segements send
	BytesSent               uint64        `json:"bytes_sent"`      // # of bytes sent
	RetransBytes            uint64        `json:"retrans_bytes"`   // # of retransmitted bytes
	SegsReceived            uint64        `json:"segs_rcvd"`       // # of segments received
	BytesReceived           uint64        `json:"bytes_rcvd"`      // # of bytes received
	OutOfOrderBytesReceived uint64        `json:"ooo_bytes_rcvd"`  // # of our-of-order bytes received
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseInfo(b []byte) (tcpopt.Option, error) {
	if len(b) < sizeofTCPConnectionInfo {
		return nil, errBufferTooShort
	}
	sti := (*sysTCPConnectionInfo)(unsafe.Pointer(&b[0]))
	i := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPCI_OPT_WSCALE != 0 {
		i.Options = append(i.Options, WindowScale(sti.Snd_wscale))
		i.PeerOptions = append(i.PeerOptions, WindowScale(sti.Rcv_wscale))
	}
	if sti.Options&sysTCPCI_OPT_SACK != 0 {
		i.Options = append(i.Options, SACKPermitted(true))
		i.PeerOptions = append(i.PeerOptions, SACKPermitted(true))
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
		SenderWindowBytes: uint(sti.Snd_cwnd),
	}
	i.Sys = &SysInfo{
		Flags:                   SysFlags(sti.Flags),
		SenderWindowBytes:       uint(sti.Snd_wnd),
		SenderInUseBytes:        uint(sti.Snd_sbbytes),
		SRTT:                    time.Duration(sti.Srtt) * time.Millisecond,
		SegsSent:                uint64(sti.Txpackets),
		BytesSent:               uint64(sti.Txbytes),
		RetransBytes:            uint64(sti.Txretransmitbytes),
		SegsReceived:            uint64(sti.Rxpackets),
		BytesReceived:           uint64(sti.Rxbytes),
		OutOfOrderBytesReceived: uint64(sti.Rxoutoforderbytes),
	}
	return i, nil
}

func parseCCAlgorithmInfo(name string, b []byte) (CCAlgorithmInfo, error) { return nil, errOpNoSupport }
