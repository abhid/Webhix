package core

import (
	"runtime"

	"github.com/GaIsBAX/Webhix/internal/domain"
)

var (
	WebhixVersion = "unknown"
	Commit        = "unknown"
	Built         = "unknown"
)

type Version struct{}

func NewVersion() *Version {
	return &Version{}
}

func (v *Version) Info() domain.VersionInfo {
	return domain.VersionInfo{
		Version: WebhixVersion,
		Commit:  Commit,
		Built:   Built,
		Go:      runtime.Version(),
	}
}
