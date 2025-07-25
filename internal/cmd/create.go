package cmd

import (
	"fmt"
	"path/filepath"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/package-manager/internal/archiver"
	"github.com/example/package-manager/internal/config"
	"github.com/example/package-manager/internal/matcher"
	"github.com/example/package-manager/internal/transport"
)

var createCmd = &cobra.Command{
	Use:   "create <packet.json|yaml>",
	Short: "Create archive from packet config and optionally upload via SSH",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := args[0]
		cfg, err := config.LoadPacket(cfgPath)
		if err != nil {
			return err
		}
		baseDir := filepath.Dir(cfgPath)
		files, err := matcher.Collect(baseDir, cfg.Targets)
		if err != nil {
			return err
		}
		outName := fmt.Sprintf("%s-%s.zip", cfg.Name, cfg.Ver)
		outPath := filepath.Join(baseDir, outName)
		if err := archiver.CreateZip(files, outPath); err != nil {
			return err
		}

		// optional upload
		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")
		remoteDir, _ := cmd.Flags().GetString("remote-dir")
		if host != "" {
			client, err := transport.Connect(transport.SSHConfig{Host: host, User: user, KeyPath: key})
			if err != nil {
				return err
			}
			defer client.Close()
			if err := client.Upload(outPath, remoteDir); err != nil {
				return err
			}
			fmt.Printf("Uploaded %s to %s:%s\n", outName, host, remoteDir)
		} else {
			fmt.Printf("Created %s with %d files (not uploaded)\n", outPath, len(files))
		}
		return nil
	},
}

func init() {
	createCmd.Flags().String("host", "", "ssh host (host:port) for upload")
	createCmd.Flags().String("user", os.Getenv("USER"), "ssh user for upload")
	createCmd.Flags().String("key", filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"), "path to private key for upload")
	createCmd.Flags().String("remote-dir", "/opt/pm_repo", "remote directory to upload archive")
}
