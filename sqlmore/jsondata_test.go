package sqlmore

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

type UserData struct {
	Id   int
	Info JsonSql[UserInfo]
}

type UserInfo struct {
	Name string
	Age  int
}

func Test_jsondata(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	_, err = db.ExecContext(context.Background(), `create table user_data(id integer , info JsonSql)`)
	if err != nil {
		log.Println(err)
	}
	userinfo := UserInfo{
		Name: "test",
		Age:  18,
	}

	testCase := []struct {
		Name     string
		Input    any
		WantData any
		WantErr  error
	}{
		{
			Name: "just struct",
			Input: UserData{
				Id: 1,
				Info: JsonSql[UserInfo]{
					Column: userinfo,
					Valid:  true,
				},
			},
			WantData: UserData{
				Id: 1,
				Info: JsonSql[UserInfo]{
					Column: userinfo,
					Valid:  true,
				},
			},
			WantErr: nil,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := db.ExecContext(context.Background(), "insert into user_data(id,info) values(?,?)", tc.Input.(UserData).Id, tc.Input.(UserData).Info)
			assert.Equal(t, tc.WantErr, err)
			rows, _ := db.QueryContext(context.Background(), "select * from user_data")
			tm := &UserData{}
			for rows.Next() {
				rows.Scan(&tm.Id, &tm.Info)
			}
			assert.Equal(t, tc.WantData, *tm)
		})
	}
}
