package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add the different sub-commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(whoIsCmd)
	rootCmd.AddCommand(IAmCmd)
	rootCmd.AddCommand(ReadPropertyServerCmd)
	rootCmd.AddCommand(ReadPropertyClientCmd)
	rootCmd.AddCommand(WritePropertyServerCmd)
	rootCmd.AddCommand(WritePropertyClientCmd)

	rootCmd.PersistentFlags().StringVar(&rAddr, "remote-address", "127.0.0.1:47808", "Remote IP:Port tuple to connect to.")
	rootCmd.PersistentFlags().StringVar(&bAddr, "broadcast-address", ":47808", "Default broadcast address to bind to.")
}

var (
	rAddr string
	bAddr string

	commit string

	rootCmd = &cobra.Command{
		Use:   "bacnet-examples",
		Short: "A collection of BACnet/IP examples showcasing and verifying our implementation.",
		Long: "The idea is to execute a given example whilst capturoing traffic on WireShark. We\n" +
			"Are leveraging the `bacnet-stack` implementation as a compliant server to check we can\n" +
			"communicate and interoperate with other standard compliant devices. As stated above, bear\n" +
			"in mind we only support the BACnet/IP protocol leveraging UDP at the transport layer!",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get the examples associated git commit hash.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("built commit: %s\n", commit)
		},
	}
)

func argValidation(cmd *cobra.Command, args []string) error {
	if false {
		return fmt.Errorf("false is true now?")
	}

	return nil
}

func execute() {
	log.SetFlags(log.Lshortfile)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	execute()
}
