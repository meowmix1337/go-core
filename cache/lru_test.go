package cache

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LRUTestSuite struct {
	suite.Suite
	cache *lruCache
}

func TestLRUTestSuite(t *testing.T) {
	suite.Run(t, new(LRUTestSuite))
}

func (s *LRUTestSuite) SetupSuite() {

}

// SetupTest runs before each test in the suite
func (s *LRUTestSuite) SetupTest() {
	s.cache = NewLRUCache(10)
}

// TearDownTest runs after each test in the suite
func (s *LRUTestSuite) TearDownTest() {
}

func (s *LRUTestSuite) Test_SetGet() {
	key := "new"
	value := "test"

	s.cache.Set(context.Background(), key, value, 5)

	item, err := s.cache.Get(context.Background(), "new")
	if err != nil {
		s.T().Error(err)
	}

	if item != value {
		s.T().Errorf("expected %s, but got %s", value, item)
	}
}

func (s *LRUTestSuite) Test_SetExistsAlreadyAndGet() {
	for i := range 5 {
		s.cache.Set(context.Background(), strconv.Itoa(i), i, 5)
	}

	s.cache.Set(context.Background(), "0", 0, 5)

	value, err := s.cache.Get(context.Background(), "0")
	if err != nil {
		s.T().Error(err)
	}

	if value != 0 {
		s.T().Errorf("expected %v, but got %v", 0, value)
	}
}

func (s *LRUTestSuite) TestLRUCache_CacheMiss() {
	_, err := s.cache.Get(context.Background(), "does_not_exist")
	if !errors.Is(err, CacheMissErr) {
		s.T().Errorf("expected error %v, got %v", CacheMissErr, err)
	}
}

func (s *LRUTestSuite) TestLRUCache_Delete() {
	key := "new"
	value := "test"

	s.cache.Set(context.Background(), key, value, 5)

	_, err := s.cache.Get(context.Background(), key)
	if err != nil {
		s.T().Error(err)
	}

	s.cache.Delete(context.Background(), key)

	_, err = s.cache.Get(context.Background(), key)
	if !errors.Is(err, CacheMissErr) {
		s.T().Errorf("expected error %v, got %v", CacheMissErr, err)
	}
}

func (s *LRUTestSuite) TestLRUCache_DeleteNotExist() {
	err := s.cache.Delete(context.Background(), "does_not_exist")
	if err != nil {
		s.T().Error(err)
	}
}

func (s *LRUTestSuite) TestLRUCache_SizeAndPurge() {
	for i := range 5 {
		s.cache.Set(context.Background(), strconv.Itoa(i), i, 5)
	}

	if s.cache.Size(context.Background()) != 5 {
		s.T().Errorf("expected 5, but got %v", s.cache.Size(context.Background()))
	}

	s.cache.Purge(context.Background())
	if s.cache.Size(context.Background()) != 0 {
		s.T().Errorf("expected 0, but got %v", s.cache.Size(context.Background()))
	}
}

func (s *LRUTestSuite) TestLRUCache_Eviction() {
	// add ten items to cache
	for i := range 10 {
		s.cache.Set(context.Background(), strconv.Itoa(i), i, 5)
	}

	// add 11th item which should evict idx 0 since we inserted 0 first
	s.cache.Set(context.Background(), "11", 11, 60)

	_, err := s.cache.Get(context.Background(), "0")
	if !errors.Is(err, CacheMissErr) {
		s.T().Errorf("expected error %v, got %v", CacheMissErr, err)
	}

	value, err := s.cache.Get(context.Background(), "11")
	if err != nil {
		s.T().Error(err)
	}

	if value != 11 {
		s.T().Errorf("expected %d, but got %v", 11, value)
	}
}

func (s *LRUTestSuite) TestLRUCache_Expired() {

	s.cache.Set(context.Background(), "1", 1, 0)

	if s.cache.Size(context.Background()) != 1 {
		s.T().Errorf("expected 1, but got %v", s.cache.Size(context.Background()))
	}

	_, err := s.cache.Get(context.Background(), "1")
	if !errors.Is(err, CacheMissErr) {
		s.T().Errorf("expected error %v, got %v", CacheMissErr, err)
	}

	if s.cache.Size(context.Background()) != 0 {
		s.T().Errorf("expected 0, but got %v", s.cache.Size(context.Background()))
	}
}
