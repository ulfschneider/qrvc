package services

import (
	"github.com/ulfschneider/qrvc/internal/application/ports"
)

type BomService struct {
	bomProvider ports.BomProvider
}

func NewBomService(bomProvider ports.BomProvider) BomService {
	return BomService{bomProvider: bomProvider}
}

func (bs *BomService) WriteBomJSON() error {
	err := bs.bomProvider.WriteBomJSON()
	return err
}
