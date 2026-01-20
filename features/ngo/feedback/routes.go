package feedback

import (
	"foodlink_backend/middleware"
	"net/http"
)

func SetupRoutes(service *Service, handler *Handler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			handler.GetAllFeedback(w, r)
		case r.Method == http.MethodPost:
			handler.CreateFeedback(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/stories", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			handler.GetAllStories(w, r)
		case r.Method == http.MethodPost:
			handler.CreateStory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return middleware.Chain(authMiddleware)(mux)
}
