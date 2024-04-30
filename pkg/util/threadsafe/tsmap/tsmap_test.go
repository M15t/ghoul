package tsmap

// Sets a value v to a key k in a map m.
import (
	"sync"
	"testing"

	"github.com/M15t/gram/pkg/util/threadsafe"
	"github.com/stretchr/testify/assert"
)

func TestSetWithValue(t *testing.T) {
	// Arrange
	mux := &sync.Mutex{}
	m := make(map[int]string)
	k := 1
	v := "value"

	// Act
	Set(mux, m, k, v)

	// Assert
	assert.Equal(t, v, m[k])
}

func TestSetWithNilMutex(t *testing.T) {
	// Arrange
	var mux threadsafe.Locker
	m := make(map[int]string)
	k := 1
	v := "value"

	// Act & Assert
	assert.Panics(t, func() { Set(mux, m, k, v) })
}

// Locks the mutex before setting the value.
func TestSetLocksMutexBeforeSettingValue(t *testing.T) {
	// Arrange
	mux := &sync.Mutex{}
	m := make(map[int]string)
	k := 1
	v := "value"

	// Act
	Set(mux, m, k, v)

	// Assert
	assert.Equal(t, v, m[k])
}

// Unlocks the mutex after setting the value.
func TestSetUnlocksMutexAfterSettingValue(t *testing.T) {
	// Arrange
	mux := &sync.Mutex{}
	m := make(map[int]string)
	k := 1
	v := "value"

	// Act
	Set(mux, m, k, v)

	// Assert
	assert.Equal(t, v, m[k])
}

// Can be safely used between multiple goroutines.
func TestSetCanBeSafelyUsedBetweenMultipleGoroutines(t *testing.T) {
	// Arrange
	mux := &sync.Mutex{}
	m := make(map[int]string)
	k := 1
	v := "value"

	var wg sync.WaitGroup
	wg.Add(2)

	// Act
	go func() {
		defer wg.Done()
		Set(mux, m, k, v)
	}()

	go func() {
		defer wg.Done()
		Set(mux, m, k, v)
	}()

	wg.Wait()

	// Assert
	assert.Equal(t, v, m[k])
}
