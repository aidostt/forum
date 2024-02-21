package main

import (
	"context"
	"fmt"
	"forum.aidostt-buzuk/internal/data"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgtype"
	"math/rand"
	"time"
)

var signingKey = "akosjfq[wo0r-0skldfj409"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

//TODO: create const TTL variables for accessToken and refreshToken, secret key

func (app *application) NewSession(ctx context.Context, user *data.User) (res Tokens, err error) {
	res.AccessToken, err = newAccessToken(user.ID, time.Minute*10)
	if err != nil {
		return Tokens{}, err
	}
	res.RefreshToken, err = newRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	session := data.Session{
		RefreshToken: res.RefreshToken,
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 30),
	}
	err = app.models.Tokens.SetSession(user, session)
	return
}

func newAccessToken(userID pgtype.UUID, ttl time.Duration) (string, error) {
	//TODO: convert uuid to string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(signingKey))
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
