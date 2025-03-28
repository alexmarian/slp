package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"

	tests := []struct {
		name           string
		makePassword   string
		verifyPassword string
		expiration     time.Duration
		wantErr        bool
	}{
		{name: "Correct password",
			makePassword:   password1,
			verifyPassword: password1,
			expiration:     1 * time.Hour,
			wantErr:        false,
		},
		{name: "Wrong  password",
			makePassword:   password1,
			verifyPassword: password2,
			expiration:     1 * time.Hour,
			wantErr:        true,
		},
		{name: "Expiration password",
			makePassword:   password1,
			verifyPassword: password1,
			expiration:     1 * time.Millisecond,
			wantErr:        true,
		},
	}
	userId := uuid.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt, err := MakeJWT(userId, tt.makePassword, tt.expiration)
			if err != nil {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(2 * time.Millisecond)
			_, err = ValidateJWT(jwt, tt.verifyPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestGetBearerToken(t *testing.T) {

	tests := []struct {
		name    string
		headers map[string][]string
		wantErr bool
	}{
		{name: "Empty headers",
			headers: make(map[string][]string, 0),
			wantErr: true,
		},
		{name: "No Authorization headers",
			headers: map[string][]string{"Bla1": {"1"}, "Bla2": {"1"}},
			wantErr: true,
		},
		{name: "Empty Authorization headers",
			headers: map[string][]string{"Authorization": {}, "Bla": {"2"}},
			wantErr: true,
		},
		{name: "No bearer Authorization headers",
			headers: map[string][]string{"Authorization": {"12312"}},
			wantErr: true,
		},
		{name: "Nice authorization headers",
			headers: map[string][]string{"Authorization": {"Bearer 12312"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
