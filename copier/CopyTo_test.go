package copier

import (
	"fmt"
	"genkit/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Intr struct {
	Ia int
	IB []string
}

type Src struct {
	A string
	B *int
	c string
	D []int
	E Intr
	F string
}

type Dst struct {
	A string
	B *int
	c string
	D []int
	E Intr
	F int
	G string
}

var simpleDst = &SimpleDst{}
var bascDst = &BasicDst{}
var embedDst = &EmbedDst{}
var complexDst = &ComplexDst{}
var specialDst = &SpecialDst{}
var notMatchDet = &NotMatchDst{}
var multiPtrDst = &MultiPtrDst{}
var a int = 10
var b int = 20
var c *int = &a

//var simpleDst *SimpleDst

func Test_isStructPoint(t *testing.T) {
	_, _, err := isStructPoint(simpleDst)
	fmt.Println(err)
}
func Test_copy(t *testing.T) {
	t.Parallel()
	testCase := []struct {
		name     string
		copyFunc func() error
		dst      any
		wantErr  error
		wantDst  any
	}{
		{
			name: "simple struct",
			copyFunc: func() error {
				copier := ReflectCopier[SimpleSrc, SimpleDst]{}
				return copier.CopyTo(&SimpleSrc{
					Name:    "大明",
					Age:     &a,
					Friends: []string{"Tom", "Jerry"},
				}, simpleDst)
			},
			wantDst: &SimpleDst{
				Name:    "大明",
				Age:     &a,
				Friends: []string{"Tom", "Jerry"},
			},
		},
		{
			name: "基础类型的 struct",
			copyFunc: func() error {
				copier := ReflectCopier[BasicSrc, BasicDst]{}
				return copier.CopyTo(&BasicSrc{
					Name:    "大明",
					Age:     10,
					CNumber: complex(1, 2),
				}, bascDst)
			},
			wantDst: &BasicDst{
				Name:    "大明",
				Age:     10,
				CNumber: complex(1, 2),
			},
		},
		{
			name: "src 是基础类型",
			copyFunc: func() error {
				copier := ReflectCopier[int, int]{}
				i := 10
				j := 2
				return copier.CopyTo(&i, &j)
			},
			wantErr: errs.NewErrTypeError(reflect.TypeOf(10)),
		},
		{
			name: "dst 是基础类型",
			copyFunc: func() error {
				copier := ReflectCopier[SimpleSrc, string]{}
				j := "ok"
				return copier.CopyTo(&SimpleSrc{
					Name:    "大明",
					Age:     &a,
					Friends: []string{"Tom", "Jerry"},
				}, &j)
			},
			wantErr: errs.NewErrTypeError(reflect.TypeOf("")),
		},
		{
			name: "接口类型",
			copyFunc: func() error {
				copier := ReflectCopier[InterfaceSrc, InterfaceDst]{}
				i := InterfaceSrc(10)
				j := InterfaceDst(1)
				return copier.CopyTo(&i, &j)
			},
			wantErr: errs.NewErrTypeError(reflect.TypeOf(new(InterfaceSrc)).Elem()),
		},
		{
			name: "simple struct 空切片,空指针",
			copyFunc: func() error {
				copier := ReflectCopier[SimpleSrc, SimpleDst]{}
				return copier.CopyTo(&SimpleSrc{
					Name: "大明",
				}, simpleDst)
			},
			wantDst: &SimpleDst{
				Name: "大明",
			},
		},
		{
			name: "组合 struct,浅拷贝",
			copyFunc: func() error {
				copier := ReflectCopier[EmbedSrc, EmbedDst]{}
				return copier.CopyTo(&EmbedSrc{
					SimpleSrc: SimpleSrc{
						Name:    "xiaoli",
						Age:     &a,
						Friends: []string{},
					},
					BasicSrc: &BasicSrc{
						Name:    "xiaowang",
						Age:     20,
						CNumber: complex(2, 2),
					},
				}, embedDst)
			},
			wantDst: &EmbedDst{
				SimpleSrc: SimpleSrc{
					Name:    "xiaoli",
					Age:     &a,
					Friends: []string{},
				},
				BasicSrc: &BasicSrc{
					Name:    "xiaowang",
					Age:     20,
					CNumber: complex(2, 2),
				},
			},
		},
		{
			name: "复杂struct,浅拷贝",
			copyFunc: func() error {
				copier := ReflectCopier[ComplexSrc, ComplexDst]{}
				return copier.CopyTo(&ComplexSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     &a,
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     &b,
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
				}, complexDst)
			},
			wantDst: &ComplexDst{
				Simple: SimpleDst{
					Name:    "xiaohong",
					Age:     &a,
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedDst{
					SimpleSrc: SimpleSrc{
						Name:    "xiaopeng",
						Age:     &b,
						Friends: []string{"la", "ha", "le"},
					},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     22,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			},
		},
		{
			name: "特殊类型",
			copyFunc: func() error {
				copier := ReflectCopier[SpecialSrc, SpecialDst]{}
				return copier.CopyTo(&SpecialSrc{
					Arr: [3]float32{1, 2, 3},
					M: map[string]int{
						"ha": 1,
						"o":  2,
					},
				}, specialDst)
			},
			wantDst: &SpecialDst{
				Arr: [3]float32{1, 2, 3},
				M: map[string]int{
					"ha": 1,
					"o":  2,
				},
			},
		},
		{
			name: "复杂struct,不匹配",
			copyFunc: func() error {
				copier := ReflectCopier[NotMatchSrc, NotMatchDst]{}
				return copier.CopyTo(&NotMatchSrc{
					Simple: SimpleSrc{
						Name:    "xiaohong",
						Age:     &a,
						Friends: []string{"ha", "ha", "le"},
					},
					Embed: &EmbedSrc{
						SimpleSrc: SimpleSrc{
							Name:    "xiaopeng",
							Age:     &b,
							Friends: []string{"la", "ha", "le"},
						},
						BasicSrc: &BasicSrc{
							Name:    "wang",
							Age:     22,
							CNumber: complex(2, 1),
						},
					},
					BasicSrc: BasicSrc{
						Name:    "wang11",
						Age:     22,
						CNumber: complex(2, 1),
					},
					S: struct{ A string }{A: "a"},
				}, notMatchDet)
			},
			wantDst: &ComplexDst{
				Simple: SimpleDst{
					Name:    "xiaohong",
					Age:     &a,
					Friends: []string{"ha", "ha", "le"},
				},
				Embed: &EmbedDst{
					SimpleSrc: SimpleSrc{
						Name:    "xiaopeng",
						Age:     &b,
						Friends: []string{"la", "ha", "le"},
					},
					BasicSrc: &BasicSrc{
						Name:    "wang",
						Age:     22,
						CNumber: complex(2, 1),
					},
				},
				BasicSrc: BasicSrc{
					Name:    "wang11",
					Age:     22,
					CNumber: complex(2, 1),
				},
			},
			wantErr: errs.NewErrKindNotMatchError(reflect.TypeOf("a").Kind(), reflect.TypeOf(1).Kind(), "A"),
		},
		{
			name: "多重指针",
			copyFunc: func() error {
				copier := ReflectCopier[MultiPtrSrc, MultiPtrDst]{}
				return copier.CopyTo(&MultiPtrSrc{
					Name:    "a",
					Age:     &c,
					Friends: nil,
				}, multiPtrDst)
			},
			wantDst: &MultiPtrDst{
				Name:    "a",
				Age:     &c,
				Friends: nil,
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.copyFunc()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			switch reflect.TypeOf(tc.wantDst).Elem() {
			case reflect.TypeOf(SimpleDst{}):
				assert.Equal(t, tc.wantDst, simpleDst)
			case reflect.TypeOf(BasicDst{}):
				assert.Equal(t, tc.wantDst, bascDst)
			case reflect.TypeOf(EmbedDst{}):
				assert.Equal(t, tc.wantDst, embedDst)
			case reflect.TypeOf(SpecialDst{}):
				assert.Equal(t, tc.wantDst, specialDst)
			case reflect.TypeOf(MultiPtrDst{}):
				assert.Equal(t, tc.wantDst, multiPtrDst)
			}

		})
	}
}

type SimpleSrc struct {
	Name    string
	Age     *int
	Friends []string
}

type SimpleDst struct {
	Name    string
	Age     *int
	Friends []string
}

type BasicSrc struct {
	Name    string
	Age     int
	CNumber complex64
}

type BasicDst struct {
	Name    string
	Age     int
	CNumber complex64
}

type InterfaceSrc interface {
}

type InterfaceDst interface {
}

type EmbedSrc struct {
	SimpleSrc
	*BasicSrc
}

type EmbedDst struct {
	SimpleSrc
	*BasicSrc
}

type ComplexSrc struct {
	Simple SimpleSrc
	Embed  *EmbedSrc
	BasicSrc
}

type ComplexDst struct {
	Simple SimpleDst
	Embed  *EmbedDst
	BasicSrc
}

type SpecialSrc struct {
	Arr [3]float32
	M   map[string]int
}

type SpecialDst struct {
	Arr [3]float32
	M   map[string]int
}

type NotMatchSrc struct {
	Simple SimpleSrc
	Embed  *EmbedSrc
	BasicSrc
	S struct {
		A string
	}
}

type NotMatchDst struct {
	Simple SimpleDst
	Embed  *EmbedDst
	BasicSrc
	S struct {
		A int
	}
}

type MultiPtrSrc struct {
	Name    string
	Age     **int
	Friends []string
}

type MultiPtrDst struct {
	Name    string
	Age     **int
	Friends []string
}
