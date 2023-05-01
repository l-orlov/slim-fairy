package service

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/l-orlov/slim-fairy/internal/store"
	"github.com/pkg/errors"
)

// RegisterUser registers user
func (svc *Service) RegisterUser(
	ctx context.Context, userToReg *model.UserToRegister,
) (*model.User, error) {
	user := &model.User{
		Name:   userToReg.Name,
		Email:  userToReg.Email,
		Phone:  userToReg.Phone,
		Age:    userToReg.Age,
		Weight: userToReg.Weight,
		Height: userToReg.Height,
		Gender: userToReg.Gender,
	}

	password, err := model.HashPassword(userToReg.Password)
	if err != nil {
		return nil, errors.Wrap(err, "model.HashPassword")
	}

	authData := &model.AuthData{
		SourceType: model.AuthDataSourceTypeUser,
		Password:   password,
	}

	txErr := svc.storage.WithTransaction(ctx, func(tx store.Tx) error {
		err = svc.storage.CreateUserTx(ctx, tx, user)
		if err != nil {
			return errors.Wrap(err, "svs.storage.CreateUserTx")
		}

		authData.SourceID = user.ID
		err = svc.storage.CreateAuthDataTx(ctx, tx, authData)
		if err != nil {
			return errors.Wrap(err, "svs.storage.CreateAuthDataTx")
		}

		return nil
	})
	if txErr != nil {
		return nil, errors.Wrap(txErr, "svs.storage.WithTransaction")
	}

	return user, nil
}
