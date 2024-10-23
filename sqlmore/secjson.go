package sqlmore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"database/sql/driver"
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
}

// Value 返回加密后的值
// 如果 T 是基本类型，那么会对 T 进行直接加密
// 否则，将 T 按照 JSON 序列化之后进行加密，返回加密后的数据
func (e EncryptColumn[T]) Value() (driver.Value, error) {
	if !e.Valid {
		return nil, nil
	}
	return sqlEncode(e.Val)
}

// Scan 方法会把写入的数据转化进行解密，
// 并将解密后的数据进行反序列化，构造 T
func (e *EncryptColumn[T]) Scan(value any) error {
	if value == nil {
		e.Val, e.Valid = *new(T), false
	}
	e.Valid = true
	key := md5.Sum([]byte(reflect.TypeOf(e.Val).String()))
	switch s := value.(type) {
	case string:
		json.Unmarshal(sqlDecode([]byte(s), key[:16]), &e.Val)
	case []byte:
		json.Unmarshal(sqlDecode(s, key[:16]), &e.Val)
	}
	return nil
}

func sqlEncode(plain any) ([]byte, error) {
	key := md5.Sum([]byte(reflect.TypeOf(plain).String()))
	block, err := aes.NewCipher(key[:16])
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
	// 加密数据并附加认证标签
	plaintext, err := json.Marshal(plain)
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
