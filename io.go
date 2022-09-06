package cfutil

import (
	"io"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

func ReadRequest(r *http.Request, message interface{}) error {
	// check content type
	if r.Header.Get("Content-Type") == "application/json" {

		if err := jsonpb.Unmarshal(r.Body, message.(protoiface.MessageV1)); err != nil {
			return err
		}

		return nil
	}

	content, _ := io.ReadAll(r.Body)

	if err := proto.Unmarshal(content, message.(proto.Message)); err != nil {
		return err
	}

	return nil
}
