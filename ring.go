// Copyright (c) 2019 Tanner Ryan. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ring

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	errElements      = errors.New("error: elements must be greater than 0")
	errFalsePositive = errors.New("error: falsePositive must be greater than 0 and less than 1")
)

// Ring contains the information for a ring data store.
type Ring struct {
	elements int64
	fpr      float64
	size     uint64  // number of bits (bit array is size/8+1)
	bits     []uint8 // main bit array
	hash     uint64  // number of hash rounds
	mx       sync.RWMutex
}

// Init initializes and returns a new ring, or an error. Given a number of
// elements, it accurately states if data is not added. Within a falsePositive
// rate, it will indicate if the data has been added.
func Init(elements int, falsePositive float64) (*Ring, error) {
	if elements <= 0 {
		return nil, errElements
	}
	if falsePositive <= 0 || falsePositive >= 1 {
		return nil, errFalsePositive
	}

	r := Ring{elements: int64(elements), fpr: falsePositive}
	// number of bits
	m := (-1 * float64(elements) * math.Log(falsePositive)) / math.Pow(math.Log(2), 2)
	// number of hash operations
	k := (m / float64(elements)) * math.Log(2)

	r.size = uint64(math.Ceil(m))
	r.hash = uint64(math.Ceil(k))
	r.bits = make([]uint8, r.size/8+1)
	return &r, nil
}

// Add adds the data to the ring.
func (r *Ring) Add(data []byte) (new bool) {
	// generate hashes
	hash := generateMultiHash(data)

	r.mx.Lock()
	defer r.mx.Unlock()

	for i := uint64(0); i < r.hash; i++ {
		index := getRound(hash, i) % r.size

		if (r.bits[index/8] & (1 << (index % 8))) == 0 {
			new = true
		}

		r.bits[index/8] |= 1 << (index % 8)
	}
	return
}

// Reset clears the ring.
func (r *Ring) Reset() {
	r.mx.Lock()
	r.bits = make([]uint8, r.size/8+1)
	r.mx.Unlock()
}

// Test returns a bool if the data is in the ring. True indicates that the data
// may be in the ring, while false indicates that the data is not in the ring.
func (r *Ring) Test(data []byte) bool {
	// generate hashes
	hash := generateMultiHash(data)
	r.mx.RLock()
	defer r.mx.RUnlock()
	for i := uint64(0); i < r.hash; i++ {
		index := getRound(hash, i) % r.size
		// check if index%8-th bit is not active
		if (r.bits[index/8] & (1 << (index % 8))) == 0 {
			return false
		}

	}
	return true
}

// Merges the sent Ring into itself.
func (r *Ring) Merge(m *Ring) error {
	if r.size != m.size || r.hash != m.hash {
		return errors.New("rings must have the same m/k parameters")
	}

	r.mx.Lock()
	m.mx.RLock()
	defer r.mx.Unlock()
	defer m.mx.RUnlock()
	for i := 0; i < len(m.bits); i++ {
		r.bits[i] |= m.bits[i]
	}
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (r *Ring) MarshalBinary() ([]byte, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	out := make([]byte, len(r.bits)+17+16)
	// store a version for future compatibility
	out[0] = 2
	binary.BigEndian.PutUint64(out[1:9], r.size)
	binary.BigEndian.PutUint64(out[9:17], r.hash)
	binary.BigEndian.PutUint64(out[17:25], uint64(r.elements))
	binary.BigEndian.PutUint64(out[25:33], math.Float64bits(r.fpr))
	copy(out[33:], r.bits)
	return out, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (r *Ring) UnmarshalBinary(data []byte) error {
	// 17 bytes for version + size + hash and 1 byte at least for bits
	if len(data) < 17+1 {
		return fmt.Errorf("incorrect length: %d", len(data))
	}
	if data[0] == 2 {
		if len(data) < 17+1+16 {
			return fmt.Errorf("incorrect length: %d", len(data))
		}
	} else if data[0] != 1 {
		return fmt.Errorf("unexpected version: %d", data[0])
	}

	r.mx.Lock()
	defer r.mx.Unlock()
	r.size = binary.BigEndian.Uint64(data[1:9])
	r.hash = binary.BigEndian.Uint64(data[9:17])
	// sanity check against the bits being the wrong size
	if len(r.bits) != int(r.size/8+1) {
		r.bits = make([]uint8, r.size/8+1)
	}

	if data[0] == 2 {
		r.elements = int64(binary.BigEndian.Uint64(data[17:25]))
		r.fpr = math.Float64frombits(binary.BigEndian.Uint64(data[25:33]))
		copy(r.bits, data[33:])
	} else {
		copy(r.bits, data[17:])
	}
	return nil
}

func (r *Ring) GetElementsCount() int64 {
	return r.elements
}

func (r *Ring) GetFpr() float64 {
	return r.fpr
}
