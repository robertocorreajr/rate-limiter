package storage

type Storage interface {
	Increment(key string) (int, error)
	TTL(key string) (int, error)
	SetTTL(key string, seconds int) error
	GetLimit(key string) (int, error)
	SetLimit(key string, value int) error
	Block(key string, seconds int) error
	IsBlocked(key string) (bool, error)
}
