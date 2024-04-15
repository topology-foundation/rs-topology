package pbcodec

import "google.golang.org/protobuf/proto"

type ProtoBufSerializer struct{}

func (p *ProtoBufSerializer) Marshal(data interface{}) ([]byte, error) {
	return proto.Marshal(data.(proto.Message))
}

func (p *ProtoBufSerializer) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
