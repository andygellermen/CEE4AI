package http

import (
	stdhttp "net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andygellermen/CEE4AI/internal/answers"
	"github.com/andygellermen/CEE4AI/internal/packages"
	"github.com/andygellermen/CEE4AI/internal/questions"
	"github.com/andygellermen/CEE4AI/internal/results"
	"github.com/andygellermen/CEE4AI/internal/scoring"
	"github.com/andygellermen/CEE4AI/internal/sessions"
)

func NewRouter(pool *pgxpool.Pool) stdhttp.Handler {
	mux := stdhttp.NewServeMux()

	questionRepo := questions.NewRepository(pool)
	sessionRepo := sessions.NewRepository(pool)
	packageRepo := packages.NewRepository(pool)
	answerRepo := answers.NewRepository(pool)
	resultRepo := results.NewRepository(pool)
	scoringService := scoring.NewService(pool)
	packageService := packages.NewService(packageRepo, questionRepo)
	sessionService := sessions.NewService(sessionRepo, packageService)
	questionService := questions.NewService(sessionRepo, answerRepo, questionRepo, packageService)
	answerService := answers.NewService(sessionRepo, questionRepo, answerRepo, packageService, scoringService)
	resultService := results.NewService(sessionRepo, answerRepo, questionRepo, resultRepo, scoringService)

	sessionsHandler := newSessionHandler(sessionService)
	questionsHandler := newQuestionsHandler(questionService)
	answersHandler := newAnswersHandler(answerService)
	resultsHandler := newResultsHandler(resultService)

	mux.HandleFunc("GET /healthz", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /api/v1/sessions", sessionsHandler.create)
	mux.HandleFunc("GET /api/v1/questions/next", questionsHandler.next)
	mux.HandleFunc("POST /api/v1/answers", answersHandler.create)
	mux.HandleFunc("GET /api/v1/results", resultsHandler.get)

	return mux
}
