/*
* @Author: qiuling
* @Date: 2019-06-17 19:32:16
* @Last Modified by: qiuling
* @Last Modified time: 2020-07-03 18:26:57
 */

package biz

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	. "github.com/wlxpkg/base"
	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	Uid  string `json:"uid"`
	Name string `json:"cj"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	Uid          string `json:"uid"`
	Name         string `json:"cj"`
	RefreshToken string `json:"token"`
	jwt.StandardClaims
}

type Jwt struct {
	JwtStr       string `json:"jwt"`
	Refresh      string `json:"refresh"`
	RefreshToken string `json:"refresh_token"`
}

func CreateJwt(user_id string) (jwtData Jwt, err error) {
	issuer := Config.Jwt.Iss
	priv, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(Config.Rsa.Private))
	if err != nil {
		log.Err(err)
		return
	}

	uid := []byte(user_id)
	uidRsa, _ := rsaEncrypt(uid)

	// base64 decode
	uidRsaStr := base64.StdEncoding.EncodeToString(uidRsa)

	// Create the Claims
	claims := JwtClaims{
		uidRsaStr,
		Config.Jwt.Uid,
		jwt.StandardClaims{
			Id:        RandStr(32),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), // 12小时有效期
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtData.JwtStr, err = token.SignedString(priv)

	if err != nil {
		log.Err(err)
		return
	}

	jwtData.RefreshToken = RandStr(32)

	refreshClaims := RefreshClaims{
		uidRsaStr,
		Config.Jwt.Uid,
		jwtData.RefreshToken,
		jwt.StandardClaims{
			Id:        RandStr(32),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), // 30天有效期
			Issuer:    issuer,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	jwtData.Refresh, err = refreshToken.SignedString(priv)
	return
}

func Jwt2Token(tokenString string) (uid string, err error) {

	issuer := Config.Jwt.Iss

	/* jwtToken, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		  // Don't forget to validate the alg is what you expect:
		  if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			  return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
		  }
		  // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		  return []byte(Config.Jwt.Secret), nil
		  // return *rsa.PrivateKey, nil
	  }) */

	jwtToken, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		publicKey := Config.Rsa.Public
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		// return *rsa.PrivateKey, nil
	})

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		fmt.Printf("claims:%#v\n", claims)
		// fmt.Println(claims["zwyd"], claims["token"])
		if claims[issuer] == nil || claims[issuer].(string) != Config.Jwt.Uid {
			err = errors.New("ERR_INVALID_TOKEN")
			return
		}

		uid = ""
		if _, ok := claims["uid"]; ok {
			uid = claims["uid"].(string)
		}

		// base64 decode
		var data []byte
		data, err = base64.StdEncoding.DecodeString(uid)
		if err != nil {
			log.Err(err)
			err = errors.New("ERR_INVALID_TOKEN")
			return
		}

		// RSA 解密
		var user_id []byte
		user_id, err = rsaDecrypt(data)
		if err != nil {
			log.Err(err)
			err = errors.New("ERR_INVALID_TOKEN")
			return
		}

		uid = string(user_id)
		return
	} else {
		fmt.Println(err)
		err = errors.New("ERR_INVALID_TOKEN")
		return
	}
}

// rsaDecrypt rsa 解密
func rsaDecrypt(ciphertext []byte) ([]byte, error) {
	priv, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(Config.Rsa.Private))
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// rsaEncrypt rsa 加密
func rsaEncrypt(ciphertext []byte) ([]byte, error) {
	pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(Config.Rsa.Public))
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pub, ciphertext)
}
