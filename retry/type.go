package retry

type RestryStrategy interface {
	Next() error
}
