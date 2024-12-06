package stat

import (
	"fmt"
	"link-shortner-api/configs"
	"link-shortner-api/pkg/middleware"
	"link-shortner-api/pkg/response"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDependencies struct {
	*StatRepository
	Config *configs.Config
}

type StatHandler struct {
	*StatRepository
}

func NewStatHandler(router *http.ServeMux, dependencies *StatHandlerDependencies) {
	handler := &StatHandler{
		StatRepository: dependencies.StatRepository,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), dependencies.Config))
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		from, err := time.Parse("2006-01-02", req.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid from query: %s", err.Error()), http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", req.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid to query: %s", err.Error()), http.StatusBadRequest)
			return
		}

		by := req.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by query", http.StatusBadRequest)
			return
		}
		stats := handler.GetStats(by, from, to)
		response.ReturnJSON(w, http.StatusOK, stats)
	}
}
