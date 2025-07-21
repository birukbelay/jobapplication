package providers

type KeyValServ interface {
	Get(key string) any
	set(key string, val any)
}
