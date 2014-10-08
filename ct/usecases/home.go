package usecases

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

type HomeInteractor struct {
	ClientRepository domain.ClientRepository
	Logger           Logger
}

func NewHomeInteractor(clientRepository domain.ClientRepository, logger Logger) *HomeInteractor {
	return &HomeInteractor{
		ClientRepository: clientRepository,
		Logger:           logger,
	}
}

func (hi *HomeInteractor) GetClient(clientId int64) (*domain.Client, error) {
	client, err := hi.ClientRepository.FindById(clientId)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, ErrInvalidClientId
	}
	return client, err
}
