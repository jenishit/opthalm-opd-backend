package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}

	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func nullUUID(value uuid.UUID) uuid.NullUUID {
	if value == uuid.Nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  value,
		Valid: true,
	}
}

func nullStringPtr(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *value,
		Valid:  true,
	}
}

func nullTimePtr(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  *value,
		Valid: true,
	}
}

func nullUUIDPtr(value *uuid.UUID) uuid.NullUUID {
	if value == nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  *value,
		Valid: true,
	}
}
