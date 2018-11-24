package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {
	// Create instance
	e := echo.New()

	// mongo
	//	host: "host" -> yaml
	// MONGO_HOST -> env
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	mongoHost := viper.GetString("mongo.host")
	mongoUser := viper.GetString("mongo.user")
	mongoPass := viper.GetString("mongo.pass")
	port := ":" + viper.GetString("port")

	connString := fmt.Sprintf("%s:%s@%s", mongoUser, mongoPass, mongoHost)
	// "root:root@127.0.0.1" - local
	session, err := mgo.Dial(connString)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	// Declare struct and put session into handler
	// TODO: ??
	h := &handler{m: session}

	// Middleware
	// Log will written very handler has run
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Apis
	e.GET("/todos", h.list)
	e.GET("/todos/:id", h.view)
	e.PUT("/todos/:id", h.done)
	e.POST("/todos", h.create)
	e.DELETE("/todos/:id", h.delete)

	// Start server
	e.Logger.Fatal(e.Start(port))
}

// Todo struct
// ------------- `json:"var"` -> decode in json form it will assigned by var
// mongo received bson
type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id" `
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

type handler struct {
	m *mgo.Session
}

// echo.Context is interface
func (h *handler) create(c echo.Context) error {
	// Copy session to use and close but main session still exist
	session := h.m.Copy()
	defer session.Close()

	var t todo
	// Unmarshal context to todo
	if err := c.Bind(&t); err != nil {
		return err
	}

	t.ID = bson.NewObjectId()

	col := session.DB("workshop").C("todos")

	// Insert data to database
	if err := col.Insert(t); err != nil {
		return err
	}

	// Receive json as params -> for test
	return c.JSON(http.StatusOK, t)
}

// (h handler) -> copy handler
// (h *handler) -> point to exist handler
func (h *handler) list(c echo.Context) error {
	// Copy session to use and close but main session still exist
	session := h.m.Copy()
	defer session.Close()

	var ts []todo
	col := session.DB("workshop").C("todos")

	// Find all and assign to ts with type []todo
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ts)
}

func (h *handler) view(c echo.Context) error {
	// Copy session to use and close but main session still exist
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	var t todo
	col := session.DB("workshop").C("todos")

	// Find all and assign to t with type todo
	// FindId with bson param
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

// Get then update
func (h *handler) done(c echo.Context) error {
	// Copy session to use and close but main session still exist
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	var t todo
	col := session.DB("workshop").C("todos")

	// Find all and assign to t with type todo
	// FindId with bson param
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}

	t.Done = true
	// Update by id
	if err := col.UpdateId(id, t); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (h *handler) delete(c echo.Context) error {
	// Copy session to use and close but main session still exist
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))

	col := session.DB("workshop").C("todos")

	if err := col.RemoveId(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}
