package services

import (
	"github.com/ulfschneider/qrvc/internal/application/config"
	"github.com/ulfschneider/qrvc/internal/application/ports"
)

type QRCardService struct {
	settings config.Settings
	repo     ports.Repository
	editor   ports.VCardEditor
}

func NewQRCardService(settings config.Settings, repo ports.Repository, editor ports.VCardEditor) QRCardService {

	return QRCardService{
		settings: settings,
		repo:     repo,
		editor:   editor,
	}
}

func (qs *QRCardService) TransformCard() error {

	card, err := qs.repo.ReadOrCreateVCard()
	if err != nil {
		return err
	}

	if qs.settings.Silent == false {
		if err = qs.editor.Edit(card); err != nil {
			return err
		}
	}

	if err = qs.repo.WriteVCard(card); err != nil {
		return err
	}

	if err = qs.repo.WriteQRCode(card); err != nil {
		return err
	}

	return nil
}
