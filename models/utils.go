package models

import (
	"database/sql"
)

// MakeNullString is a convenience function for creating nullable strings
func MakeNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
