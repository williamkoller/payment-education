package port_cryptography

type Bcrypt interface {
	Hash(plaintext string) (string, error)
	HashComparer(plaintext string, hashed string) bool
}
