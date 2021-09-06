package userdao

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/muchlist/berita_acara/dao"
	"testing"
)

func TestUserDao_Insert(t *testing.T) {
	id := 2
	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	sqlStatement, args, err := sb.Select(
		dao.B(keyRolesID),
		dao.A(keyID),
		dao.A(keyEmail),
		dao.A(keyName),
		dao.A(keyPassword),
		dao.A(keyCreatedAt),
		dao.A(keyUpdatedAt),
		dao.C("role_name"),
	).
		Distinct().
		From(keyUserTable + " A").
		Join(keyUsersRolesTable + " B ON A.id = B.users_id").
		Join("roles C ON C.id = B.roles_id").
		Where(squirrel.Eq{
			dao.A(keyID): id,
		}).
		ToSql()
	fmt.Println(sqlStatement)
	fmt.Printf("%v", args)
	fmt.Println(err)
}
