package mock

// Index ...
type Index interface {
	Get(key string) any
	GetTwo(key1, key2 string) (v1, v2 any)
	Put(key string, value any)
}

type Embed interface {
	RegularMethod()
	Embedded
}

type Embedded interface {
	EmbeddedMethod()
}
