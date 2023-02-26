package service

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/pkg/errors"
)

// RegisterClient registers client
func (svc *Service) RegisterClient(
	ctx context.Context, clientToReg *model.ClientToRegister,
) (*model.Client, error) {
	client := &model.Client{
		Name:   clientToReg.Name,
		Email:  clientToReg.Email,
		Phone:  clientToReg.Phone,
		Age:    clientToReg.Age,
		Weight: clientToReg.Weight,
		Height: clientToReg.Height,
		Gender: clientToReg.Gender,
	}

	password, err := model.HashPassword(clientToReg.Password)
	if err != nil {
		return nil, errors.Wrap(err, "model.HashPassword")
	}

	authData := &model.AuthData{
		SourceType: model.AuthDataSourceTypeClient,
		Password:   password,
	}

	txErr := svc.storage.WithTransaction(ctx, func(tx store.Tx) error {
		err = svc.storage.CreateClientTx(ctx, tx, client)
		if err != nil {
			return errors.Wrap(err, "svs.storage.CreateClientTx")
		}

		authData.SourceID = client.ID
		err = svc.storage.CreateAuthDataTx(ctx, tx, authData)
		if err != nil {
			return errors.Wrap(err, "svs.storage.CreateAuthDataTx")
		}

		return nil
	})
	if txErr != nil {
		return nil, errors.Wrap(txErr, "svs.storage.WithTransaction")
	}

	return client, nil
}
