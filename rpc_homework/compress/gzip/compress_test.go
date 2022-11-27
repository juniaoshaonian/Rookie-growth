package gzip

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompresser(t *testing.T) {
	testcases := []struct {
		name    string
		val     []byte
		wantErr error
	}{
		{
			name: "xxx",
			val:  []byte("hello world"),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			com := NewCompresser()
			encodeVal, err := com.Compress(tc.val)
			require.NoError(t, err)
			data, err := com.Decompress(encodeVal)
			require.NoError(t, err)
			assert.Equal(t, tc.val, data)
		})
	}

}
