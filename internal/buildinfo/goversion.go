package buildinfo

import "runtime"

func goVersion() string {
	return runtime.Version()
}
