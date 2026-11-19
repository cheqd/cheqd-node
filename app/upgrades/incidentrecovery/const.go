package incidentrecovery

// UpgradeName is the exact name this upgrade must be proposed under via a
// software-upgrade governance proposal. It is matched against
// cfg.UpgradeName in Config.Validate (see config.go) so the config carried in
// plan.Info can never silently run under the wrong upgrade plan name.
const UpgradeName = "v4.3.0-incident-recovery"
