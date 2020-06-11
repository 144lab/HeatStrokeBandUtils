package frontend

import "github.com/nobonobo/wecty"

type Inform struct {
	wecty.Core
	RawSize          int
	RriSize          int
	EnvSize          int
	FirmwareRevision string
	LastRri          Rri
	LastEnv          Env
}
