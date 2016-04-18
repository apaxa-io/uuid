//Origins:
//nu7hatch/gouuid
//james-lawrence/gouuid
//applift/gouuid
//an2deg/gouuid

// This package provides immutable UUID structs and the functions
// NewV3, NewV4, NewV5 and Parse() for generating versions 3, 4
// and 5 UUIDs as specified in RFC 4122.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

const uuidLen = 16                  // UUID length in bytes
const uuidStringLen = uuidLen*2 + 4 // Default UUID string representation length in chars

// The UUID reserved variants.
const (
	ReservedNCS       byte = 0x80
	ReservedRFC4122   byte = 0x40
	ReservedMicrosoft byte = 0x20
	ReservedFuture    byte = 0x00
)

// TODO hide variable from world, overwise smbd may change them
// The following standard UUIDs are for use with NewV3() or NewV5().
var (
	NamespaceDNS, _  = ParseHex("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = ParseHex("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = ParseHex("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

// Pattern used to parse hex string representation of the UUID.
// FIXME: do something to consider both brackets at one time,
// current one allows to parse string with only one opening
// or closing bracket.
const hexPattern = "^(urn\\:uuid\\:)?\\{?([a-z0-9]{8})-([a-z0-9]{4})-" +
	"([1-5][a-z0-9]{3})-([a-z0-9]{4})-([a-z0-9]{12})\\}?$"

var re = regexp.MustCompile(hexPattern)

var nullUUID = UUID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// A UUID representation compliant with specification in
// RFC 4122 document.
type UUID [uuidLen]byte

// Parse creates a UUID object from given bytes slice.
func Parse(b []byte) (u UUID, err error) {
	if len(b) != uuidLen {
		err = errors.New("Given slice is not valid UUID sequence")
		return
	}
	copy(u[:], b)
	return
}

// ParseHex creates a UUID object from given hex string
// representation. Function accepts UUID string in following
// formats:
//
//     uuid.ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//     uuid.ParseHex("{6ba7b814-9dad-11d1-80b4-00c04fd430c8}")
//     uuid.ParseHex("urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//
func ParseHex(s string) (u UUID, err error) {
	md := re.FindStringSubmatch(strings.ToLower(s))
	if md == nil {
		err = errors.New("Invalid UUID string " + s)
		return
	}
	hash := md[2] + md[3] + md[4] + md[5] + md[6]
	b, err := hex.DecodeString(hash) // TODO redeclaration of err???
	if err != nil {
		return
	}
	copy(u[:], b)
	return
}

// ParseHex creates a UUID object from given hex string
// in lower case without any additional chars.
//
//     uuid.ParseClean("6ba7b8149dad11d180b400c04fd430c8")
//
func ParseClean(s string) (u UUID, err error) {
	if len(s) != uuidLen*2 {
		errors.New("Invalid UUID clean string " + s)
		return
	}

	b, err := hex.DecodeString(s) // TODO redeclaration of err???
	if err != nil {
		return
	}
	copy(u[:], b)
	return
}

// Variant returns the UUID Variant, which determines the internal
// layout of the UUID. This will be one of the constants: RESERVED_NCS,
// RFC_4122, RESERVED_MICROSOFT, RESERVED_FUTURE.
func (u UUID) Variant() byte {
	if u[8]&ReservedNCS == ReservedNCS {
		return ReservedNCS
	} else if u[8]&ReservedRFC4122 == ReservedRFC4122 {
		return ReservedRFC4122
	} else if u[8]&ReservedMicrosoft == ReservedMicrosoft {
		return ReservedMicrosoft
	}
	return ReservedFuture
}

// Version returns a version number of the algorithm used to
// generate the UUID sequence.
func (u UUID) Version() uint {
	return uint(u[6] >> 4)
}

func (u UUID) IsNull() bool {
	return u == nullUUID
}

// Returns string representation of UUID.
func (u UUID) String() string {
	const s byte = '-'
	buf := make([]byte, uuidStringLen)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = s
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = s
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = s
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = s
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

func (u UUID) GoString() string {
	return "\"" + u.String() + "\""
}

// Returns clean string representation of UUID (without any additional chars and in lower case)
func (u UUID) CleanString() string {
	return hex.EncodeToString(u[:])
}
