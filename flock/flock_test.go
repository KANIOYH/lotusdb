package flock

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"sync/atomic"
	"testing"
)

func TestAcquireFileLock(t *testing.T) {
	testFn := func(readOnly bool, times int, actual int) {
		path := "/tmp/FLOCK"
		defer os.Remove(path)
		var count uint32

		wg := &sync.WaitGroup{}
		wg.Add(times)
		for i := 0; i < times; i++ {
			go func() {
				defer wg.Done()
				_, err := AcquireFileLock(path, readOnly)
				if err != nil {
					atomic.AddUint32(&count, 1)
				}
			}()
		}
		wg.Wait()
		assert.Equal(t, count, uint32(actual))
	}

	t.Run("exclusive-1", func(t *testing.T) {
		testFn(false, 1, 0)
	})

	t.Run("exclusive-2", func(t *testing.T) {
		testFn(false, 10, 9)
	})

	t.Run("exclusive-3", func(t *testing.T) {
		testFn(false, 500, 499)
	})

	t.Run("shared-1", func(t *testing.T) {
		testFn(true, 1, 0)
	})

	t.Run("shared-2", func(t *testing.T) {
		testFn(true, 500, 0)
	})
}

func TestFileLockGuard_Release(t *testing.T) {
	path := "/tmp/FLOCK"
	defer os.Remove(path)

	lock, err := AcquireFileLock(path, false)
	assert.Nil(t, err)
	err = lock.Release()
	assert.Nil(t, err)
}

func TestSyncDir(t *testing.T) {
	path := "/tmp/test-sync/"
	err := os.MkdirAll(path, os.ModePerm)
	assert.Nil(t, err)
	_, err = os.OpenFile(path+"test.txt", os.O_CREATE, 0644)
	assert.Nil(t, err)
	defer os.RemoveAll(path)

	err = SyncDir(path)
	assert.Nil(t, err)
}