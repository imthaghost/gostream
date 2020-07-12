package tests

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/fatih/color"
	"github.com/imthaghost/gostream/server"
)

// Todo: Edge Cases
func TestServer(t *testing.T) {
	// instantiate a server
	s := server.NewServer()
	// start a server
	go s.Start(":8000")
	// GET request to server root route
	resp, err := http.Get("http://127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
	}
	// stats code
	result := resp.StatusCode

	if result != 200 {
		t.Error()
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("%s Server Failed: expected 200 got %d \n", red("[-]"), result)

	} else {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s Passing: %d \n", green("[+]"), result)
	}
}
