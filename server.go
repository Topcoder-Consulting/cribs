package main

import (
  "os"
  "net/http"
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
  "github.com/codegangsta/martini-contrib/binding"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)

// the Crib struct that we can serialize and deserialize into Mongodb
type Crib struct {
  Handle string `form:"handle"`
  URL  string `form:"url"`
  Type string `form:"type"`  
  Description string `form:"description"`
}

/* 
   the function returns a martini.Handler which is called on each request. We simply clone 
   the session for each request and close it when the request is complete. The call to c.Map 
   maps an instance of *mgo.Database to the request context. Then *mgo.Database
   is injected into each handler function.
*/
func DB() martini.Handler {
  session, err := mgo.Dial(os.Getenv("MONGO_URL")) // mongodb://localhost
  if err != nil {
    panic(err)
  }

  return func(c martini.Context) {
    s := session.Clone()
    c.Map(s.DB(os.Getenv("MONGO_DB"))) // local
    defer s.Close()
    c.Next()
  }
}

// function to return an array of all Cribs from mondodb
func All(db *mgo.Database) []Crib {
  var cribs []Crib
  db.C("cribs").Find(nil).All(&cribs)
  return cribs
}

// function to return a specific Crib by handle
func Fetch(db *mgo.Database, handle string) Crib {
  var crib Crib
  db.C("cribs").Find(bson.M{"handle": handle}).One(&crib)
  return crib
}

func main() {

  m := martini.Classic()
  // specify the layout to use when rendering HTML
  m.Use(render.Renderer(render.Options {
    Layout: "layout",
  }))
  // use the Mongo middleware
  m.Use(DB())

  // list of all cribs
  m.Get("/", func(r render.Render, db *mgo.Database) {
    r.HTML(200, "list", All(db))
  })  

  /* 
    create a new crib the form submission. Contains some martini magic. The call 
    to binding.Form(Crib{}) parses out form data when the request comes in. 
    It binds the data to the struct, maps it to the request context  and
    injects into our next handler function to insert into Mongodb.
 */   
  m.Post("/", binding.Form(Crib{}), func(crib Crib, r render.Render, db *mgo.Database) {
    db.C("cribs").Insert(crib)
    r.HTML(200, "list", All(db))
  })  

  // display the crib for a specific user
  m.Get("/:handle", func(params martini.Params, r render.Render, db *mgo.Database) {
    r.HTML(200, "display", Fetch(db, params["handle"]))    
  })    

  http.ListenAndServe(":8080", m)

}