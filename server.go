package main

import "github.com/go-martini/martini"

func main() {

  m := martini.Classic()

  m.Get("/", func() string {
    return "Hello world!"
  })

  m.Get("/hello/:name", func(params martini.Params) string {
    return "Hello " + params["name"]
  })  
  
  m.NotFound(func() {
    // handle 404
  })

  m.Run()

}