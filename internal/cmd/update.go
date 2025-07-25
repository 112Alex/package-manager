package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/example/package-manager/internal/config"
)

var updateCmd = &cobra.Command{
    Use:   "update <packages.json|yaml>",
    Short: "Download archives via SSH and extract based on packages config (stub)",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        cfgPath := args[0]
        cfg, err := config.LoadPackages(cfgPath)
        if err != nil {
            return err
        }
        // TODO: implement SSH download and extraction
        fmt.Printf("Would update %d packages (not yet implemented)\n", len(cfg.Packages))
        return nil
    },
}
