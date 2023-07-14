package main

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/ulbios/bacnet"
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/services"
)

var (
	IAmCmd = &cobra.Command{
		Use:   "iam",
		Short: "Send IAm requests.",
		Long: "This example will wait until it receives a WhoIs request. Upon reception\n" +
			"it'll just reply with the configured IAm fields",
		Args: argValidation,
		Run:  IAmExample,
	}
)

func IAmExample(cmd *cobra.Command, args []string) {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %s", err)
	}

	ifaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("couldn't get interface information: %v\n", err)
	}

	listenConn, err := net.ListenPacket("udp", bAddr)
	if err != nil {
		log.Fatalf("failed to begin listening for packets: %v\n", err)
	}
	defer listenConn.Close()

	mIAm, err := bacnet.NewIAm(321, 31)
	if err != nil {
		log.Fatalf("error generating initial IAm: %v\n", err)
	}

	reqRaw := make([]byte, 1024)

	var nBytes int
	var remoteAddr net.Addr
	for {
		nBytes, remoteAddr, err = listenConn.ReadFrom(reqRaw)
		if err != nil {
			log.Fatalf("error reading incoming packet: %v\n", err)
		}

		if common.IsLocalAddr(ifaceAddrs, remoteAddr) {
			log.Printf("got our own broadcast, back to listening...\n")
			continue
		}

		log.Printf("read %d bytes from %s: %x\n", nBytes, remoteAddr, reqRaw[:nBytes])

		serviceMsg, err := bacnet.Parse(reqRaw[:nBytes])
		if err != nil {
			log.Fatalf("error parsing the received message: %v\n", err)
		}

		whoIsMessage, ok := serviceMsg.(*services.UnconfirmedWhoIs)
		if !ok {
			log.Fatalf("we didn't receive a WhoIs reply...\n")
		}

		log.Printf("received a WhoIs request!\n")

		log.Printf("\n\tunmarshalled WhoIs BVLC: %#v\n", whoIsMessage.BVLC)
		log.Printf("\n\tunmarshalled WhoIs NPDU: %#v\n", whoIsMessage.NPDU)
		log.Printf("\n\tunmarshalled WhoIs APDU: %#v\n", whoIsMessage.APDU)

		if _, err := listenConn.WriteTo(mIAm, remoteUDPAddr); err != nil {
			log.Fatalf("error sending our IAm response: %v\n", err)
		}
	}
}
