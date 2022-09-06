package cfutil

import (
	"net/http"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func ApplyCors(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	return nil
}

func ApplyContentType(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", "application/x-protobuf")
	}

	return nil
}

var jsonpbMarshaler jsonpb.Marshaler
var once sync.Once

func WriteResponse(w http.ResponseWriter, r *http.Request, resp proto.Message) error {
	if r.Header.Get("Accept") == "application/json" {
		once.Do(func() {
			jsonpbMarshaler.OrigName = true
		})
		if err := jsonpbMarshaler.Marshal(w, resp); err != nil {
			return err
		}

		return nil
	}

	payload, _ := proto.Marshal(resp)
	if _, err := w.Write(payload); err != nil {
		return err
	}

	return nil
}
