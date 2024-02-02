//go:build !windows

package win_process

func RunMeElevated() {}
func AmAdmin() bool  { return true }
