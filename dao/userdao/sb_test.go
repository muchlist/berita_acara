package userdao

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"testing"
)

func TestUserDao_Insert(t *testing.T) {

	search := "asdadad"
	cursor := 0
	limit := 0

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sqlfrom := sb.Select(keyID, keyEmail, keyName, keyCreatedAt, keyUpdatedAt).
		From(keyUserTable)

	sqlfrom = sqlfrom.Where(squirrel.And{
		squirrel.Gt{keyID: cursor},
		squirrel.Like{"name": fmt.Sprint("%", search, "%")},
	})

	sqlStatement, args, err := sqlfrom.OrderBy(keyID + " ASC").
		Limit(uint64(limit)).
		ToSql()

	fmt.Println(sqlStatement)
	fmt.Printf("%v\n", args)
	fmt.Println(err)
}
