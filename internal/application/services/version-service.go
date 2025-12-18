package services

import (
	"github.com/ulfschneider/qrvc/internal/application/ports"
)

type VersionService struct {
	versionProvider ports.VersionProvider
}

func NewVersionService(versionProvider ports.VersionProvider) VersionService {
	return VersionService{versionProvider: versionProvider}
}

func (vs *VersionService) Version() string {
	return vs.versionProvider.Version()
}

func (vs *VersionService) Commit() string {
	return vs.versionProvider.Commit()
}

func (vs *VersionService) Time() string {
	return vs.versionProvider.Time()
}
