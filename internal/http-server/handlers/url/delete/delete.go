package deleteurl

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

type URLRemover interface {
	DeleteUrl(alias string) error
}

func HandleDeleteUrl(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.HandleDeleteUrl"

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

		err := urlRemover.DeleteUrl(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("url not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete url", customlog.Error(err))
			render.JSON(w, r, resp.Error("failed to delete url"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())
	}
}
