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
	ReadPropertyServerCmd = &cobra.Command{
		Use:   "rps",
		Short: "Reply ReadProperty requests with Complex ACKs.",
		Long: "This example will wait until it receives a ReadProperty request. Upon reception\n" +
			"it'll just reply with the configured Complex ACK fields",
		Args: argValidation,
		Run:  ReadPropertyServerExample,
	}
)

func ReadPropertyServerExample(cmd *cobra.Command, args []string) {
	listenConn, err := net.ListenPacket("udp", bAddr)
	if err != nil {
		log.Fatalf("failed to begin listening for packets: %v\n", err)
	}
	defer listenConn.Close()

	mCACKs := [][]byte{}
	for i := 0; i < 2; i++ {
		mCACK, err := bacnet.NewCACK(
			services.ServiceConfirmedReadProperty,
			objects.ObjectTypeAnalogOutput,
			uint32(i),
			objects.PropertyIdPresentValue,
			1.1*float32((i+1)),
		)
		if err != nil {
			log.Fatalf("error generating CACK %d: %v\n", i, err)
		}
		mCACKs = append(mCACKs, mCACK)
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

		readPropertyMessage, ok := serviceMsg.(*services.ConfirmedReadProperty)
		if !ok {
			log.Fatalf("we didn't receive a ReadProperty request: %x\n", reqRaw[:nBytes])
		}
		log.Printf("received a ReadProperty request!\n")

		decodedReadPropertyMessage, err := readPropertyMessage.Decode()
		if err != nil {
			log.Fatalf("error decoding the ReadProperty message: %v\n", err)
		}

		log.Printf("decoded ReadProperty message:\n\tObjectType: %d\n\tInstance ID: %d\n\tProperty ID: %d\n",
			decodedReadPropertyMessage.ObjectType, decodedReadPropertyMessage.InstanceId,
			decodedReadPropertyMessage.PropertyId)

		if decodedReadPropertyMessage.InstanceId >= uint32(len(mCACKs)) {
			bErr, err := bacnet.NewError(
				services.ServiceConfirmedReadProperty, objects.ErrorClassObject, objects.ErrorCodeUnknownObject)
			if err != nil {
				log.Fatalf("error generating Error reply: %v\n", err)
			}
			if _, err := listenConn.WriteTo(bErr, remoteAddr); err != nil {
				log.Fatalf("error sending our Error reply: %v\n", err)
			}
			log.Printf("we were asked for a wrong instance ID!\n")
			continue
		}

		if _, err := listenConn.WriteTo(mCACKs[decodedReadPropertyMessage.InstanceId], remoteAddr); err != nil {
			log.Fatalf("error sending our CACK reply: %v\n", err)
		}

		log.Printf("replied with our CACK!\n")
	}
}
