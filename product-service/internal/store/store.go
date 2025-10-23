package store

import (
	"sync"

	"github.com/phanthehoang2503/product-service/internal/model"
)

type Store struct {
	sync.Mutex
	nextID   int64
	products map[int64]model.Product
}

func NewStore() *Store {
	return &Store{
		nextID:   1,
		products: make(map[int64]model.Product),
	}
}

func (s *Store) Create(p model.Product) model.Product {
	s.Lock()
	defer s.Unlock()
	p.ID = s.nextID
	s.nextID++
	s.products[p.ID] = p
	return p
}

func (s *Store) List() []model.Product {
	s.Lock()
	defer s.Unlock()

	out := make([]model.Product, 0, len(s.products))
	for _, v := range s.products {
		out = append(out, v)
	}
	return out
}

func (s *Store) Get(id int64) (model.Product, bool) {
	s.Lock()
	defer s.Unlock()

	p, ok := s.products[id]
	return p, ok
}

func (s *Store) Update(id int64, p model.Product) (model.Product, bool) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.products[id]
	if !ok {
		return model.Product{}, false
	}

	p.ID = id
	s.products[id] = p

	return p, true
}

func (s *Store) Delete(id int64) bool {
	s.Lock()
	defer s.Unlock()

	_, ok := s.products[id]
	if !ok {
		return false
	}
	delete(s.products, id)
	return true
}
