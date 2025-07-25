package matcher

import (
    "path/filepath"
    "strings"
    "os"

    "github.com/bmatcuk/doublestar/v4"
    "github.com/example/package-manager/internal/config"
)

// Collect iterates over TargetSpec list and returns list of matching file paths (unique).
func Collect(baseDir string, targets []config.TargetSpec) ([]string, error) {
    seen := make(map[string]struct{})
    var out []string

    for _, t := range targets {
        var includePattern string
        var excludePattern string

        if t.Raw != "" {
            includePattern = t.Raw
        } else {
            p := t.Path
            if t.Pattern != "" {
                p = filepath.Join(p, t.Pattern)
            }
            includePattern = p
            excludePattern = t.Exclude
        }

        // Resolve relative to baseDir
        pattern := filepath.ToSlash(includePattern)
        fsys := os.DirFS(baseDir)
        matches, err := doublestar.Glob(fsys, pattern)
        if err != nil {
            return nil, err
        }
        for _, rel := range matches {
            // skip dirs (Glob returns file paths but double-check)
            if strings.HasSuffix(rel, "/") {
                continue
            }
            // Exclude pattern check
            if excludePattern != "" {
                ok, _ := doublestar.Match(excludePattern, rel)
                if ok {
                    continue
                }
            }
            abs := filepath.Join(baseDir, rel)
            if _, dup := seen[abs]; !dup {
                seen[abs] = struct{}{}
                out = append(out, abs)
            }
        }
    }
    return out, nil
}
