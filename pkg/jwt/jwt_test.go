package jwt_test

import (
	"api/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "test@test.com"
	jwtService := jwt.NewJWT("y9aAOdn0O6ZjMU-Wrdzfaem2d7ZeTYTl-RwNNb3jemw")
	token, err := jwtService.Create(&jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Error(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is not valid")
	}
	if data.Email != email {
		t.Fatalf("Email %s not equal %s", data.Email, email)
	}
}
