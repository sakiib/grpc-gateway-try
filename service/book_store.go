package service

import (
	"errors"
	pb "github.com/sakiib/grpc-gateway-demo/gen/go/proto"
	"sync"
)

type BookStore interface {
	Set(book *pb.Book) error
	Get(id string) (*pb.Book, error)
}

type InMemStore struct {
	mutex sync.Mutex
	data  map[string]*pb.Book
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		data: make(map[string]*pb.Book),
	}
}

func (store *InMemStore) Set(book *pb.Book) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if book == nil {
		return errors.New("book not found")
	}

	if _, exists := store.data[book.Id]; exists {
		return errors.New("book with the given id already exists")
	}

	store.data[book.Id] = book
	return nil
}

func (store *InMemStore) Get(id string) (*pb.Book, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, exists := store.data[id]; !exists {
		return nil, errors.New("book with the given id not found")
	}
	return store.data[id], nil
}