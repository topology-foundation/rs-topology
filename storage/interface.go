package storage

type Storage interface {
	Has(key []byte) (bool, error)
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) error
	Close() error
}
