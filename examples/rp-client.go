// Copyright 2020 bacnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/ulbios/bacnet"
	"github.com/ulbios/bacnet/services"
)

func init() {
	ReadPropertyClientCmd.Flags().Uint16Var(&rpObjectType, "object-type", 1, "Object type to read.")
	ReadPropertyClientCmd.Flags().Uint32Var(&rpInstanceId, "instance-id", 0, "Instance ID to read.") // Analog-input
	ReadPropertyClientCmd.Flags().Uint8Var(&rpPropertyId, "property-id", 85, "Property ID to read.") // Current-value
	ReadPropertyClientCmd.Flags().IntVar(&rpPeriod, "period", 1, "Period, in seconds, between requests.")
	ReadPropertyClientCmd.Flags().IntVar(&rpN, "messages", 1, "Number of messages to send, being 0 unlimited.")
}

var (
	rpObjectType uint16
	rpInstanceId uint32
	rpPropertyId uint8
	rpPeriod     int
	rpN          int

	ReadPropertyClientCmd = &cobra.Command{
		Use:   "rpc",
		Short: "Send ReadProperty requests.",
		Long:  "There's not much more really. This command sends a configurable ReadProperty request.",
		Args:  argValidation,
		Run:   ReadPropertyClientExample,
	}
)

func ReadPropertyClientExample(cmd *cobra.Command, args []string) {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %s", err)
	}

	listenConn, err := net.ListenPacket("udp", bAddr)
	if err != nil {
		log.Fatalf("failed to begin listening for packets: %v\n", err)
	}
	defer listenConn.Close()

	mReadProperty, err := bacnet.NewReadProperty(rpObjectType, rpInstanceId, rpPropertyId)
	if err != nil {
		log.Fatalf("error generating initial ReadProperty: %v\n", err)
	}

	replyRaw := make([]byte, 1024)
	sentRequests := 0
	for {
		if _, err := listenConn.WriteTo(mReadProperty, remoteUDPAddr); err != nil {
			log.Fatalf("Failed to write the request: %s\n", err)
		}

		log.Printf("sent: %x", mReadProperty)

		nBytes, remoteAddr, err := listenConn.ReadFrom(replyRaw)

		log.Printf("read %d bytes from %s: %x\n", nBytes, remoteAddr, replyRaw[:nBytes])

		serviceMsg, err := bacnet.Parse(replyRaw[:nBytes])
		if err != nil {
			log.Fatalf("error parsing the received message: %v\n", err)
		}

		cACKEnc, ok := serviceMsg.(*services.ComplexACK)
		if !ok {
			log.Fatalf("we didn't receive a CACK reply...\n")
		}

		log.Printf("unmarshalled BVLC: %#v\n", cACKEnc.BVLC)
		log.Printf("unmarshalled NPDU: %#v\n", cACKEnc.NPDU)

		decodedCACK, err := cACKEnc.Decode()
		if err != nil {
			log.Fatalf("couldn't decode the CACK reply: %v\n", err)
		}

		log.Printf(
			"decoded CACK reply:\n\tObject Type: %d\n\tInstance Id: %d\n\tProperty Id: %d\n\tValue: %f\n",
			decodedCACK.ObjectType, decodedCACK.InstanceId, decodedCACK.PropertyId, decodedCACK.PresentValue,
		)

		sentRequests++

		if sentRequests == rpN {
			break
		}

		time.Sleep(time.Duration(rpPeriod) * time.Second)
	}
}
