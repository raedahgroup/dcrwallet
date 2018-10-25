package chacharng

import (
	"errors"
	"fmt"
	"io"

	"github.com/tmthrgd/go-rand"
)

// RandBytes returns random [rndsize]byte slice with provided seed
// based on chacha20. The seed size must be 32.
func RandBytes(seed []byte, rndsize int) ([]byte, error) {
	r, err := rand.New(seed[:])
	if err != nil {
		fmt.Println("error creates chacharng bytes with seed")
		return nil, err
	}

	ret := make([]byte, rndsize)

	n, err := r.Read(ret)

	if n != rndsize {
		return nil, errors.New("error returns wrong size of bytes")
	}

	return ret, nil

}

// NewReaderBytes generates random [rndsize]byte slice with provided seed and reader for next random
// based on chacha20. The seed size must be 32.
func NewReaderBytes(seed []byte, rndsize int) (io.Reader, []byte, error) {
	r, err := NewRandReader(seed[:])
	if err != nil {
		fmt.Println("error creates chacharng bytes with seed")
		return nil, nil, err
	}

	ret := make([]byte, rndsize)

	n, _ := r.Read(ret)

	if n != rndsize {
		return nil, nil, errors.New("error returns wrong size of bytes")
	}

	return r, ret, nil

}

// NewRandReader creates a new rand reader from the provided seed.
func NewRandReader(seed []byte) (io.Reader, error) {
	r, err := rand.New(seed[:])
	if err != nil {
		fmt.Println("error new chacharng with seed")
		return nil, err
	}

	return r, nil
}
