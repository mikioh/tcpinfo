// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd netbsd

package tcpinfo

import (
	"runtime"
	"time"
	"unsafe"

	"github.com/mikioh/tcpopt"
)

var options = [soMax]option{
	soInfo:   {ianaProtocolTCP, sysTCP_INFO, parseInfo},
	soCCInfo: {ianaProtocolTCP, -1, nil},
	soCCAlgo: {ianaProtocolTCP, -1, nil},
}

// Marshal implements the Marshal method of tcpopt.Option interface.
func (i *Info) Marshal() ([]byte, error) { return (*[sizeofTCPInfo]byte)(unsafe.Pointer(i))[:], nil }

// A SysInfo represents platform-specific information.
type SysInfo struct {
	SenderWindowBytes uint `json:"snd_wnd_bytes"`   // advertised sender window in bytes
	SenderWindowSegs  uint `json:"snd_wnd_segs"`    // advertised sender window in # of segments
	NextEgressSeq     uint `json:"egress_seq"`      // next egress seq. number
	NextIngressSeq    uint `json:"ingress_seq"`     // next ingress seq. number
	RetransSegs       uint `json:"retrans_segs"`    // # of retransmit segments sent
	OutOfOrderSegs    uint `json:"ooo_segs"`        // # of out-of-order segments received
	ZeroWindowUpdates uint `json:"zerownd_updates"` // # of zero-window updates sent
	Offloading        bool `json:"offloading"`      // TCP offload processing
}

var sysStates = [11]State{Closed, Listen, SynSent, SynReceived, Established, CloseWait, FinWait1, Closing, LastAck, FinWait2, TimeWait}

func parseInfo(b []byte) (tcpopt.Option, error) {
	if len(b) < sizeofTCPInfo {
		return nil, errBufferTooShort
	}
	sti := (*sysTCPInfo)(unsafe.Pointer(&b[0]))
	i := &Info{State: sysStates[sti.State]}
	if sti.Options&sysTCPI_OPT_WSCALE != 0 {
		i.Options = append(i.Options, WindowScale(sti.Pad_cgo_0[0]>>4))
		i.PeerOptions = append(i.PeerOptions, WindowScale(sti.Pad_cgo_0[0]&0x0f))
	}
	if sti.Options&sysTCPI_OPT_SACK != 0 {
		i.Options = append(i.Options, SACKPermitted(true))
		i.PeerOptions = append(i.PeerOptions, SACKPermitted(true))
	}
	if sti.Options&sysTCPI_OPT_TIMESTAMPS != 0 {
		i.Options = append(i.Options, Timestamps(true))
		i.PeerOptions = append(i.PeerOptions, Timestamps(true))
	}
	i.SenderMSS = MaxSegSize(sti.Snd_mss)
	i.ReceiverMSS = MaxSegSize(sti.Rcv_mss)
	i.RTT = time.Duration(sti.Rtt) * time.Microsecond
	i.RTTVar = time.Duration(sti.Rttvar) * time.Microsecond
	i.RTO = time.Duration(sti.Rto) * time.Microsecond
	i.ATO = time.Duration(sti.X__tcpi_ato) * time.Microsecond
	i.LastDataSent = time.Duration(sti.X__tcpi_last_data_sent) * time.Microsecond
	i.LastDataReceived = time.Duration(sti.Last_data_recv) * time.Microsecond
	i.LastAckReceived = time.Duration(sti.X__tcpi_last_ack_recv) * time.Microsecond
	i.FlowControl = &FlowControl{
		ReceiverWindow: uint(sti.Rcv_space),
	}
	i.CongestionControl = &CongestionControl{
		SenderSSThreshold:   uint(sti.Snd_ssthresh),
		ReceiverSSThreshold: uint(sti.X__tcpi_rcv_ssthresh),
	}
	i.Sys = &SysInfo{
		NextEgressSeq:     uint(sti.Snd_nxt),
		NextIngressSeq:    uint(sti.Rcv_nxt),
		RetransSegs:       uint(sti.Snd_rexmitpack),
		OutOfOrderSegs:    uint(sti.Rcv_ooopack),
		ZeroWindowUpdates: uint(sti.Snd_zerowin),
	}
	if sti.Options&sysTCPI_OPT_TOE != 0 {
		i.Sys.Offloading = true
	}
	switch runtime.GOOS {
	case "freebsd":
		i.CongestionControl.SenderWindowBytes = uint(sti.Snd_cwnd)
		i.Sys.SenderWindowBytes = uint(sti.Snd_wnd)
	case "netbsd":
		i.CongestionControl.SenderWindowSegs = uint(sti.Snd_cwnd)
		i.Sys.SenderWindowSegs = uint(sti.Snd_wnd)
	}
	return i, nil
}

func parseCCAlgorithmInfo(name string, b []byte) (CCAlgorithmInfo, error) { return nil, errOpNoSupport }
