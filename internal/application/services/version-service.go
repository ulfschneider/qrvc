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

func (vs *VersionService) Version() (string, error) {
	return vs.versionProvider.Version()
}
