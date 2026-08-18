package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func decodeBody(request *http.Request, target any) error {
	defer request.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return errMultipleDocuments
		}
		return err
	}
	return nil
}

func rejectMalformedBody(writer http.ResponseWriter, err error) {
	writeProblem(writer, http.StatusBadRequest, "invalid_json", err.Error())
}
