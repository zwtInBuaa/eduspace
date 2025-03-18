package utils

import (
	"EDU_TH_2_backend/gin/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenUtil interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type tokenUtil struct {
	secretKey     []byte
	expirationSec int64
}

func NewTokenUtil(secretKey string, expirationSec int64) TokenUtil {
	// secretKey:用于签署JWT令牌的密钥，expirationSec:JWT令牌的过期时间（以秒为单位）
	return &tokenUtil{secretKey: []byte(secretKey), expirationSec: expirationSec}
}

func (tu *tokenUtil) GenerateToken(user *models.User) (string, error) {
	var role_string string
	if user.Role == 0 {
		role_string = "管理员"
	} else if user.Role == 1 {
		role_string = "老师"
	} else if user.Role == 2 {
		role_string = "助教"
	} else if user.Role == 3 {
		role_string = "学生"
	}

	// 使用HS256算法对JWT进行签名，生成一个新的token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // 定义JWT载荷
		"id":   user.ID,     // 将用户id作为token的一个声明（claim）
		"role": role_string, // 将用户姓名作为token的一个声明
		//"email":    user.Email,                                                           // 将用户email作为token的一个声明
		"exp": time.Now().Add(time.Duration(tu.expirationSec) * time.Second).Unix(), // 设置token的过期时间
	})
	// 使用secret key对token进行签名并转换为字符串格式
	tokenString, err := token.SignedString(tu.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

var invalidTokens = make(map[string]bool) // 全局变量，用于存储失效的 token

func InvalidateToken(tokenString string) error {
	// 使 token 失效
	invalidTokens[tokenString] = true
	return nil
}

func (tu *tokenUtil) ValidateToken(tokenString string) (*jwt.Token, error) {
	if invalidTokens[tokenString] {
		return nil, fmt.Errorf("Token 已失效")
	}

	// 解析token，并校验签名是否正确
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法是否正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回secret key用于校验签名
		return tu.secretKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("Token 已过期")
			}
		}
		return nil, err
	}
	return token, nil
}
