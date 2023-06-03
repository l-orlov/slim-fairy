package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/l-orlov/slim-fairy/bot/internal/config"
	"github.com/l-orlov/slim-fairy/bot/internal/model"
	"github.com/pkg/errors"
)

// Keys for token claims map
const (
	// Standard claims
	claimKeyExpiresAt = "exp" // Expires at
	claimKeyId        = "jti" // Token id
	claimKeyIssuedAt  = "iat" // Issued at
	claimKeyIssuer    = "iss" // Issuer
	claimKeyNotBefore = "nbf" // Not before
	claimKeySubject   = "sub" // Subject
	// Custom claims
	claimKeySourceType = "type" // Claim for source_type
)

// Claims has key fields from token
type Claims struct {
	SourceID   uuid.UUID
	SourceType model.AuthDataSourceType
}

// New returns new JWT token for source_id, source_type
func New(sourceID uuid.UUID, sourceType model.AuthDataSourceType) (string, error) {
	jwtCfg := config.Get().JWTToken

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimKeyExpiresAt:  time.Now().Add(jwtCfg.Lifitime).Unix(),
		claimKeyId:         uuid.NewString(),
		claimKeyIssuedAt:   time.Now().Unix(),
		claimKeyIssuer:     jwtCfg.Issuer,
		claimKeyNotBefore:  time.Now().Unix(),
		claimKeySubject:    sourceID.String(),
		claimKeySourceType: sourceType,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return "", errors.Wrap(err, "token.SignedString")
	}

	return tokenString, nil
}

// ParseAndValidate parses and validates JWT token
func ParseAndValidate(tokenString string) (*Claims, error) {
	jwtCfg := config.Get().JWTToken

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return secret key
		return []byte(jwtCfg.Secret), nil
	}, jwt.WithIssuer(jwtCfg.Issuer))
	if err != nil {
		return nil, errors.Wrap(err, "jwt.Parse")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.Errorf("not valid token claims")
	}

	// Get key fields from token
	sourceIDStr := claims[claimKeySubject].(string)
	sourceID, err := uuid.Parse(sourceIDStr)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse")
	}

	sourceType := claims[claimKeySourceType].(string)

	return &Claims{
		SourceID:   sourceID,
		SourceType: model.AuthDataSourceType(sourceType),
	}, nil
}
