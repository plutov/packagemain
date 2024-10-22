package is

import "testing"

func FuzzEqual(f *testing.F) {
	f.Add([]byte{'f', 'u', 'z', 'z'}, []byte{'t', 'e', 's', 't'})

	f.Fuzz(func(t *testing.T, a []byte, b []byte) {
		Equal(a, b)
	})
}
