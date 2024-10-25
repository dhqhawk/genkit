package compare

// 该接口用于用户自定义某些类型
// 目前我们将提供切片以及map的比较
// TODO:这里我们考虑是否把any转化为compareable的类型
type ComponentChecker[T any] interface {
	IsSubset(SuperSet []T, SubSet ...T) (bool, error)
	GetIntersection(SetA []T, SetB []T) T
}
