package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mahmud3253/Project/Api-Gateway/pkg/logger"
)

type JwtHandler struct {
	Sub       string
	Iss       string
	Exp       string
	Iat       string
	Aud       []string
	Role      string
	Token     string
	SigninKey string
	Log       logger.Logger
}

func (jwtHandler *JwtHandler) GenerateJWT() (access, refresh string, err error) {
	var (
		accessToken, refreshToken *jwt.Token
		claims                    jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)
	claims = accessToken.Claims.(jwt.MapClaims)
	claims["sub"] = jwtHandler.Sub
	claims["iss"] = jwtHandler.Iss
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	claims["aud"] = jwtHandler.Aud

	access, err = accessToken.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil {
		jwtHandler.Log.Error("error generating access token", logger.Error(err))
		return
	}

	refresh, err = refreshToken.SignedString([]byte(jwtHandler.SigninKey))
	if err != nil {
		jwtHandler.Log.Error("error generating access token", logger.Error(err))
		return
	}
	return access, refresh, nil
}

func (jwtHandler *JwtHandler) ExtractClaims() (jwt.MapClaims, error) {
	var(
		token *jwt.Token
		err error
	)

	token,err=jwt.Parse(jwtHandler.Token,func(t *jwt.Token)(interface{},error){
		return []byte(jwtHandler.SigninKey),nil
	})

	if err!=nil{
		return nil,err
	}

	claims,ok:=token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid){
		jwtHandler.Log.Error(("Invalid jwt token"))
		return nil,err
	}


	return claims,nil
}
