package crypto

// Service holds the business rules. It knows the store, never the transport.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}
