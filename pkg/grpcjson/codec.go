package grpcjson

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

type Codec struct{}

func (c Codec) Name() string {
	return "json"
}

func (c Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func Register() {
	encoding.RegisterCodec(Codec{})
}
