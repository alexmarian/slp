package auth

import (
	"log"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	wpassword := "wpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}
	log.Println(hash)
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("Error checking password/expect match: %s", err)
	}
	err = CheckPasswordHash(hash, wpassword)
	if err == nil {
		t.Errorf("Error checking password/expect no match: %s", err)
	}
}
