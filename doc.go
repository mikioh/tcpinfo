// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tcpinfo implements encoding and decoding of TCP-level
// socket options regarding connection information.
//
// The Transmission Control Protocol (TCP) is defined in RFC 793.
// TCP Selective Acknowledgment Options is defined in RFC 2018.
// Management Information Base for the Transmission Control Protocol
// (TCP) is defined in RFC 4022.
// TCP Congestion Control is defined in RFC 5681.
// Computing TCP's Retransmission Timer is described in RFC 6298.
// TCP Options and Maximum Segment Size (MSS) is defined in RFC 6691.
// Shared Use of Experimental TCP Options is defined in RFC 6994.
// TCP Extensions for High Performance is defined in RFC 7323.
//
// Example:
//
//	import (
//		"github.com/mikioh/tcp"
//		"github.com/mikioh/tcpinfo"
//	)
//
//	c, err := net.Dial("tcp", "golang.org:80")
//	if err != nil {
//		// error handling
//	}
//	defer c.Close()
//
//	tc, err := tcp.NewConn(c)
//	if err != nil {
//		// error handling
//	}
//	var o tcpinfo.Info
//	var b [256]byte
//	i, err := tc.Option(o.Level(), o.Name(), b[:])
//	if err != nil {
//		// error handling
//	}
//	txt, err := json.Marshal(i)
//	if err != nil {
//		// error handling
//	}
//	fmt.Println(txt)
package tcpinfo
