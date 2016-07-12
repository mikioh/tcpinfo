// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !darwin,!freebsd,!linux,!netbsd

package tcpinfo

import "github.com/mikioh/tcpopt"

var options = [soMax]option{
	soInfo:   {ianaProtocolTCP, -1, nil},
	soCCInfo: {ianaProtocolTCP, -1, nil},
	soCCAlgo: {ianaProtocolTCP, -1, nil},
}

// Marshal implements the Marshal method of tcpopt.Option interface.
func (i *Info) Marshal() ([]byte, error) { return nil, errOpNoSupport }

// A SysInfo represents platform-specific information.
type SysInfo struct{}

func parseInfo(b []byte) (tcpopt.Option, error)                           { return nil, errOpNoSupport }
func parseCCAlgorithmInfo(name string, b []byte) (CCAlgorithmInfo, error) { return nil, errOpNoSupport }
