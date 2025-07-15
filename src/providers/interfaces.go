package providers

type KeyVal interface {
	Get(key string) any
	set(key string, val any)
}
