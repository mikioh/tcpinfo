// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tcpinfo implements encoding and decoding of TCP-level
// socket options regarding connection information.
//
// Example:
//
//	import "github.com/mikioh/tcp"
//	import "github.com/mikioh/tcpinfo"
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
