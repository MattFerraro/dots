package main

import (
  "github.com/codegangsta/martini"
  "github.com/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  // render html templates from templates directory
  m.Use(render.Renderer())

  m.Get("/", func(r render.Render) {
    r.HTML(200, "base", "yolo")
  })

  m.Run()
}