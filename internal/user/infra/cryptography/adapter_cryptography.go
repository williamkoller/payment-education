package user_cryptography

import (
	portCryptography "github.com/williamkoller/system-education/internal/user/port/cryptography"
	"golang.org/x/crypto/bcrypt"
)

type BcryptAdapter struct {
	cost int
}

var _ portCryptography.Bcrypt = (*BcryptAdapter)(nil)

func NewBcryptHasher(cost int) *BcryptAdapter {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptAdapter{cost: cost}
}

func (b *BcryptAdapter) Hash(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), b.cost)
	return string(bytes), err
}

func (b *BcryptAdapter) HashComparer(plaintext string, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	if err != nil {
		return false, err
	}
	return true, nil
}
