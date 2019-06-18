/*
 * @Author: qiuling
 * @Date: 2019-06-17 19:32:16
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-18 11:16:55
 */

package biz

import (
	. "artifact/pkg/config"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func Jwt2Token(tokenString string) (token string, err error) {

	jwtToken, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(Config.Jwt.Secret), nil
	})

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		// fmt.Printf("claims:%+v\n", claims)
		// fmt.Println(claims["artifact"], claims["token"])
		if claims["artifact"].(string) != Config.Jwt.Uid {
			err = errors.New("ERR_INVALID_TOKEN")
			return
		}
		token = claims["token"].(string)
		return
	} else {
		fmt.Println(err)
		err = errors.New("ERR_INVALID_TOKEN")
		return
	}
}
