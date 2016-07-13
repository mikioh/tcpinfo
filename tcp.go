// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcpinfo

// An OptionKind represents an option kind.
type OptionKind int

const (
	KindMaxSegSize    OptionKind = 2
	KindWindowScale   OptionKind = 3
	KindSACKPermitted OptionKind = 4
	KindTimestamps    OptionKind = 8
)

var optionKinds = map[OptionKind]string{
	KindMaxSegSize:    "mss",
	KindWindowScale:   "wscale",
	KindSACKPermitted: "sack perm",
	KindTimestamps:    "tmstamps",
}

func (k OptionKind) String() string {
	s, ok := optionKinds[k]
	if !ok {
		return "<nil>"
	}
	return s
}

// An Option represents an option.
type Option interface {
	Kind() OptionKind
}

// A MaxSegSize represents a maxiumum sengment size option.
type MaxSegSize uint

// Kind returns an option kind field.
func (mss MaxSegSize) Kind() OptionKind { return KindMaxSegSize }

// A WindowScale represents a windows scale option.
type WindowScale int

// Kind returns an option kind field.
func (ws WindowScale) Kind() OptionKind { return KindWindowScale }

// A SACKPermitted reports whether a selective acknowledgment
// permitted option is enabled.
type SACKPermitted bool

// Kind returns an option kind field.
func (sp SACKPermitted) Kind() OptionKind { return KindSACKPermitted }

// A Timestamps reports whether a timestamps option is enabled.
type Timestamps bool

// Kind returns an option kind field.
func (ts Timestamps) Kind() OptionKind { return KindTimestamps }
