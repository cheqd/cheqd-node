package incidentrecovery

// UpgradeName is the exact name this upgrade must be proposed under via a
// software-upgrade governance proposal. It is matched against
// cfg.UpgradeName in Config.Validate (see config.go) so the loaded config
// can never silently run under the wrong upgrade plan name.
const UpgradeName = "v4.3.0-incident-recovery"

// ConfigPathEnvVar names the environment variable used to locate the
// incident-recovery config YAML file at upgrade-execution time. Deliberately
// read lazily, inside the registered upgrade handler itself (see app.go's
// RegisterUpgradeHandlers wiring), rather than at app construction time —
// this upgrade handler only needs to run once, on the small number of nodes
// executing this specific upgrade plan; a missing or invalid config file
// must never prevent the binary from starting or running any other,
// unrelated upgrade.
const ConfigPathEnvVar = "CHEQD_INCIDENT_RECOVERY_CONFIG"

// DefaultConfigPath is used if ConfigPathEnvVar is unset.
const DefaultConfigPath = "/etc/cheqd-incident-recovery/config.yaml"
