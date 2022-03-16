package jwtutils

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func Sign(subject string, claims map[string]interface{}, secret string) (string, error) {
	now := time.Now()
	registeredClaims := jwt.RegisteredClaims{
		Issuer:    "auth0",
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		extraClaims: map[string]interface{}{
			subject: claims,
		},
		RegisteredClaims: registeredClaims,
	})
	return token.SignedString([]byte(secret))
}

func ExtractToken(bearer string, secret string, mapper func(*jwt.Token) interface{}) (interface{}, error) {
	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return mapper(token), nil
}

type customClaims struct {
	extraClaims map[string]interface{}
	jwt.RegisteredClaims
}

func (cc customClaims) MarshalJSON() ([]byte, error) {
	claims, err := cc.mergeClaims()
	if err != nil {
		return nil, err
	}
	return json.Marshal(claims)
}

func (cc customClaims) Valid() error {
	return nil
}

func (cc *customClaims) mergeClaims() (map[string]json.RawMessage, error) {
	registeredClaimsJson, err := json.Marshal(cc.RegisteredClaims)
	if err != nil {
		return nil, err
	}
	var result map[string]json.RawMessage
	err = json.Unmarshal(registeredClaimsJson, &result)
	if err != nil {
		return nil, err
	}
	for name, val := range cc.extraClaims {
		value, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		result[name] = value
	}
	return result, nil
}
