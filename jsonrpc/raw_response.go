package jsonrpc

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/node-real/go-pkg/log"
)

// RawResponse keeps Result and Error as unparsed JSON
// It is meant to be used to deserialize JSONPRC responses from downstream components
// while Response is meant to be used to craft our own responses to clients.
type RawResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      ID               `json:"id"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *json.RawMessage `json:"error,omitempty"`
}

// MarshalJSON implements json.Marshaler and adds the "jsonrpc":"2.0"
// property.
func (r RawResponse) MarshalJSON() ([]byte, error) {

	if r.Error != nil {
		response := struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      ID              `json:"id"`
			Error   json.RawMessage `json:"error,omitempty"`
		}{
			JSONRPC: "2.0",
			ID:      r.ID,
			Error:   *r.Error,
		}

		return json.Marshal(response)
	} else {
		response := struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      ID              `json:"id"`
			Result  json.RawMessage `json:"result,omitempty"`
		}{
			JSONRPC: "2.0",
			ID:      r.ID,
			Result:  r.Result,
		}

		if response.Result == nil {
			response.Result = jsonNull
		}

		return json.Marshal(response)
	}
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *RawResponse) UnmarshalJSON(data []byte) error {
	type tmpType RawResponse

	if err := json.Unmarshal(data, (*tmpType)(r)); err != nil {
		return err
	}
	return nil
}

// =====================================================================================================================

const VSN = "2.0"

func NewResultRawResponse(result []byte, id ID) *RawResponse {
	return &RawResponse{
		JSONRPC: VSN,
		ID:      id,
		Result:  result,
	}
}

func NewErrorResponse(error *Error, id ID) *RawResponse {
	var rawError json.RawMessage
	rawError, err := jsoniter.Marshal(&error)
	if err != nil {
		log.Errorf("err:%v, when marshal in errorMessage", err)
	}
	return &RawResponse{
		JSONRPC: VSN,
		ID:      id,
		Error:   &rawError,
	}
}

func NewNullResponse(id ID) *RawResponse {
	return &RawResponse{
		JSONRPC: VSN,
		ID:      id,
		Result:  []byte("null"),
	}
}

func NewInvalidInput(message string, id ID) *RawResponse {
	return NewErrorResponse(InvalidInput(message), id)
}

func NewInvalidParams(id ID) *RawResponse {
	return NewErrorResponse(InvalidParams("invalid params"), id)
}

func NewMissingParams(id ID, index int) *RawResponse {
	return NewErrorResponse(InvalidParams(fmt.Sprintf("missing value for required argument %v", index)), id)
}

func NewInvalidParamsWithMessage(message string, id ID) *RawResponse {
	return NewErrorResponse(InvalidParams(message), id)
}

func NewInternalError(id ID) *RawResponse {
	return NewErrorResponse(InternalError("internal error"), id)
}

func NewParseError(id ID) *RawResponse {
	return NewErrorResponse(ParseError("parse error"), id)
}

func NewResourceUnavailable(message string, id ID) *RawResponse {
	return NewErrorResponse(ResourceUnavailable(message), id)
}

func NewResourceNotFound(message string, id ID) *RawResponse {
	return NewErrorResponse(ResourceNotFound(message), id)
}
