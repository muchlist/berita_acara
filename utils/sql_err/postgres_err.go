package sql_err

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/muchlist/berita_acara/utils/rest_err"
)

func ParseError(err error) rest_err.APIError {
	if err == pgx.ErrNoRows {
		return rest_err.NewBadRequestError(fmt.Sprintf("tidak ada data yang sesuai dengan id yang diberikan"))
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return rest_err.NewBadRequestError("input mengalami konflik dengan data existing")
		case pgerrcode.ForeignKeyViolation:
			return rest_err.NewBadRequestError("parent id tidak ditemukan")
		case pgerrcode.UndefinedColumn:
			return rest_err.NewInternalServerError("galat pada query database, column tidak tersedia", err)
		}
	}
	return rest_err.NewInternalServerError("galat", err)
}
