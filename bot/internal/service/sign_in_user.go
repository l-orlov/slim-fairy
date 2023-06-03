package service

import (
	"context"

	"github.com/l-orlov/slim-fairy/bot/internal/jwttoken"
	model2 "github.com/l-orlov/slim-fairy/bot/internal/model"
	"github.com/pkg/errors"
)

// SignInUser do user sign-in
func (svc *Service) SignInUser(
	ctx context.Context, userToSignIn *model2.UserToSignIn,
) (token string, err error) {
	user, err := svc.storage.GetUserByEmail(ctx, userToSignIn.Email)
	if err != nil {
		return "", errors.Wrap(err, "svs.storage.GetUserByEmail")
	}

	sourceType := model2.AuthDataSourceTypeUser
	authData, err := svc.storage.GetAuthDataBySourceIDAndType(ctx, user.ID, sourceType)
	if err != nil {
		return "", errors.Wrap(err, "svs.storage.GetAuthDataBySourceIDAndType")
	}

	if !model2.CheckPasswordHash(authData.Password, userToSignIn.Password) {
		return "", ErrWrongPassword
	}

	token, err = jwttoken.New(user.ID, sourceType)
	if err != nil {
		return "", errors.Wrap(err, "jwttoken.New")
	}

	return token, nil
}
