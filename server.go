package main

import (
  "os"
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/codegangsta/martini-contrib/binding"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type Crib struct {
  Handle string `form:"handle"`
  URL  string `form:"url"`
  Type string `form:"type"`  
  Description string `form:"description"`
}

func DB() martini.Handler {
  session, err := mgo.Dial(os.Getenv("MONGO_URL")) // mongodb://localhost
  if err != nil {
    panic(err)
  }

  return func(c martini.Context) {
    s := session.Clone()
    c.Map(s.DB(os.Getenv("MONGO_DB")))
    defer s.Close()
    c.Next()
  }
}

func All(db *mgo.Database) []Crib {
  var cribs []Crib
  db.C("cribs").Find(nil).All(&cribs)
  return cribs
}

func Fetch(db *mgo.Database, handle string) Crib {
  var crib Crib
  db.C("cribs").Find(bson.M{"handle": handle}).One(&crib)
  return crib
}

func main() {

  m := martini.Classic()
  m.Use(render.Renderer(render.Options {
    Layout: "layout",
  }))
  m.Use(DB())

  m.Get("/", func(r render.Render, db *mgo.Database) {
    r.HTML(200, "list", All(db))
  })  

  m.Post("/", binding.Form(Crib{}), func(crib Crib, r render.Render, db *mgo.Database) {
    db.C("cribs").Insert(crib)
    r.HTML(200, "list", All(db))
  })  

  m.Get("/:handle", func(params martini.Params, r render.Render, db *mgo.Database) {
    r.HTML(200, "display", Fetch(db, params["handle"]))    
  })    

  m.Run()

}