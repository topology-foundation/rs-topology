package storage

type Storage interface {
	Set(key, value []byte) error
}
