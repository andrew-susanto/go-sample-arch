package localcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewGoCache_SetGet(t *testing.T) {
	cache := NewGoCache()
	rsc := NewRepository(cache)

	// get non existing key
	val, exists := rsc.Get("key1")
	assert.Nil(t, val)
	assert.Equal(t, false, exists)

	// set new key
	err := rsc.SetNoExpiry("key1", "val1")
	assert.Nil(t, err)

	// set existing key
	err = rsc.SetNoExpiry("key1", "val2")
	assert.Nil(t, err)

	// get existing key; overwrited
	val, exists = rsc.Get("key1")
	assert.Equal(t, "val2", val)
	assert.Equal(t, true, exists)

	// get existing key as string
	val, exists = rsc.GetString("key1")
	assert.Equal(t, "val2", val)
	assert.Equal(t, true, exists)
}

func Test_GetSet_NonString(t *testing.T) {
	cache := NewGoCache()
	rsc := NewRepository(cache)

	tests := []struct {
		name string
		args interface{}
	}{
		{
			name: "mapStringInt",
			args: map[string]interface{}{"id": 123},
		},
		{
			name: "float",
			args: 123.456,
		},
		{
			name: "nil",
			args: nil,
		},
		{
			name: "struct",
			args: struct{ name string }{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := rsc.SetNoExpiry(tt.name, tt.args)
			assert.Nil(t, err)

			val, exists := rsc.Get(tt.name)
			assert.Equal(t, tt.args, val)
			assert.Equal(t, true, exists)
		})
	}

}

func Test_Delete(t *testing.T) {
	cache := NewGoCache()
	rsc := NewRepository(cache)

	// set new key
	err := rsc.SetNoExpiry("key1", "val1")
	assert.Nil(t, err)

	// delete non existing key
	err = rsc.Delete("key0")
	assert.Nil(t, err)

	// delete  existing key
	err = rsc.Delete("key1")
	assert.Nil(t, err)

	// get deleted key
	val, exists := rsc.Get("key1")
	assert.Nil(t, val)
	assert.Equal(t, false, exists)
}
