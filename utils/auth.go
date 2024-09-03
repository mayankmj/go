package utils

import (
    "time"
    "blog-app/models"
    "github.com/golang-jwt/jwt/v5"
    "errors"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

var secretKey = []byte("UrL4/zwWqS7zEwRrNZK3Yp1KCnW+dnTXg1o6jaLBCCw=")
func GenerateJWT(user *models.User) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: user.Username,
        Email:    user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, err
    }

    return claims, nil
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return secretKey, nil
    })
    if err != nil {
        return "", err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", errors.New("invalid token")
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", errors.New("user_id not found in token")
    }

    return userID, nil
}

