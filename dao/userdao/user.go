package userdao

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/berita_acara/db"
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/rest_err"
	"github.com/muchlist/berita_acara/utils/sql_err"
	"strings"
)

const (
	keyUserTable = "users"
	keyUsername  = "username"
	keyEmail     = "email"
	keyName      = "name"
	keyPassword  = "password"
	keyRole      = "role"
	keyCreatedAt = "created_at"
	keyUpdatedAt = "updated_at"
)

type userDao struct {
	db *pgxpool.Pool
	sb squirrel.StatementBuilderType
}

func New(db *pgxpool.Pool) UserDaoAssumer {
	return &userDao{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (u *userDao) Insert(user dto.User) (*string, rest_err.APIError) {

	sqlStatement, args, err := u.sb.Insert(keyUserTable).Columns(keyUsername, keyEmail, keyName, keyPassword, keyRole, keyCreatedAt, keyUpdatedAt).
		Values(user.Username, user.Email, user.Name, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).
		Suffix(returning(keyUsername)).
		ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	var userName dto.UppercaseString
	err = u.db.QueryRow(context.Background(), sqlStatement, args).Scan(&userName)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	usernameString := string(userName)
	return &usernameString, nil
}

func (u *userDao) Edit(input dto.User) (*dto.User, rest_err.APIError) {

	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyEmail:     input.Email,
			keyName:      input.Name,
			keyRole:      input.Role,
			keyUpdatedAt: input.UpdatedAt,
		}).
		Where(keyUsername, input.Name).
		Suffix(returning(keyUsername, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	var user dto.User
	err = db.DB.QueryRow(
		context.Background(),
		sqlStatement, args).Scan(&user.Username, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u *userDao) ChangePassword(input dto.User) (*dto.User, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyPassword:  input.Password,
			keyUpdatedAt: input.UpdatedAt,
		}).
		Where(keyUsername, input.Username).
		Suffix(returning(keyUsername, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	var user dto.User
	err = db.DB.QueryRow(context.Background(), sqlStatement, args).
		Scan(&user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u *userDao) Delete(userName string) rest_err.APIError {

	sqlStatement, args, err := u.sb.Delete(keyUserTable).
		Where(keyUsername, dto.UppercaseString(userName)).
		ToSql()
	if err != nil {
		return rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	res, err := db.DB.Exec(context.Background(), sqlStatement, args)
	if err != nil {
		return rest_err.NewInternalServerError("gagal saat penghapusan user", err)
	}

	if res.RowsAffected() != 1 {
		return rest_err.NewBadRequestError(fmt.Sprintf("User dengan username %s tidak ditemukan", userName))
	}

	return nil
}

func (u *userDao) Get(userName string) (*dto.User, rest_err.APIError) {

	sqlStatement, args, err := u.sb.Select(keyUsername, keyEmail, keyName, keyPassword, keyRole, keyCreatedAt, keyUpdatedAt).
		From(keyUserTable).
		Where(keyUsername, dto.UppercaseString(userName)).
		ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	row := db.DB.QueryRow(context.Background(), sqlStatement, args)

	var user dto.User
	err = row.Scan(&user.Username, &user.Email, &user.Name, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u *userDao) Find() ([]dto.User, rest_err.APIError) {
	sqlStatement, _, err := u.sb.Select(keyUsername, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt).
		From(keyUserTable).
		ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	rows, err := db.DB.Query(context.Background(), sqlStatement)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar user", err)
	}

	defer rows.Close()
	var users []dto.User
	for rows.Next() {
		user := dto.User{}
		err := rows.Scan(&user.Username, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func returning(columns ...string) string {
	sb := strings.Builder{}
	sb.WriteString("RETURNING ")
	for _, key := range columns {
		sb.WriteString(key + ", ")
	}
	return strings.TrimSuffix(sb.String(), ", ")
}
