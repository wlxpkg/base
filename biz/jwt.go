/*
 * @Author: qiuling
 * @Date: 2019-06-17 19:32:16
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-19 16:25:03
 */

package biz

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"

	"github.com/dgrijalva/jwt-go"
)

func Jwt2Token(tokenString string) (uid string, err error) {

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

	// R(jwtToken, "jwtToken")

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		// fmt.Printf("claims:%#v\n", claims)
		// fmt.Println(claims["zwyd"], claims["token"])
		if claims["zwyd"] == nil || claims["zwyd"].(string) != Config.Jwt.Uid {
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
