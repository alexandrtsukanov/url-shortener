package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/customlog"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func HandleGetUrl(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.HandleGetUrl"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("invalid URL", customlog.Error(errors.New("invalid URL")))
			render.JSON(w, r, resp.Error("invalid URL"))

			return
		}

		resUrl, err := urlGetter.GetUrl(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("url not found"))

			return
		}

		if err != nil {
			log.Error("failed to get url", customlog.Error(err))
			render.JSON(w, r, resp.Error("failed to get url"))

			return
		}

		log.Info("url found", slog.String("url", resUrl))

		// http.Redirect(w, r, resUrl, http.StatusFound)
	}
}
