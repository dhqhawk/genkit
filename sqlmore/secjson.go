package sqlmore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"io"
	"reflect"
)

// EncryptColumn 代表一个加密的列
// 一般来说加密可以选择依赖于数据库进行加密
// EncryptColumn 并不打算使用极其难破解的加密算法
// 而是选择使用 AES GCM 模式。
// 如果你觉得安全性不够，那么你可以考虑自己实现类似的结构体.
type EncryptColumn[T any] struct {
	Val T
	// Valid 为 true 的时候，Val 才有意义
	Valid bool
	get   getkey
}

func (e *EncryptColumn[T]) GetKey(fn getkey) *EncryptColumn[T] {
	e.get = fn
	return e
}

type getkey func() [16]byte

// Value 返回加密后的值
// 如果 T 是基本类型，那么会对 T 进行直接加密
// 否则，将 T 按照 JSON 序列化之后进行加密，返回加密后的数据
func (e EncryptColumn[T]) Value() (driver.Value, error) {
	if !e.Valid {
		return nil, nil
	}
	var key [16]byte
	if e.get != nil {
		key = e.get()
	} else {
		key = md5.Sum([]byte(reflect.TypeOf(e.Val).String()))
	}
	return sqlEncode(e.Val, key)
}

// Scan 方法会把写入的数据转化进行解密，
// 并将解密后的数据进行反序列化，构造 T
func (e *EncryptColumn[T]) Scan(value any) error {
	if value == nil {
		e.Val, e.Valid = *new(T), false
	}
	e.Valid = true
	var key [16]byte
	if e.get != nil {
		key = e.get()
	} else {
		key = md5.Sum([]byte(reflect.TypeOf(e.Val).String()))
	}
	switch s := value.(type) {
	case string:
		json.Unmarshal(sqlDecode([]byte(s), key[:]), &e.Val)
	case []byte:
		json.Unmarshal(sqlDecode(s, key[:]), &e.Val)
	}
	return nil
}

func sqlEncode(plain any, key [16]byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// 生成随机的初始化向量 (IV)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	var plaintext []byte
	switch p := plain.(type) {
	case string:
		plaintext = []byte(p)
	case int:
		tmp := int64(p)
		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, tmp)
		plaintext = buffer.Bytes()
	case uint:
		tmp := uint64(p)
		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, tmp)
		plaintext = buffer.Bytes()
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		buffer := new(bytes.Buffer)
		err = binary.Write(buffer, binary.BigEndian, p)
		plaintext = buffer.Bytes()
	case []byte:
		plaintext = p
	default:
		plaintext, err = json.Marshal(p)
	}
	// 加密数据并附加认证标签
	if err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func sqlDecode(ciphertext []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	// 创建 GCM 模式的 AEAD 解密器
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据并验证认证标签
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil
	}
	return plaintext
}
