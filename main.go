// to test locally: start username with "debug"

package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "embed"
)

//go:embed templates
var templates embed.FS

var addr = flag.String("addr", ":8100", "http service address")

func serveRegister(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, templates, "templates/register.html")
}

func servePlay(w http.ResponseWriter, r *http.Request) {
    http.ServeFileFS(w, r, templates, "templates/play.html")
}

func main() {
    fmt.Println("Dr. H. Bridge")

    game := NewGame()

    flag.Parse()
    hub := newHub(game)
    go hub.run()

    mux := http.NewServeMux()
    mux.HandleFunc("GET /drhbridge", serveRegister)
    mux.HandleFunc("GET /drhbridge/{$}", serveRegister)
    mux.HandleFunc("GET /drhbridge/play", servePlay)
    mux.HandleFunc("/drhbridge/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    srv := &http.Server{
        Addr: *addr,
        Handler: mux,
    }
    fmt.Printf("serving drhbridge at %v\n", *addr)

    fmt.Printf("drh hand %v\n", game.DrHHand())
    fmt.Printf("west hand %v\n", game.WestHand())
    fmt.Printf("teddy hand %v\n", game.TeddyHand())
    fmt.Printf("east hand %v\n", game.EastHand())

    err := srv.ListenAndServe()
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
