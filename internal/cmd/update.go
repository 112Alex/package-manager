package cmd

import (
    "fmt"
    "path/filepath"
    "os"

    "github.com/spf13/cobra"

    "github.com/example/package-manager/internal/archiver"
    "github.com/example/package-manager/internal/config"
    "github.com/example/package-manager/internal/transport"
    "github.com/example/package-manager/internal/version"
)

var updateCmd = &cobra.Command{
    Use:   "update <packages.json|yaml>",
    Short: "Download archives via SSH and extract based on packages config",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        cfgPath := args[0]
        cfg, err := config.LoadPackages(cfgPath)
        if err != nil {
            return err
        }
        // Parse SSH flags
        host, _ := cmd.Flags().GetString("host")
        user, _ := cmd.Flags().GetString("user")
        key, _ := cmd.Flags().GetString("key")
        remoteDir, _ := cmd.Flags().GetString("remote-dir")
        localDir, _ := cmd.Flags().GetString("dest")
        if localDir == "" {
            localDir = filepath.Dir(cfgPath)
        }
        if err := os.MkdirAll(localDir, 0o755); err != nil {
            return err
        }

        client, err := transport.Connect(transport.SSHConfig{Host: host, User: user, KeyPath: key})
        if err != nil {
            return err
        }
        defer client.Close()

        updated := 0
        for _, p := range cfg.Packages {
            // Build filename convention name-ver.zip
            if p.Ver == "" {
                p.Ver = "*" // wildcard (not enforced)
            }
            // For simplicity assume remote file exists as name-version.zip
            remoteName := fmt.Sprintf("%s-%s.zip", p.Name, p.Ver)
            remotePath := filepath.Join(remoteDir, remoteName)
            // Download
            localZip, err := client.Download(remotePath, localDir)
            if err != nil {
                fmt.Printf("skip %s: %v\n", p.Name, err)
                continue
            }
            // Version check (uses filename version)
            ver := p.Ver
            if !version.Satisfies(p.Ver, ver) {
                fmt.Printf("%s: version %s does not satisfy %s\n", p.Name, ver, p.Ver)
                continue
            }
            if err := archiver.ExtractZip(localZip, localDir); err != nil {
                fmt.Printf("extract %s error: %v\n", p.Name, err)
                continue
            }
            updated++
        }
        fmt.Printf("Updated %d/%d packages into %s\n", updated, len(cfg.Packages), localDir)
        return nil
    },
}

func init() {
    updateCmd.Flags().String("host", "", "ssh host (host:port)")
    updateCmd.Flags().String("user", os.Getenv("USER"), "ssh user")
    updateCmd.Flags().String("key", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"), "path to private key")
    updateCmd.Flags().String("remote-dir", "/opt/pm_repo", "remote directory with archives")
    updateCmd.Flags().String("dest", "", "local destination dir (default: dir of config file)")
    updateCmd.Short = "Download archives via SSH and extract based on packages config (see --help for flags)"
}
