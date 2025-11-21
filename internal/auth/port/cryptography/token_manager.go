package port_auth_cryptography

type TokenManager interface {
	Sign(data map[string]interface{}) (string, error)
	Verify(tokenStr string) (map[string]interface{}, error)
}
