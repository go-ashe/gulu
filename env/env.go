package env

type RunMode string

const (
	RunModeLocal RunMode = "local"
	RunModeDev   RunMode = "dev"
	RunModeTest  RunMode = "test"
	RunModeProd  RunMode = "prod"
)
