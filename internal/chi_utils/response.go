package chi_utils

import (
	"net/http"

	"github.com/go-chi/render"
)

type errResponse struct {
	HTTPStatusCode int           `json:"http_status_code"` // http response status code
	Description    errorResponse `json:"description"`
}

type Error struct {
	Description string `json:"description"`
	Code        int32  `json:"code"`
}

type errorResponse struct {
	Errors []Error `json:"errors,omitempty"`
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func InvalidRequest(err error) render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusBadRequest,
		Description: errorResponse{
			Errors: []Error{
				Error{
					Code:        http.StatusBadRequest,
					Description: err.Error(),
				},
			},
		},
	}
}

func InternalServerError(err error) render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		Description: errorResponse{
			Errors: []Error{
				Error{
					Code:        http.StatusInternalServerError,
					Description: err.Error(),
				},
			},
		},
	}
}

func NotImplementedError() render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusNotImplemented,
		Description: errorResponse{
			Errors: []Error{
				Error{
					Code:        http.StatusNotImplemented,
					Description: http.StatusText(http.StatusNotImplemented),
				},
			},
		},
	}
}

func Forbidden(err error) render.Renderer {
	return &errResponse{
		HTTPStatusCode: http.StatusForbidden,
		Description: errorResponse{
			Errors: []Error{
				Error{
					Code:        http.StatusForbidden,
					Description: err.Error(),
				},
			},
		},
	}
}
