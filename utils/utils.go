package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/influxdata/toml"
)

const (
	// CharsetUTF8 ...
	CharsetUTF8 = "charset=utf-8"

	// MediaTypes
	ApplicationJSON                  = "application/json"
	ApplicationJSONCharsetUTF8       = ApplicationJSON + "; " + CharsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + CharsetUTF8
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + CharsetUTF8
	ApplicationForm                  = "application/x-www-form-urlencoded"
	ApplicationProtobuf              = "application/protobuf"
	TextHTML                         = "text/html"
	TextHTMLCharsetUTF8              = TextHTML + "; " + CharsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + CharsetUTF8
	ApplicationToml                  = "application/toml"
	ApplicationTomlCharsetUTF8       = ApplicationToml + "; " + CharsetUTF8
	MultipartForm                    = "multipart/form-data"
)

var (
	// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
	ErrJSONPayloadEmpty = errors.New("JSON payload is empty")

	// ErrTomlPayloadEmpty is returned when the Toml payload is empty.
	ErrTomlPayloadEmpty = errors.New("Toml payload is empty")

	startTime = time.Now().Format("2006-01-02 15:04:05")
)

// ResponseJOSN reponse josn
func ResponseJOSN(w http.ResponseWriter, code int, v interface{}) {
	var re []byte
	var err error

	re, err = json.Marshal(v)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", ApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	w.Write(re)
}

// QueryJSON decode json from http.Request.Body
func QueryJSON(r *http.Request, v interface{}) error {
	rb := NewRequestBody(r.Body)

	content, err := rb.Bytes()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return ErrJSONPayloadEmpty
	}
	return json.Unmarshal(content, v)
}

// ResponseToml reponse toml
func ResponseToml(w http.ResponseWriter, code int, v interface{}) {
	var re []byte
	var err error

	re, err = toml.Marshal(v)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", ApplicationTomlCharsetUTF8)
	w.WriteHeader(code)
	w.Write(re)
}

// QueryToml decode json from http.Request.Body
func QueryToml(r *http.Request, v interface{}) error {
	rb := NewRequestBody(r.Body)

	content, err := rb.Bytes()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return ErrTomlPayloadEmpty
	}
	return toml.Unmarshal(content, v)
}
