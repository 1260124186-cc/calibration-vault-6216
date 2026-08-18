package httpapi

import "errors"

var errMultipleDocuments = errors.New("request body must contain one JSON object")
