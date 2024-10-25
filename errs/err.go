package errs

import (
	"errors"
	"fmt"
	"reflect"
)

func NewErrIndexOutOfRange(length, index int) error {
	return fmt.Errorf("ekit: 下标超出范围，长度 %d, 下标 %d", length, index)
}

func NewErrNeedStructPoint() error {
	return fmt.Errorf("ekit: 只支持一级结构体指针")
}

func NewErrSrcInValid() error {
	return fmt.Errorf("ekit: src值无效")
}

func NewErrTypeError(typ reflect.Type) error {
	return fmt.Errorf("ekit: copier 入口只支持 Struct 不支持类型 %v, 种类 %v", typ, typ.Kind())
}

func NewErrDstInValid() error {
	return fmt.Errorf("ekit: dst值无效")
}

func NewErrKindNotMatchError(src reflect.Kind, dst reflect.Kind, field string) error {
	return fmt.Errorf("ekit: 字段 %s 的 Kind 不匹配, src: %v, dst: %v", field, src, dst)
}

func NewErrorBeyondMaxcnt() error {
	return errors.New("ekit: 超过最大重试次数")
}
