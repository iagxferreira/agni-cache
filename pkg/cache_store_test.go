package pkg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCacheStore(t *testing.T) {
	store := NewCacheStore(100 * time.Millisecond)
	defer store.StopCleanup()

	t.Run("store item", func(t *testing.T) {
		value := []byte("test-value")
		item := NewCacheItem(value, 5*time.Second)

		store.Set(item)

		retrievedItem, found := store.Get(item.ID)
		if !found {
			t.Fatal("item not found")
		}

		if string(retrievedItem.GetValue()) != string(value) {
			t.Errorf("expected value %s, got %s", string(value), string(retrievedItem.GetValue()))
		}
	})

	t.Run("store item", func(t *testing.T) {
		value := []byte("test-value")
		item := NewCacheItem(value, 5*time.Second)

		store.Set(item)

		newValue := []byte("new-value")
		item.SetValue(newValue)

		retrievedItem, found := store.Get(item.ID)
		if !found {
			t.Fatal("item not found")
		}

		if string(retrievedItem.GetValue()) != string(newValue) {
			t.Errorf("expected value %s, got %s", string(newValue), string(retrievedItem.GetValue()))
		}
	})

	t.Run("delete item", func(t *testing.T) {
		value := []byte("test-value")
		item := NewCacheItem(value, 5*time.Second)

		store.Set(item)
		store.Delete(item.ID)

		_, found := store.Get(item.ID)
		if found {
			t.Fatal("item not deleted")
		}
	})

	t.Run("not found item", func(t *testing.T) {
		_, found := store.Get(uuid.New())
		if found {
			t.Fatal("item found, but should not exist")
		}
	})

	t.Run("item expired", func(t *testing.T) {
		value := []byte("test-value")
		item := NewCacheItem(value, 1*time.Millisecond)

		store.Set(item)

		time.Sleep(200 * time.Millisecond)

		_, found := store.Get(item.ID)
		if found {
			t.Fatal("item not expired")
		}
	})
}
