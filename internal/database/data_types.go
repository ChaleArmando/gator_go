package database

import (
	"database/sql"
	"time"
)

func StringNull(str string) sql.NullString {
	if str != "" {
		return sql.NullString{String: str, Valid: true}
	}
	return sql.NullString{}
}

func TimeNull(tm string) sql.NullTime {
	publishedAt := sql.NullTime{}
	if t, err := time.Parse(time.RFC1123Z, tm); err == nil {
		publishedAt = sql.NullTime{
			Time:  t,
			Valid: true,
		}
	}
	return publishedAt
}
