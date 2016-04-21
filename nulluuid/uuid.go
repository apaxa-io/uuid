package nulluuid

import (
	"fmt"
	"github.com/apaxa-io/uuid"
	"github.com/jackc/pgx"
)

// UUID represent UUID that may be NULL.
// UUID implements the pgx.Scanner and pgx.Encoder interfaces.
type UUID struct {
	UUID  uuid.UUID
	Valid bool // Valid is true if UUID is not NULL
}

// Scan implements the pgx.Scanner interface.
func (u *UUID) Scan(vr *pgx.ValueReader) error {
	if vr.Type().DataType != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Scan cannot decode %s (OID %d)", vr.Type().DataTypeName, vr.Type().DataType))
	}

	if vr.Len() == -1 {
		u.UUID, u.Valid = uuid.Null(), false
		return nil
	}

	u.Valid = true
	return u.UUID.Scan(vr)
}

// FormatCode implements the pgx.Encoder interface.
func (u UUID) FormatCode() int16 {
	return pgx.BinaryFormatCode
}

// Encode implements the pgx.Encoder interface.
func (u UUID) Encode(w *pgx.WriteBuf, oid pgx.Oid) error {
	if oid != pgx.UuidOid {
		return pgx.SerializationError(fmt.Sprintf("UUID.Encode cannot encode into OID %d", oid))
	}

	if !u.Valid {
		w.WriteInt32(-1)
		return nil
	}

	return u.UUID.Encode(w, oid)
}

// FromUUID construct nullable UUID from given UUID.
func FromUUID(u uuid.UUID) UUID {
	return UUID{UUID: u, Valid: true}
}

// Null construct invalid UUID.
func Null() UUID {
	return UUID{Valid: false}
}
