package os

import "runtime"

func IsLinux() bool {
	return runtime.GOOS == `linux`
}

func IsWindows() bool {
	return runtime.GOOS == `windows`
}

func IsDarwin() bool {
	return runtime.GOOS == `darwin`
}
