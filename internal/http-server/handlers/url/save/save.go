package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/customlog"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLen uint8 = 6

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func CreateUrl(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.NewUrlSaver"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", customlog.Error(err))
			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validationErr, ok := err.(validator.ValidationErrors)

			if !ok {
				var validationErrMsg string = "failed to validate"

				log.Error(validationErrMsg, customlog.Error(errors.New(validationErrMsg)))
				render.JSON(w, r, resp.Error(validationErrMsg))

				return
			}

			log.Error("invalid request body", customlog.Error(err))

			render.JSON(w, r, resp.ValidationError(validationErr))

			return
		}

		alias := req.Alias

		if alias == "" {
			aliasGenerated, err := random.NewRandomString(aliasLen)

			if err != nil {
				log.Error("internal error", customlog.Error(err))
				render.JSON(w, r, resp.Error("internal error: "+err.Error()))

				return
			}

			alias = aliasGenerated
		}

		id, err := urlSaver.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save url", customlog.Error(err))

			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}

		log.Info("url added", slog.Int64("url_id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
