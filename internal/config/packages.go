package config

// PackagesConfig corresponds to packages.json|yaml for update command.
// It lists required packages with optional version constraints.

type PackagesConfig struct {
    Packages []DepSpec `json:"packages" yaml:"packages"`
}
