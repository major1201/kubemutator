package session

import (
	"math/rand"
	"strconv"
)

// RequestID of a session.
type RequestID uint32

// NewRequestID generates a new ID. The generated ID is high likely to be unique, but not cryptographically secure.
// The generated ID will never be 0.
func NewRequestID() RequestID {
	for {
		id := RequestID(rand.Uint32())
		if id != 0 {
			return id
		}
	}
}

// ToUint32 returns the uint32 form of ID
func (id RequestID) ToUint32() uint32 {
	return uint32(id)
}

// ToString returns the string form of ID
func (id RequestID) ToString() string {
	return strconv.FormatUint(uint64(id), 10)
}
