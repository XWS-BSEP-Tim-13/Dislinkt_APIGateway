package api

import (
	"encoding/json"
	dom "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/domain"
	"io"
	"net/http"
)

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func decodePostDtoBody(r io.Reader) (*dom.PostDto, error) {

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var rt dom.PostDto
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}
