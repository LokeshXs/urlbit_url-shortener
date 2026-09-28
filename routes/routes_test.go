package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lokeshxs/url-shortener/db"
	"github.com/gin-gonic/gin"
)

func TestDatabaseOutage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, connection := range []string{"", "postgres://localhost:1/unavailable?sslmode=disable&connect_timeout=1"} {
		t.Run(connection, func(t *testing.T) {
			t.Setenv("POSTGRES_URL", connection)
			if err := db.InitDB(); err == nil {
				t.Fatal("expected database initialization error")
			}
			t.Cleanup(func() {
				if db.DB != nil {
					db.DB.Close()
					db.DB = nil
				}
			})
			server := gin.New()
			RoutingHandler(server)
			for _, endpoint := range []struct{ method, path string }{
				{http.MethodPost, "/shorten"},
				{http.MethodGet, "/stats"},
				{http.MethodGet, "/abcdef"},
				{http.MethodDelete, "/abcdef"},
				{http.MethodPost, "/webhook/signup"},
			} {
				response := httptest.NewRecorder()
				server.ServeHTTP(response, httptest.NewRequest(endpoint.method, endpoint.path, nil))
				if response.Code != http.StatusInternalServerError || response.Body.String() != `{"message":"Internal server error"}` {
					t.Errorf("%s %s: got %d %s", endpoint.method, endpoint.path, response.Code, response.Body.String())
				}
			}
			response := httptest.NewRecorder()
			server.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health-check", nil))
			if response.Code != http.StatusOK {
				t.Fatalf("health check during outage: got %d", response.Code)
			}
		})
	}
}
