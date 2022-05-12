package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "practl",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var configFile string

	runCmd := &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			RunServer(configFile)
		},
	}

	runCmd.Flags().StringVarP(&configFile, "config", "c", "/etc/configs", "Config file for this run loop")
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Error running engine:", err.Error())
	}
}
