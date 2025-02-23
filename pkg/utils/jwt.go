package utils

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (uint, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return 0, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint(claims["user_id"].(float64))
        return userID, nil
    }

    return 0, jwt.ErrSignatureInvalid
}
