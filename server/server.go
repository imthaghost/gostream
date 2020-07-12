package server

import (
	"github.com/imthaghost/gostream/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StreamAPI structure
type StreamAPI struct {
	e *echo.Echo
}

// NewServer Instance of Echo
func NewServer() *StreamAPI {

	return &StreamAPI{
		e: echo.New(),
	}
}

// Start server functionality
func (s *StreamAPI) Start(port string) {
	// Logger
	s.e.Use(middleware.Logger())
	// CORS
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// root endpoint
	s.e.GET("/", controllers.Hello)
	// Start Server
	s.e.Logger.Fatal(s.e.Start(port))
}

// Close server functionality
func (s *StreamAPI) Close() {
	s.e.Close()
}
