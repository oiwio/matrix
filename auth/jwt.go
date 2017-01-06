package auth

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"matrix/config"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var configuration config.Config

func init() {
	configuration = config.New()
	// log.Infoln(configuration)
}

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(configuration.Authentication.TokenDuration)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = userId
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}
	return tokenString, nil
}

// func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
// 	if validity, ok := timestamp.(float64); ok {
// 		tm := time.Unix(int64(validity), 0)
// 		remainer := tm.Sub(time.Now())
// 		if remainer > 0 {
// 			return int(remainer.Seconds() + expireOffset)
// 		}
// 	}
// 	return expireOffset
// }

// func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
// 	redisConn := redis.Connect()
// 	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims["exp"]))
// }

// func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
// 	redisConn := redis.Connect()
// 	redisToken, _ := redisConn.GetValue(token)

// 	if redisToken == nil {
// 		return false
// 	}

// 	return true
// }

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(configuration.Authentication.PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(configuration.Authentication.PublicKeyPath)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
