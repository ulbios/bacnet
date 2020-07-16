// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/kazukiigeta/bacnet"
)

func main() {

	var (
		addr     = flag.String("raddr", "127.0.0.1:47808", "Remote IP and Port to connect to.")
		interval = flag.Duration("whois-interval", 1, "Interval (sec) for Unconfirmed request of WhoIs.")
	)
	flag.Parse()

	udpAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %s", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Failed to dial UDP: %s", err)
	}
	defer conn.Close()

	bvlc := bacnet.NewBVLC(bacnet.BVLCFuncBroadcast)
	npdu := bacnet.NewNPDU(false, false, false, false)
	u := bacnet.NewUnconfirmedWhoIs(bvlc, npdu)
	b, err := u.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to marshal binary of Unconfimed request WhoIs packet: %s", err)
	}

	for {
		if _, err := conn.Write(b); err != nil {
			log.Fatalf("Failed to write Unconfimed request WhoIs packet: %s", err)
		}
		log.Printf("Sent: %x", b)
		time.Sleep(*interval * time.Second)
	}
}
