package utils

import (
	"fmt"
	"runtime/debug"
)

func SetClientVersion() string {
	version := "development"
	build, ok := debug.ReadBuildInfo()
	if ok {
		if build.Main.Version != "" {
			version = build.Main.Version
		}
	}
	return version
}

func SetUserAgent(ua string) string {
	if ua == "" {
		return fmt.Sprintf("TinesGoSdk/%s", SetClientVersion())
	}
	return ua
}
