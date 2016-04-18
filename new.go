package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"hash"
)

// Set the four most significant bits (bits 12 through 15) of the
// time_hi_and_version field to the 4-bit version number.
func (u *UUID) setVersion(v byte) {
	u[6] = (u[6] & 0xF) | (v << 4)
}

// Generate a MD5 hash of a namespace and a name, and copy it to the
// UUID slice.
func (u *UUID) setBytesFromHash(hash hash.Hash, ns, name []byte) {
	hash.Write(ns[:])
	hash.Write(name)
	copy(u[:], hash.Sum([]byte{})[:16])
}

// Set the two most significant bits (bits 6 and 7) of the
// clock_seq_hi_and_reserved to zero and one, respectively.
func (u *UUID) setVariant(v byte) {
	switch v {
	case ReservedNCS:
		u[8] = (u[8] | ReservedNCS) & 0xBF
	case ReservedRFC4122:
		u[8] = (u[8] | ReservedRFC4122) & 0x7F
	case ReservedMicrosoft:
		u[8] = (u[8] | ReservedMicrosoft) & 0x3F
	}
}

// Generate a UUID based on the MD5 hash of a namespace identifier
// and a name.
func NewV3(ns UUID, name []byte) (u UUID, err error) {
	// Set all bits to MD5 hash generated from namespace and name.
	u.setBytesFromHash(md5.New(), ns[:], name)
	u.setVariant(ReservedRFC4122)
	u.setVersion(3)
	return
}

// Generate a random UUID.
func NewV4() (u UUID, err error) {
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	_, err = rand.Read(u[:])
	if err != nil {
		return
	}
	u.setVariant(ReservedRFC4122)
	u.setVersion(4)
	return
}

// Generate a UUID based on the SHA-1 hash of a namespace identifier
// and a name.
func NewV5(ns UUID, name []byte) (u UUID, err error) {
	// Set all bits to truncated SHA1 hash generated from namespace
	// and name.
	u.setBytesFromHash(sha1.New(), ns[:], name)
	u.setVariant(ReservedRFC4122)
	u.setVersion(5)
	return
}
