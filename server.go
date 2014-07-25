package main

import (
  "os"
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/codegangsta/martini-contrib/binding"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)

type Crib struct {
  Handle string `form:"handle"`
  URL  string `form:"url"`
  Type string `form:"type"`  
  Description string `form:"description"`
}

/* Multi-
   called we initialize a Mongo session on localhost. DB() returns a martini.Handler which will be called on every request. We simply clone the session for every request and make sure it is closed once the request is done being processed. The important bit is the call to c.Map. This maps an instance of *mgo.Database to our request context. This allows all subsequent handler functions to specify a *mgo.Database as an argument and get it injected.
*/
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