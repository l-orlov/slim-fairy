package service

import (
	"context"

	"github.com/l-orlov/slim-fairy/internal/jwttoken"
	"github.com/l-orlov/slim-fairy/internal/model"
	"github.com/pkg/errors"
)

// SignInClient do client sign-in
func (svc *Service) SignInClient(
	ctx context.Context, clientToSignIn *model.ClientToSignIn,
) (token string, err error) {
	client, err := svc.storage.GetClientByEmail(ctx, clientToSignIn.Email)
	if err != nil {
		return "", errors.Wrap(err, "svs.storage.GetClientByEmail")
	}

	sourceType := model.AuthDataSourceTypeClient
	authData, err := svc.storage.GetAuthDataBySourceIDAndType(ctx, client.ID, sourceType)
	if err != nil {
		return "", errors.Wrap(err, "svs.storage.GetAuthDataBySourceIDAndType")
	}

	hashedPassword, err := model.HashPassword(clientToSignIn.Password)
	if err != nil {
		return "", errors.Wrap(err, "model.HashPassword")
	}

	if hashedPassword != authData.Password {
		return "", ErrWrongPassword
	}

	token, err = jwttoken.New(client.ID, sourceType)
	if err != nil {
		return "", errors.Wrap(err, "jwttoken.New")
	}

	return token, nil
}
