package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Keyfile []byte

type Claims struct {
	Subject string
	Kind    string
}

func (k *Keyfile) New(id string, kind string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    "PongleHub",
			Audience:  jwt.ClaimStrings{kind},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	)

	tokenString, err := token.SignedString(k)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %+v", err)
	}

	return tokenString, nil
}

func (k *Keyfile) Parse(token string) (Claims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return k, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return Claims{}, fmt.Errorf("couldn't parse non-token object: %s", token)
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return Claims{}, fmt.Errorf("invalid signature: %s", t.Signature)
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return Claims{}, errors.New("token expired or not ready yet")
		} else {
			return Claims{}, fmt.Errorf("error parsing token: %+v", err)
		}
	}

	issuer, err := t.Claims.GetIssuer()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting issuer: %+v", err)
	}

	if issuer != "PongleHub" {
		return Claims{}, fmt.Errorf("invalid issuer, expected PongleHub, got %s", issuer)
	}

	audience, err := t.Claims.GetAudience()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting audience: %+v", err)
	}

	if len(audience) != 1 {
		return Claims{}, fmt.Errorf("expected one audience, got %d", len(audience))
	}

	subject, err := t.Claims.GetSubject()
	if err != nil {
		return Claims{}, fmt.Errorf("error getting subject: %+v", err)
	}

	return Claims{
		Subject: subject,
		Kind:    audience[0],
	}, nil
}
