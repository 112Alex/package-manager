package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "pm",
    Short: "pm is a lightweight Go package manager for packing and deploying archives over SSH",
    Long:  `pm allows you to create archives defined by packet.json|yaml and update packages defined in packages.json|yaml via SSH.`,
}

// Execute is the entrypoint called from main.go
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "pm: %v\n", err)
        os.Exit(1)
    }
}
