// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcpinfo

import (
	"encoding/binary"
	"unsafe"

	"github.com/mikioh/tcpopt"
)

var nativeEndian binary.ByteOrder

func init() {
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		nativeEndian = binary.LittleEndian
	} else {
		nativeEndian = binary.BigEndian
	}
	for _, o := range options {
		if o.name <= 0 || o.parseFn == nil {
			continue
		}
		tcpopt.Register(o.level, o.name, o.parseFn)
	}
}

const (
	ianaProtocolTCP = 0x6
)

const (
	soInfo = iota
	soCCInfo
	soCCAlgo
	soMax
)

// An option represents a binding for socket option.
type option struct {
	level   int // option level
	name    int // option name, must be equal or greater than 1
	parseFn func([]byte) (tcpopt.Option, error)
}
