package main

import (
	"net/http"
	"time"

	"testtask/api/types"

	"github.com/gin-gonic/gin"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"
)

type Data struct {
	Id   bson.ObjectId `form:"id" bson:"_id,omitempty"`
	Data string        `form:"data" bson:"data"`
}

type MongoDB struct {
	Host             string
	Port             string
	Addrs            string
	Database         string
	EventTTLAfterEnd time.Duration
	StdEventTTL      time.Duration
	Info             *mgo.DialInfo
	Session          *mgo.Session
}

func (mongo *MongoDB) GetData() (dates []Data, err error) { // {{{
	session := mongo.Session.Clone()
	defer session.Close()

	err = session.DB(mongo.Database).C("Data").Find(bson.M{}).All(&dates)
	return dates, err
}

func (mongo *MongoDB) PostData(data *Data) (err error) { // {{{
	session := mongo.Session.Clone()
	defer session.Close()

	err = session.DB(mongo.Database).C("Data").Insert(&data)
	return err
}

func (mongo *MongoDB) CreateUser(user *types.User) (err error) {
	session := mongo.Session.Clone()
	defer session.Close()
	err = session.DB(mongo.Database).C("Data").Insert(&user)
	return err
}

func MiddleDB(mongo *MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := mongo.SetSession()
		if err != nil {
			c.Abort()
		} else {
			c.Set("mongo", mongo)
			c.Next()
		}
	}
}

func (mongo *MongoDB) SetSession() (err error) {
	mongo.Session, err = mgo.DialWithInfo(mongo.Info)

	if err != nil {
		//mongo.Session, err = mgo.Dial(mongo.Host)
		//if err != nil {
		// fmt.Printf("settts")
		// log.Fatal(err)
		return err
		//}
	}
	return err
}
func getData(c *gin.Context) { // {{{
	mongo, ok := c.Keys["mongo"].(*MongoDB)
	if !ok {
		c.JSON(400, gin.H{"message": "can't reach db", "body": nil})
	}

	session := mongo.Session.Clone()
	defer session.Close()

	var users []types.User

	err := session.DB(mongo.Database).C("User").Find(bson.M{}).All(&users)
	if err != nil {
		c.JSON(400, gin.H{"message": "error post to db", "body": nil, "err": err})
		return
	}
	//data, err := mongo.GetData()
	// fmt.Printf(err)
	// fmt.Printf("\ndata: %v, ok: %v\n", data, ok)
	//if err != nil {
	//	c.JSON(400, gin.H{"message": "can't get data from database", "body": nil})
	//} else {
	c.JSON(200, gin.H{"message": "get data sucess", "body": users})
	//}
} // }}}

func postData(c *gin.Context) { // {{{
	mongo, ok := c.Keys["mongo"].(*MongoDB)
	if !ok {
		c.JSON(400, gin.H{"message": "can't connect to db", "body": nil})
		return
	}

	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := validator.Validate(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session := mongo.Session.Clone()
	defer session.Close()
	err := session.DB(mongo.Database).C("User").Insert(&user)
	if err != nil {
		c.JSON(400, gin.H{"message": "error post to db", "body": nil, "err": err})
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "error post to db", "body": nil, "err": err})
		return
	}
	c.JSON(200, gin.H{"message": "post data sucess"})
}

func main() {
	mongo := MongoDB{}

	mongo.Info = &mgo.DialInfo{
		Addrs:    []string{"localhost: 27017"},
		Timeout:  60 * time.Second,
		Database: "context",
		Username: "test",
		Password: "test",
	}

	router := gin.Default()
	router.Use(MiddleDB(&mongo))
	router.POST("/data", postData)
	router.GET("/data", getData)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
