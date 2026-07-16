package repository

import (
	"database/sql"

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

