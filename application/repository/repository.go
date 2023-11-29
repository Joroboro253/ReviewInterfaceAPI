package repository

import "github.com/jmoiron/sqlx"

// Хранилище
type Store struct {
	config *Config
}

//func New(config *Config) *Store {
//	//return &Config{
//	//	config: config,
//	//}
//}

func (s *Store) Open() error {
	return nil
}

func (s *Store) Close() {

}

type Repository struct {
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
