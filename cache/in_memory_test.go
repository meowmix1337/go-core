package cache

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InMemoryCacheTestSuite struct {
	suite.Suite
	cache *InMemoryCache
}

func TestInMemoryCacheSuite(t *testing.T) {
	suite.Run(t, new(InMemoryCacheTestSuite))
}

func (s *InMemoryCacheTestSuite) SetupSuite() {

}

// SetupTest runs before each test in the suite
func (s *InMemoryCacheTestSuite) SetupTest() {
	s.cache = NewInMemoryCache()
}

// TearDownTest runs after each test in the suite
func (s *InMemoryCacheTestSuite) TearDownTest() {
}

func (s *InMemoryCacheTestSuite) TestInMemory_SetGet() {
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

func (s *InMemoryCacheTestSuite) TestInMemory_CacheMiss() {
	_, err := s.cache.Get(context.Background(), "does_not_exist")
	if !errors.Is(err, CacheMissErr) {
		s.T().Errorf("expected error %v, got %v", CacheMissErr, err)
	}
}

func (s *InMemoryCacheTestSuite) TestInMemory_Delete() {
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

func (s *InMemoryCacheTestSuite) TestInMemory_DeleteNotExist() {
	err := s.cache.Delete(context.Background(), "does_not_exist")
	if err != nil {
		s.T().Error(err)
	}
}

func (s *InMemoryCacheTestSuite) TestInMemory_SizeAndPurge() {
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
