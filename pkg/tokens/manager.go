package tokens

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dshurubtsov/pkg/models"
	"github.com/golang-jwt/jwt/v4"
)

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userId string) (string, error) {
	// time when token will be expired
	timeExp := jwt.NewNumericDate(time.Now().Add(time.Hour * 2))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: timeExp,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseJWT(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) ValidateRefreshToken(tokens models.Token) (bool, error) {
	sha512 := sha512.New()
	io.WriteString(sha512, m.signingKey)

	salt := string(sha512.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return false, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}

	data, err := base64.URLEncoding.DecodeString(tokens.RefreshToken)
	if err != nil {
		return false, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, err
	}

	if string(plain) != tokens.AccessToken {
		return false, errors.New("invalid tokens")
	}

	return true, nil
}

func (m *Manager) NewRefreshToken(accessToken string) (string, error) {
	sha512 := sha512.New()
	io.WriteString(sha512, m.signingKey)

	salt := string(sha512.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	refreshToken := base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(accessToken), nil))

	return refreshToken, nil
}
