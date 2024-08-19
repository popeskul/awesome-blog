package hash

//go:generate mockgen -destination=mocks/mock_hash_service.go -package=mockshash -source=hash_service.go github.com/popeskul/awesome-blog/backend/internal/domain/repository HashService

// HashService defines methods for hashing passwords.
type HashService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}
