package config

import (
    "encoding/json"
    "fmt"
)

// TargetSpec represents an include pattern with optional exclusion.
// It supports two JSON/YAML forms:
// 1. String: "./dir/*.txt"
// 2. Object: {"path": "./dir", "exclude": "*.tmp"}

type TargetSpec struct {
    Path    string `json:"path" yaml:"path"`
    Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"` // optional separate pattern
    Exclude string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
    Raw  string `json:"-" yaml:"-"` // holds raw string form when provided
}

// UnmarshalJSON implements custom json unmarshalling to allow string or object.
func (t *TargetSpec) UnmarshalJSON(b []byte) error {
    // Check if it is quoted string
    if len(b) > 0 && b[0] == '"' {
        var s string
        if err := json.Unmarshal(b, &s); err != nil {
            return err
        }
        t.Raw = s
        return nil
    }
    // Otherwise decode as struct
    type alias TargetSpec
    var tmp alias
    if err := json.Unmarshal(b, &tmp); err != nil {
        return err
    }
    *t = TargetSpec(tmp)
    return nil
}

// String returns a human-readable value.
func (t TargetSpec) String() string {
    if t.Raw != "" {
        return t.Raw
    }
    return fmt.Sprintf("%s (exclude=%s)", t.Path, t.Exclude)
}

// DepSpec represents dependency on another packet with optional semver constraint.
type DepSpec struct {
    Name string `json:"name" yaml:"name"`
    Ver  string `json:"ver,omitempty" yaml:"ver,omitempty"`
}

// PacketConfig corresponds to packet.json|yaml for create command.
type PacketConfig struct {
    Name    string      `json:"name" yaml:"name"`
    Ver     string      `json:"ver" yaml:"ver"`
    Targets []TargetSpec `json:"targets" yaml:"targets"`
    Packets []DepSpec   `json:"packets,omitempty" yaml:"packets,omitempty"`
}
