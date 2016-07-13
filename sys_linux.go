// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcpinfo

import (
	"strings"
	"time"
	"unsafe"

	"github.com/mikioh/tcpopt"
)

var options = [soMax]option{
	soInfo:   {ianaProtocolTCP, sysTCP_INFO, parseInfo},
	soCCInfo: {ianaProtocolTCP, sysTCP_CC_INFO, parseCCInfo},
	soCCAlgo: {ianaProtocolTCP, sysTCP_CONGESTION, parseCCAlgorithm},
}

// Marshal implements the Marshal method of tcpopt.Option interface.
func (i *Info) Marshal() ([]byte, error) { return (*[sizeofTCPInfo]byte)(unsafe.Pointer(i))[:], nil }

// A CAState represents a state of congestion avoidance.
type CAState int

var caStates = map[CAState]string{
	CAOpen:     "open",
	CADisorder: "disorder",
	CACWR:      "congestion window reduced",
	CARecovery: "recovery",
	CALoss:     "loss",
}

func (st CAState) String() string {
	s, ok := caStates[st]
	if !ok {
		return "<nil>"
	}
	return s
}

// A SysInfo represents platform-specific information.
type SysInfo struct {
	PathMTU         uint       `json:"path_mtu"`     // path maximum transmission unit
	AdvertisedMSS   MaxSegSize `json:"adv_mss"`      // advertised maximum segment size
	CAState         CAState    `json:"ca_state"`     // state of congestion avoidance
	KeepAliveProbes uint       `json:"ka_probes"`    // # of keep alive probes sent
	UnackSegs       uint       `json:"unack_segs"`   // # of unack'd segments in transmission queue
	SackSegs        uint       `json:"sack_segs"`    // # of sack'd segments in tranmission queue
	LostSegs        uint       `json:"lost_segs"`    // # of lost segments in transmission queue
	RetransSegs     uint       `json:"retrans_segs"` // # of retransmitting segments in transmission queue
	ForwardAckSegs  uint       `json:"fack_segs"`    // # of forward ack'd segments in transmission queue
}

var sysStates = [12]State{Unknown, Established, SynSent, SynReceived, FinWait1, FinWait2, TimeWait, Closed, CloseWait, LastAck, Listen, Closing}

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
	if sti.Options&sysTCPI_OPT_TIMESTAMPS != 0 {
		i.Options = append(i.Options, Timestamps(true))
		i.PeerOptions = append(i.PeerOptions, Timestamps(true))
	}
	i.SenderMSS = MaxSegSize(sti.Snd_mss)
	i.ReceiverMSS = MaxSegSize(sti.Rcv_mss)
	i.RTT = time.Duration(sti.Rtt) * time.Microsecond
	i.RTTVar = time.Duration(sti.Rttvar) * time.Microsecond
	i.RTO = time.Duration(sti.Rto) * time.Microsecond
	i.ATO = time.Duration(sti.Ato) * time.Microsecond
	i.LastDataSent = time.Duration(sti.Last_data_sent) * time.Millisecond
	i.LastDataReceived = time.Duration(sti.Last_data_recv) * time.Millisecond
	i.LastAckReceived = time.Duration(sti.Last_ack_recv) * time.Millisecond
	i.FlowControl = &FlowControl{
		ReceiverWindow: uint(sti.Rcv_space),
	}
	i.CongestionControl = &CongestionControl{
		SenderSSThreshold:   uint(sti.Snd_ssthresh),
		ReceiverSSThreshold: uint(sti.Rcv_ssthresh),
		SenderWindow:        uint(sti.Snd_cwnd),
	}
	i.Sys = &SysInfo{
		PathMTU:         uint(sti.Pmtu),
		AdvertisedMSS:   MaxSegSize(sti.Advmss),
		CAState:         CAState(sti.Ca_state),
		KeepAliveProbes: uint(sti.Probes),
		UnackSegs:       uint(sti.Unacked),
		SackSegs:        uint(sti.Sacked),
		LostSegs:        uint(sti.Lost),
		RetransSegs:     uint(sti.Retrans),
		ForwardAckSegs:  uint(sti.Fackets),
	}
	return i, nil
}

// A VegasInfo represents Vegas congestion control information.
type VegasInfo struct {
	Enabled    bool          `json:"enabled"`
	RoundTrips uint          `json:"rnd_trips"` // # of round-trips
	RTT        time.Duration `json:"rtt"`       // round-trip time
	MinRTT     time.Duration `json:"min_rtt"`   // minimum round-trip time
}

// Algorithm implements the Algorithm method of CCAlgorithmInfo
// interface.
func (vi *VegasInfo) Algorithm() string { return "vegas" }

// A CEState represents a state of ECN congestion encountered (CE)
// codepoint.
type CEState int

// A DCTCPInfo represents Datacenter TCP congestion control
// information.
type DCTCPInfo struct {
	Enabled         bool    `json:"enabled"`
	CEState         CEState `json:"ce_state"`    // state of ECN CE codepoint
	Alpha           uint    `json:"alpha"`       // fraction of bytes sent
	ECNAckedBytes   uint    `json:"ecn_acked"`   // # of acked bytes with ECN
	TotalAckedBytes uint    `json:"total_acked"` // total # of acked bytes
}

// Algorithm implements the Algorithm method of CCAlgorithmInfo
// interface.
func (di *DCTCPInfo) Algorithm() string { return "dctcp" }

func parseCCAlgorithmInfo(name string, b []byte) (CCAlgorithmInfo, error) {
	if strings.HasPrefix(name, "dctcp") {
		if len(b) < sizeofTCPDCTCPInfo {
			return nil, errBufferTooShort
		}
		sdi := (*sysTCPDCTCPInfo)(unsafe.Pointer(&b[0]))
		di := &DCTCPInfo{Alpha: uint(sdi.Alpha)}
		if sdi.Enabled != 0 {
			di.Enabled = true
		}
		return di, nil
	}
	if len(b) < sizeofTCPVegasInfo {
		return nil, errBufferTooShort
	}
	svi := (*sysTCPVegasInfo)(unsafe.Pointer(&b[0]))
	vi := &VegasInfo{
		RoundTrips: uint(svi.Rttcnt),
		RTT:        time.Duration(svi.Rtt) * time.Microsecond,
		MinRTT:     time.Duration(svi.Minrtt) * time.Microsecond,
	}
	if svi.Enabled != 0 {
		vi.Enabled = true
	}
	return vi, nil
}
