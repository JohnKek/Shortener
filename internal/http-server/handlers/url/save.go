package url

import (
	"UrlShort/utils"
	"UrlShort/utils/api/response"
	"UrlShort/utils/random"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
)

const (
	NEW    = "handlers.url.save.New"
	LENGTH = 10
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	Save(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", NEW),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("faiede to decode request body", utils.Err(err))
			render.JSON(w, r, response.Error("faiede to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", utils.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(LENGTH)
		}

		id, err := urlSaver.Save(req.URL, alias)

	}
}
