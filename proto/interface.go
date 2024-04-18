package pbcodec

type Serializer interface {
	Marshal(data interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}
