package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var applogger *log.Logger

type Application struct {
	Server *http.Server
}

func NewApplication() *Application {
	return &Application{
		Server: &http.Server{
			Addr:           ":8080",
			Handler:        nil,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (app *Application) SetRoutes() *Application {
	mux := http.NewServeMux()
	mux.Handle("/encode", Recover(http.HandlerFunc(Encode)))
	mux.Handle("/decode", Recover(http.HandlerFunc(Decode)))

	app.Server.Handler = mux

	return app
}

func main() {
	applogger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	app := NewApplication().SetRoutes()

	applogger.Println("Server start.")
	applogger.Fatal(app.Server.ListenAndServe())
}

func Encode(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	ip := params.Get("ip")

	fmt.Fprintf(w, "%v", IPEncode(ip))
}

func Decode(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	ip := params.Get("ip")
	ipInt, _ := strconv.Atoi(ip)

	fmt.Fprintf(w, "%v", IPDecode(ipInt))
}

func IPDecode(encoded int) string {
	a := (encoded >> 24) & 0xFF
	b := (encoded >> 16) & 0xFF
	c := (encoded >> 8) & 0xFF
	d := encoded & 0xFF

	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

func IPEncode(ip string) int {
	var nums [4]int

	_, err := fmt.Sscanf(ip, "%d.%d.%d.%d", &nums[0], &nums[1], &nums[2], &nums[3])
	if err != nil {
		panic(err)
	}

	return (nums[0] << 24) | (nums[1] << 16) | (nums[2] << 8) | nums[3]
}

func Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				applogger.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
