package environment

import "os"

const (
	gooseAction = "GOOSE_ACTION"
	gooseArgs   = "GOOSE_ARGS"
)

func GetGooseAction() string {
	return os.Getenv(gooseAction)
}

func GetGooseArgs() string {
	return os.Getenv(gooseArgs)
}
