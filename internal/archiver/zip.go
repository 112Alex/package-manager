package archiver

import (
    "github.com/mholt/archiver/v3"
)

// CreateZip archives given files (absolute paths) into outPath (zip).
func CreateZip(files []string, outPath string) error {
    return archiver.Archive(files, outPath)
}

// ExtractZip unarchives zip file to destDir.
func ExtractZip(archivePath, destDir string) error {
    return archiver.Unarchive(archivePath, destDir)
}
