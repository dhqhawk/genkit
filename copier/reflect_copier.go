package copier

import (
	"genkit/errs"
	"reflect"
)

type ReflectCopier[Src any, Dst any] struct {
	deep bool
}

func (r *ReflectCopier[Src, Dst]) CopyTo(src *Src, dst *Dst) error {
	r.deep = true
	//检验src是否为一阶结构体指针
	srcVal, srcTyp, err := isStructPoint(src)
	if err != nil {
		return err
	}
	//检验dst是否为一阶结构体指针
	dstVal, _, err := isStructPoint(dst)
	if err != nil {
		return err
	}

	srcMetnum := srcTyp.NumField()
	for i := 0; i < srcMetnum; i++ {
		srcField := srcTyp.Field(i)
		//避免非导出字段
		if !srcField.IsExported() {
			continue
		}
		//具有相同名字的可导出字段以及该字段类型相同
		if dstfield := dstVal.FieldByName(srcField.Name); dstfield.CanSet() {
			if dstfield.Kind() != srcField.Type.Kind() {
				return errs.NewErrKindNotMatchError(srcField.Type.Kind(), dstfield.Kind(), srcField.Name)
			}
			if r.deep {
				err := deepCopy(srcVal.Field(i), dstfield)
				if err != nil {
					return err
				}
			} else {
				dstfield.Set(srcVal.Field(i))
			}
		}
	}
	return nil
}

func (r *ReflectCopier[Src, Dst]) Copy(src *Src) (*Dst, error) {

	//TODO implement me
	panic("implement me")

}

func isStructPoint(entity any) (reflect.Value, reflect.Type, error) {
	// 先检查 entity 是否为 nil，避免进入反射时的问题
	if entity == nil {
		return reflect.Value{}, nil, errs.NewErrTypeError(reflect.TypeOf(entity).Elem())
	}
	entityVal := reflect.ValueOf(entity)
	if !entityVal.IsValid() {
		return reflect.Value{}, nil, errs.NewErrTypeError(reflect.TypeOf(entity).Elem())
	}
	if entityVal.IsZero() {
		return reflect.Value{}, nil, errs.NewErrTypeError(reflect.TypeOf(entity).Elem())
	}
	entityTyp := reflect.TypeOf(entity)
	if entityTyp.Kind() != reflect.Ptr || entityTyp.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, nil, errs.NewErrTypeError(reflect.TypeOf(entity).Elem())
	}
	entityTyp = entityTyp.Elem()
	entityVal = entityVal.Elem()
	return entityVal, entityTyp, nil
}

func deepCopy(src, dst reflect.Value) error {
	if !src.IsValid() {
		return errs.NewErrSrcInValid()
	}
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			//TODO:如果出现嵌套初始状态的话,需要在该处把之前的逻辑判断一次
			if src.Field(i).Type().Kind() != dst.Field(i).Type().Kind() {
				srcFieldType := src.Type().Field(i)
				return errs.NewErrKindNotMatchError(src.Field(i).Type().Kind(), dst.Field(i).Type().Kind(), srcFieldType.Name)
			}
			deepCopy(src.Field(i), dst.Field(i))
		}
	case reflect.Slice:
		if src.Cap() == 0 {
			return nil
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			deepCopy(src.Index(i), dst.Index(i))
		}
	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			deepCopy(src.Index(i), dst.Index(i))
		}
	case reflect.Map:
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			value := src.MapIndex(key)
			newValue := reflect.New(value.Type()).Elem()
			newKey := reflect.New(key.Type()).Elem()
			deepCopy(value, newValue)
			deepCopy(key, newKey)
			dst.SetMapIndex(newKey, newValue)
		}
	case reflect.Pointer:
		if src.IsNil() {
			dst.Set(reflect.Zero(src.Type()))
		} else {
			elemSrc := src.Elem()
			elemDst := reflect.New(dst.Type().Elem()).Elem()
			deepCopy(elemSrc, elemDst)
			dst.Set(elemDst.Addr())
		}
	default:
		dst.Set(src)
	}
	return nil
}
