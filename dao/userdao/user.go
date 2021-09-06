package userdao

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
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
	keyRolesName       = "roles_name"

	keyRolesTable = "roles"
	keyRoleName   = "role_name"
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
	if len(user.Roles) == 0 {
		return 0, rest_err.NewBadRequestError("role tidak boleh kosong")
	}

	// -------------------------------------------------------------- insert user data
	trx, err := u.db.Begin(ctx)
	defer func(trx pgx.Tx) {
		_ = trx.Rollback(context.Background())
	}(trx)

	sqlStatement, args, err := u.sb.Insert(keyUserTable).Columns(keyID, keyEmail, keyName, keyPassword, keyCreatedAt, keyUpdatedAt).
		Values(user.ID, user.Email, user.Name, user.Password, user.CreatedAt, user.UpdatedAt).
		Suffix(dao.Returning(keyID)).
		ToSql()
	if err != nil {
		return 0, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var userID int
	err = trx.QueryRow(ctx, sqlStatement, args...).Scan(&userID)
	if err != nil {
		logger.Error("error saat trx query user (Insert:0)", err)
		return 0, sql_err.ParseError(err)
	}

	// ------------------------------------------------------------- insert role data
	sqlInsert := u.sb.Insert(keyUsersRolesTable).Columns(keyRolesName, keyUsersID)
	for _, roleName := range user.Roles {
		sqlInsert = sqlInsert.Values(roleName, userID)
	}
	sqlInsert = sqlInsert.Suffix(dao.Returning(keyRolesName))
	sqlStatement, args, err = sqlInsert.ToSql()

	if err != nil {
		return 0, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	var rolesName string
	err = trx.QueryRow(ctx, sqlStatement, args...).Scan(&rolesName)
	if err != nil {
		logger.Error("error saat trx query usersRoles(Insert:1)", err)
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

	if len(input.Roles) == 0 {
		return nil, rest_err.NewBadRequestError("role tidak boleh kosong")
	}

	// ------------------------------------------------------------------------- begin
	trx, err := u.db.Begin(ctx)
	defer func(trx pgx.Tx) {
		_ = trx.Rollback(context.Background())
	}(trx)

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

	// ------------------------------------------------------------------------- user edit
	var user dto.User
	err = trx.QueryRow(
		ctx,
		sqlStatement, args...).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, sql_err.ParseError(err)
	}

	// ------------------------------------------------------------------------- role delete
	sqlStatement, args, err = u.sb.Delete(keyUsersRolesTable).
		Where(squirrel.Eq{
			keyUsersID: input.ID,
		}).ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	_, err = db.DB.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat trx exec usersRoles(ChangeRole:0)", err)
		return nil, sql_err.ParseError(err)
	}

	// ---------------------------------------------------
	sqlInsert := u.sb.Insert(keyUsersRolesTable).Columns(keyRolesName, keyUsersID)
	for _, roleName := range input.Roles {
		sqlInsert = sqlInsert.Values(roleName, user.ID)
	}
	sqlStatement, args, err = sqlInsert.ToSql()
	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	_, err = trx.Exec(ctx, sqlStatement, args...)
	if err != nil {
		logger.Error("error saat trx query usersRoles(Insert:1)", err)
		return nil, sql_err.ParseError(err)
	}

	if err := trx.Commit(ctx); err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrCommit, err)
	}

	user.Roles = input.Roles

	return &user, nil
}

func (u *userDao) ChangePassword(input dto.User) rest_err.APIError {
	sqlStatement, args, err := u.sb.Update(keyUserTable).
		SetMap(squirrel.Eq{
			keyPassword:  input.Password,
			keyUpdatedAt: input.UpdatedAt,
		}).
		Where(keyID, input.ID).
		Suffix(dao.Returning(keyID, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt)).
		ToSql()

	if err != nil {
		return rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	res, err := db.DB.Exec(context.Background(), sqlStatement, args...)
	if err != nil {
		return sql_err.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("User dengan username %d tidak ditemukan", input.ID))
	}

	return nil
}

func (u *userDao) Delete(id int) rest_err.APIError {

	sqlStatement, args, err := u.sb.Delete(keyUserTable).
		Where(keyID, id).
		ToSql()
	if err != nil {
		return rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	}

	res, err := db.DB.Exec(context.Background(), sqlStatement, args...)
	if err != nil {
		return rest_err.NewInternalServerError("gagal saat penghapusan user", err)
	}

	if res.RowsAffected() == 0 {
		return rest_err.NewBadRequestError(fmt.Sprintf("User dengan username %d tidak ditemukan", id))
	}

	return nil
}

func (u *userDao) Get(id int) (*dto.User, rest_err.APIError) {
	sqlStatement, args, err := u.sb.Select(
		dao.B(keyRolesName),
		dao.A(keyID),
		dao.A(keyEmail),
		dao.A(keyName),
		dao.A(keyPassword),
		dao.A(keyCreatedAt),
		dao.A(keyUpdatedAt),
	).
		Distinct().
		From(keyUserTable + " A").
		Join(keyUsersRolesTable + " B ON A.id = B.users_id").
		Where(squirrel.Eq{
			dao.A(keyID): id,
		}).
		ToSql()

	if err != nil {
		return nil, rest_err.NewInternalServerError(dao.ErrSqlBuilder, err)
	}

	rows, err := db.DB.Query(context.Background(), sqlStatement, args...)
	if err != nil {
		return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar user", err)
	}
	defer rows.Close()

	var userRes dto.User
	for rows.Next() {
		user := dto.User{}
		var roleName string
		err := rows.Scan(&roleName, &user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, sql_err.ParseError(err)
		}
		if userRes.ID == 0 {
			userRes.ID = user.ID
			userRes.Email = user.Email
			userRes.Name = user.Name
			userRes.Password = user.Password
			userRes.UpdatedAt = user.UpdatedAt
			userRes.CreatedAt = user.CreatedAt
		}
		userRes.Roles = append(userRes.Roles, roleName)
	}
	return &userRes, nil
}

func (u *userDao) Find() ([]dto.User, rest_err.APIError) {
	//sqlStatement, _, err := u.sb.Select(keyID, keyEmail, keyName, keyRole, keyCreatedAt, keyUpdatedAt).
	//	From(keyUserTable).
	//	ToSql()
	//if err != nil {
	//	return nil, rest_err.NewInternalServerError("kesalahan pada sql builder", err)
	//}
	//rows, err := db.DB.Query(context.Background(), sqlStatement)
	//if err != nil {
	//	return nil, rest_err.NewInternalServerError("gagal mendapatkan daftar user", err)
	//}
	//
	//defer rows.Close()
	var users []dto.User
	//for rows.Next() {
	//	user := dto.User{}
	//	err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	//	if err != nil {
	//		return nil, sql_err.ParseError(err)
	//	}
	//	users = append(users, user)
	//}
	return users, nil
}
