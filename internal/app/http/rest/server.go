package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pedro-mealha/nba-stats-api/internal/app/domain/stats"
	"github.com/pedro-mealha/nba-stats-api/internal/app/gateway/nba"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.uber.org/zap"
)

type API struct {
	logger *zap.SugaredLogger
	s      stats.Provider
}

// NewAPI creates a new router with the needed endpoints
func NewAPI(logger *zap.SugaredLogger, s stats.Provider) *API {
	return &API{logger: logger, s: s}
}

// Routes exposes rest endpoints
func (a *API) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*pedromealha.dev", "http://localhost*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Route("/stats", func(r chi.Router) {
		r.Get("/scoreboard", a.getScoreboard)
		r.Get("/boxscore", a.getBoxscore)
	})

	return &ochttp.Handler{
		Handler:     r,
		Propagation: &b3.HTTPFormat{},
	}
}

func (a *API) getScoreboard(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		cmd nba.GetScoreboardCommand
	)

	cmd.Date = r.URL.Query().Get("date")
	cmd.LeagueID = nba.ParseLeague(r.URL.Query().Get("league"))

	res, err := a.s.GetScoreboard(ctx, cmd)
	if err != nil {
		a.logger.Errorw("failed to get scoreboard", "err", err)

		http.Error(w, "failed to get scoreboard", http.StatusInternalServerError)

		return
	}

	render.JSON(w, r, res)
}

func (a *API) getBoxscore(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		cmd nba.GetBoxscoreCommand
	)

	cmd.GameID = r.URL.Query().Get("gameId")

	res, err := a.s.GetBoxscore(ctx, cmd)
	if err != nil {
		a.logger.Errorw("failed to get boxscore", "err", err)

		http.Error(w, "failed to get boxscore", http.StatusInternalServerError)

		return
	}

	render.JSON(w, r, res)
}
