package cfutil

import (
	"io"
	"net/http"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	rpcStatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

var jsonpbMarshaler jsonpb.Marshaler
var once sync.Once

func WriteResponse(w http.ResponseWriter, r *http.Request, resp interface{}) error {
	if r.Header.Get("Accept") == "application/json" {
		once.Do(func() {
			jsonpbMarshaler.OrigName = true
		})
		if err := jsonpbMarshaler.Marshal(w, resp.(protoiface.MessageV1)); err != nil {
			return err
		}

		return nil
	}

	payload, _ := proto.Marshal(resp.(proto.Message))
	if _, err := w.Write(payload); err != nil {
		return err
	}

	return nil
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) error {
	w.WriteHeader(status)

	errorMessage := &rpcStatus.Status{
		Code:    int32(status),
		Message: err.Error(),
		Details: nil,
	}

	if r.Header.Get("Accept") == "application/json" {
		once.Do(func() {
			jsonpbMarshaler.OrigName = true
		})
		if err := jsonpbMarshaler.Marshal(w, errorMessage); err != nil {
			return err
		}

		return nil
	}

	payload, _ := proto.Marshal(errorMessage)
	if _, err := w.Write(payload); err != nil {
		return err
	}

	return nil
}

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

func ProtobufHandler(w http.ResponseWriter, r *http.Request, request proto.Message, do func(proto.Message) (proto.Message, error)) {
	if err := ApplyCors(w, r); err != nil {
		WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := ApplyContentType(w, r); err != nil {
		WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := ReadRequest(r, &request); err != nil {
		WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	response, err := do(request)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := WriteResponse(w, r, response); err != nil {
		WriteError(w, r, http.StatusInternalServerError, err)
	}
}
