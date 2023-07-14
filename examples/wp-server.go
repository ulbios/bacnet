package main

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/ulbios/bacnet"
	"github.com/ulbios/bacnet/objects"
	"github.com/ulbios/bacnet/services"
)

var (
	WritePropertyServerCmd = &cobra.Command{
		Use:   "wps",
		Short: "Reply WriteProperty requests with Simple ACKs.",
		Long: "This example will wait until it receives a WriteProperty request. Upon reception\n" +
			"it'll store the value and reply with a Simple ACK.",
		Args: argValidation,
		Run:  WritePropertyServerExample,
	}
)

func WritePropertyServerExample(cmd *cobra.Command, args []string) {
	remoteUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %s", err)
	}

	listenConn, err := net.ListenPacket("udp", bAddr)
	if err != nil {
		log.Fatalf("failed to begin listening for packets: %v\n", err)
	}
	defer listenConn.Close()

	storedValues := []float32{0, 0}

	sACK, err := bacnet.NewSACK(services.ServiceConfirmedWriteProperty)
	if err != nil {
		log.Fatalf("error generating the SACK: %v\n", err)
	}

	iAm, err := bacnet.NewIAm(321, 31)
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

		log.Printf("read %d bytes from %s: %x\n", nBytes, remoteAddr, reqRaw[:nBytes])

		serviceMsg, err := bacnet.Parse(reqRaw[:nBytes])
		if err != nil {
			log.Fatalf("error parsing the received message: %v\n", err)
		}

		writePropertyMessage, ok := serviceMsg.(*services.ConfirmedWriteProperty)
		if !ok {
			_, ok := serviceMsg.(*services.UnconfirmedWhoIs)
			if !ok {
				log.Printf("we didn't receive a WriteProperty request! Back to listening...\n")
				continue
			}
			log.Printf("received a WhoIs request!\n")
			if _, err := listenConn.WriteTo(iAm, remoteUDPAddr); err != nil {
				log.Fatalf("error sending our IAm reply: %v\n", err)
			}
			continue
		}
		log.Printf("received a WriteProperty request!\n")

		decodedWritePropertyMessage, err := writePropertyMessage.Decode()
		if err != nil {
			log.Fatalf("error decoding the WriteProperty message: %v\n", err)
		}

		log.Printf(
			"decoded WriteProperty message:\n\tObjectType: %d\n\tInstance ID: %d\n\tProperty ID: %d\n\tValue: %f\n",
			decodedWritePropertyMessage.ObjectType, decodedWritePropertyMessage.InstanceId,
			decodedWritePropertyMessage.PropertyId, decodedWritePropertyMessage.Value)

		if decodedWritePropertyMessage.InstanceId >= uint32(len(storedValues)) {
			bErr, err := bacnet.NewError(
				services.ServiceConfirmedWriteProperty, objects.ErrorClassObject, objects.ErrorCodeUnknownObject)
			if err != nil {
				log.Fatalf("error generating Error reply: %v\n", err)
			}
			if _, err := listenConn.WriteTo(bErr, remoteAddr); err != nil {
				log.Fatalf("error sending our Error reply: %v\n", err)
			}
			log.Printf("we were asked for a wrong instance ID!\n")
			continue
		}

		storedValues[decodedWritePropertyMessage.InstanceId] = decodedWritePropertyMessage.Value

		if _, err := listenConn.WriteTo(sACK, remoteAddr); err != nil {
			log.Fatalf("error sending our CACK reply: %v\n", err)
		}
		log.Printf("replied with our SACK!\n")
	}
}
