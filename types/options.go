package types

type InstallOptions struct {
	RestartDiscord bool `json:"restartDiscord"`
}

type RepairOptions struct {
	DisablePlugins bool `json:"disablePlugins"`
}

type UninstallOptions struct {
	FullUninstall bool `json:"fullUninstall"`
}
