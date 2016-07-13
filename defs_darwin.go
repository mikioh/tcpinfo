// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package tcpinfo

/*
#include <netinet/tcp.h>
*/
import "C"

const (
	sysTCP_CONNECTION_INFO = C.TCP_CONNECTION_INFO

	sysTCPCI_OPT_TIMESTAMPS           = C.TCPCI_OPT_TIMESTAMPS
	sysTCPCI_OPT_SACK                 = C.TCPCI_OPT_SACK
	sysTCPCI_OPT_WSCALE               = C.TCPCI_OPT_WSCALE
	sysTCPCI_OPT_ECN                  = C.TCPCI_OPT_ECN
	sysTCPCI_FLAG_LOSSRECOVERY        = C.TCPCI_FLAG_LOSSRECOVERY
	sysTCPCI_FLAG_REORDERING_DETECTED = C.TCPCI_FLAG_REORDERING_DETECTED

	sizeofTCPConnInfo = C.sizeof_struct_tcp_connection_info
)

type sysTCPConnInfo C.struct_tcp_connection_info