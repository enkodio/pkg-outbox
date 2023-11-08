package entity

type Header interface {
	GetKey() string
	GetValue() []byte
}

type Headers interface {
	SetHeader(key string, value []byte)
	GetValueByKey(key string) []byte
}
