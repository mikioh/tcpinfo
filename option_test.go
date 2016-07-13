// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcpinfo_test

import (
	"runtime"
	"testing"

	"github.com/mikioh/tcpinfo"
	"github.com/mikioh/tcpopt"
)

func TestOption(t *testing.T) {
	opts := make([]tcpopt.Option, 0, 3)
	switch runtime.GOOS {
	case "darwin", "freebsd", "netbsd":
		opts = append(opts, &tcpinfo.Info{})
	case "linux":
		opts = append(opts, &tcpinfo.Info{})
		opts = append(opts, &tcpinfo.CCInfo{})
		opts = append(opts, tcpinfo.CCAlgorithm("vegas"))
	default:
		t.Skipf("%s/%s", runtime.GOOS, runtime.GOARCH)
	}

	for _, o := range opts {
		if o.Level() <= 0 {
			t.Fatalf("got %#x; want greater than zero", o.Level())
		}
		if o.Name() <= 0 {
			t.Fatalf("got %#x; want greater than zero", o.Name())
		}
		b, err := o.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if len(b) == 0 {
			continue
		}
		oo, err := tcpopt.Parse(o.Level(), o.Name(), b)
		if err != nil {
			t.Fatal(err)
		}
		if oo, ok := oo.(*tcpinfo.Info); ok {
			if _, err := oo.MarshalJSON(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestParseBufferOverrun(t *testing.T) {
	for _, o := range []tcpopt.Option{
		&tcpinfo.Info{},
		&tcpinfo.CCInfo{},
		tcpinfo.CCAlgorithm("vegas"),
	} {
		var b [3]byte
		tcpopt.Parse(o.Level(), o.Name(), b[:])
	}
}
