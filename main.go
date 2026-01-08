package main

import (
	"os"

	"github.com/ibuilding-x/driver-box/pkg/driverbox"
	"github.com/ibuilding-x/driver-box/pkg/exports"
	"github.com/ibuilding-x/driver-box/pkg/plugins"
)

func main() {
	os.Setenv(ENV_VERGE_BASE_URL, "http://localhost:8080")
	plugins.RegisterAllPlugins()
	exports.LoadAllExports()
	driverbox.Exports.LoadExport(NewExport())
	driverbox.Start()
	select {}
}
