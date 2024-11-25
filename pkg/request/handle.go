package request

import (
	"link-shortner-api/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		response.ReturnJSON(*w, http.StatusBadRequest, err.Error())
		return nil, err
	}
	err = Validate[T](body)
	if err != nil {
		response.ReturnJSON(*w, http.StatusBadRequest, err.Error())
		return nil, err
	}

	return &body, nil
}
