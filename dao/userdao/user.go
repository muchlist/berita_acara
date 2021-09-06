package userdao

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/berita_acara/dao"
	"github.com/muchlist/berita_acara/db"
	"github.com/muchlist/berita_acara/dto"
	"github.com/muchlist/berita_acara/utils/logger"
	"github.com/muchlist/berita_acara/utils/rest_err"
	"github.com/muchlist/berita_acara/utils/sql_err"
	"time"
)

const (
	connectTimeout = 3

	keyUserTable = "users"
	keyID        = "id"
	keyEmail     = "email"
	keyName      = "name"
	keyPassword  = "password"
	keyRole      = "role"
	keyCreatedAt = "created_at"
	keyUpdatedAt = "updated_at"

	keyUsersRolesTable = "users_roles"
	keyUsersID         = "users_id"
	keyRolesID         = "roles_id"
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

func (u *userDao) Insert(user dto.User) (int, rest_err.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	user.Prepare()
	if len(user.Role) == 0 {
		return 0, rest_err.NewBadRequestError("role tidak boleh kosong")
	}

	// --------------------------------------------------------------
	trx, err := u.db.Begin(ctx)
	sqlStatement, args, err := u.sb.Insert(keyUserTable).Columns(keyID, keyEmail, keyName, keyPassword, keyCreatedAt, keyUpdatedAt).
		Values(user.ID, user.Email, user.Name, user.Password, user.CreatedAt, user.UpdatedAt).
		Suffix(dao.Returning(keyID)).
		ToSql()
	if err != nil {
		_ = trx.Rollback(ctx)
		return 0, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var userID int
	err = trx.QueryRow(ctx, sqlStatement, args...).Scan(&userID)
	if err != nil {
		logger.Error("error saat trx query user (Insert:0)", err)
		_ = trx.Rollback(ctx)
		return 0, sql_err.ParseError(err)
	}

	// -------------------------------------------------------------
	sqlInsert := u.sb.Insert(keyUsersRolesTable).Columns(keyRolesID, keyUsersID)
	for _, roleID := range user.Role {
		sqlInsert = sqlInsert.Values(roleID, userID)
	}
	sqlInsert = sqlInsert.Suffix(dao.Returning(keyID))
	sqlStatement, args, err = sqlInsert.ToSql()

	if err != nil {
		_ = trx.Rollback(ctx)
		return 0, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var usersRolesID int
	err = trx.QueryRow(ctx, sqlStatement, args...).Scan(&usersRolesID)
	if err != nil {
		logger.Error("error saat trx query usersRoles(Insert:1)", err)
		_ = trx.Rollback(ctx)
		return 0, sql_err.ParseError(err)
	}

	if err := trx.Commit(ctx); err != nil {
		return 0, rest_err.NewInternalServerError(dao.ErrCommit, err)
	}

	return userID, nil
}

func (u *userDao) Edit(input dto.User) (*dto.User, rest_err.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyEmail:     input.Email,
			keyName:      input.Name,
			keyUpdatedAt: input.UpdatedAt,
		}).
		Where(squirrel.Eq{
			keyID: input.ID,
		}).
		Suffix(dao.Returning(keyID, keyEmail, keyName, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var user dto.User
	err = db.DB.QueryRow(
		ctx,
		sqlStatement, args...).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u *userDao) ChangeRole(userID int, role []int) (*dto.User, rest_err.APIError) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	if len(role) == 0 {
		return nil, rest_err.NewBadRequestError("role tidak boleh kosong")
	}

	// --------------------------------------------------
	trx, err := u.db.Begin(ctx)
	sqlStatement, args, err := u.sb.Delete(keyUsersRolesTable).
		Where(squirrel.Eq{
			keyID: userID,
		}).ToSql()
	if err != nil {
		_ = trx.Rollback(ctx)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	_, err = db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		_ = trx.Rollback(ctx)
		logger.Error("error saat trx exec usersRoles(ChangeRole:0)", err)
		return nil, sql_err.ParseError(err)
	}

	// ---------------------------------------------------
	sqlInsert := u.sb.Insert(keyUsersRolesTable).Columns(keyRolesID, keyUsersID)
	for _, roleID := range role {
		sqlInsert = sqlInsert.Values(roleID, userID)
	}
	sqlStatement, args, err = sqlInsert.ToSql()
	if err != nil {
		_ = trx.Rollback(ctx)
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	_, err = trx.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat trx query usersRoles(Insert:1)", err)
		_ = trx.Rollback(ctx)
		return nil, sql_err.ParseError(err)
	}

	if err := trx.Commit(ctx); err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrCommit, err)
	}

	// -----------------------------------------------------
	userRes, err := u.Get(userID)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan user", err)
	}
	userRes.Role = role

	return userRes, nil
}

func (u *userDao) ChangePassword(input dto.User) (*dto.User, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyPassword:  input.Password,
			keyUpdatedAt: input.UpdatedAt,
		}).
		Where(keyID, input.ID).
		Suffix(dao.Returning(keyID, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	var user dto.User
	err = db.DB.QueryRow(context.Background(), sqlStatement, args).
		Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u *userDao) Delete(id int) rest_err.APIError {

	sqlStatement, args, err := u.sb.Delete(keyUserTable).
		Where(keyID, id).
		ToSql()
	if err != nil {
		return rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	res, err := db.DB.Exec(context.Background(), sqlStatement, args)
	if err != nil {
		return rest_err.NewInternalServerError("gagal saat penghapusan user", err)
	}

	if res.RowsAffected() != 1 {
		return rest_err.NewBadRequestError(fmt.Sprintf("User dengan username %d tidak ditemukan", id))
	}

	return nil
}

func (u *userDao) Get(id int) (*dto.User, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Select(keyID, keyEmail, keyName, keyPassword, keyCreatedAt, keyUpdatedAt).
		From(keyUserTable).
		Where(squirrel.Eq{
			keyID: id,
		}).
		ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	row := db.DB.QueryRow(context.Background(), sqlStatement, args...)

	var user dto.User
	err = row.Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}
	return &user, nil
}

func (u *userDao) Find() ([]dto.User, rest_err.APIError) {
	sqlStatement, _, err := u.sb.Select(keyID, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt).
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
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		users = append(users, user)
	}
	return users, nil
}
