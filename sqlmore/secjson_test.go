package sqlmore

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type Userdata struct {
	Id   int
	Info EncryptColumn[Userinfo]
}

type Userinfo struct {
	Name string
	Age  int
}

func Test_secjsondata(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	_, err = db.ExecContext(context.Background(), `create table user_data(id integer , info EncryptColumn)`)
	require.NoError(t, err)
	userinfo := Userinfo{
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
			Input: Userdata{
				Id: 1,
				Info: EncryptColumn[Userinfo]{
					Val:   userinfo,
					Valid: true,
				},
			},
			WantData: Userdata{
				Id: 1,
				Info: EncryptColumn[Userinfo]{
					Val:   userinfo,
					Valid: true,
				},
			},
			WantErr: nil,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := db.ExecContext(context.Background(), "insert into user_data(id,info) values(?,?)", tc.Input.(Userdata).Id, tc.Input.(Userdata).Info)
			assert.Equal(t, tc.WantErr, err)
			rows, _ := db.QueryContext(context.Background(), "select * from user_data")
			tm := &Userdata{}
			for rows.Next() {
				rows.Scan(&tm.Id, &tm.Info)
			}
			assert.Equal(t, tc.WantData, *tm)
		})
	}
}

func Benchmark_encode_decode(b *testing.B) {

	a := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	// 创建一个固定长度为 16 的字节数组
	var c [16]byte
	// 手动将字节切片的内容复制到字节数组
	copy(c[:], a)

	b.Run("encode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			a := "hello world"
			sqlEncode(a, c)
		}
	})
	t := "hello world"

	ciphertext, _ := sqlEncode(t, c)
	key := md5.Sum([]byte(reflect.TypeOf(ciphertext).String()))
	b.Run("decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sqlDecode(ciphertext, key[:16])
		}
	})
}

func Test_byte_buffer(t *testing.T) {
	var a int32 = 1443242342
	fmt.Println(byte(a))
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, a)
	fmt.Println(buffer.Bytes())
}
