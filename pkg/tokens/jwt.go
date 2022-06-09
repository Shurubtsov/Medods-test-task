package tokens

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dshurubtsov/pkg/models"
)

type JWTMaker struct {
	SecretKey string
}

func (j JWTMaker) CreateToken(user models.User) (models.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["user_id"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwt := models.Token{}

	jwt.AccessToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return jwt, err
	}

	return j.CreateRefreshToken(jwt)
}

func (j JWTMaker) ValidateToken(accessToken string) (models.User, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(j.SecretKey), nil
	})

	user := models.User{}
	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user.Username = payload["user_id"].(string)

		return user, nil
	}

	return user, errors.New("invalid token")
}

func (j JWTMaker) ValidateRefreshToken(model models.Token) (models.User, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, j.SecretKey)

	user := models.User{}
	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return user, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return user, err
	}

	data, err := base64.URLEncoding.DecodeString(model.RefreshToken)
	if err != nil {
		return user, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return user, err
	}

	if string(plain) != model.AccessToken {
		return user, errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(model.AccessToken, claims)

	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid token")
	}

	user.Username = payload["user_id"].(string)

	return user, nil
}

func (j JWTMaker) CreateRefreshToken(token models.Token) (models.Token, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, j.SecretKey)

	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		fmt.Println(err.Error())

		return token, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return token, err
	}

	token.RefreshToken = base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token.AccessToken), nil))

	return token, nil
}
