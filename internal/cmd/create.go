package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/example/package-manager/internal/archiver"
	"github.com/example/package-manager/internal/config"
	"github.com/example/package-manager/internal/matcher"
)

var createCmd = &cobra.Command{
	Use:   "create <packet.json|yaml>",
	Short: "Create archive from packet config and (later) upload via SSH",
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
		fmt.Printf("Created %s with %d files\n", outPath, len(files))
		// TODO: upload via SSH
		return nil
	},
}
