package copier

type Copier[Src any, Dst any] interface {
	CopyTo(src *Src, dst *Dst) error
	Copy(src *Src) (*Dst, error)
}
