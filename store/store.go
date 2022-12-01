package store

type SecretStore interface {
	NewPassword(password, name, service string) error
	NewSecret(kind, secret, name, service string) error
	Password(name, service string) ([]byte, error)
	Secret(kind, name, service string) ([]byte, error)
}
