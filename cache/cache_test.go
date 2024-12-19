package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string
}

func TestEmpty(t *testing.T) {
	key := "testKey"

	result := Get(key)
	assert.Nil(t, result, "result should be nil")
}

func TestSetAndGet(t *testing.T) {
	key := "testKey"
	value := TestStruct{
		Name: "testName",
	}
	expires := DefaultExpiration()

	Set(key, value, expires)

	result := Get(key)
	assert.Equal(t, value, result, "they should be equal")
}

func TestGetExpired(t *testing.T) {
	key := "testKey"
	value := "testValue"
	expires := time.Now().Add(-time.Hour) // expired time

	Set(key, value, expires)

	result := Get(key)
	assert.Nil(t, result, "result should be nil")
}

func TestInvalidate(t *testing.T) {
	key := "testKey"
	value := "testValue"
	expires := DefaultExpiration()

	Set(key, value, expires)
	Invalidate(key)

	result := Get(key)
	assert.Nil(t, result, "result should be nil")
}
