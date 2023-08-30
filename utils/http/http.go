package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"user-segmentation-api/utils/database"
	"user-segmentation-api/utils/models"
)

type Server struct {
	db      database.Database
	address string
}

func New(c string) Server {
	var s Server
	s.address = c
	return s
}

func (s Server) StartServer() {
	s.db = database.New("mysql")
	defer s.db.Close()

	s.db.Vers()

	router := gin.Default()
	router.POST("/segments", s.postSegments)
	router.POST("/users", s.postUsers)
	router.GET("/segments", s.getUserSegmentsByID)
	router.DELETE("/segments", s.deleteSegmentByName)
	router.DELETE("/users", s.deleteUserByID)
	router.POST("/users/:id", s.changeUserSegmentsByID)

	err := router.Run(s.address)
	if err != nil {
		panic(err.Error())
	}
}

func (s Server) postSegments(c *gin.Context) {
	var newSegment models.NameReq

	if err := c.BindJSON(&newSegment); err != nil {
		panic(err.Error())
	}

	segment := s.db.FindSegment(newSegment.Name)
	if segment.Name != "" {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Segment already exists"})
		return
	}

	s.db.CreateSegment(newSegment.Name)
	c.IndentedJSON(http.StatusCreated, newSegment)
}

func (s Server) postUsers(c *gin.Context) {
	var newUser models.IDReq

	if err := c.BindJSON(&newUser); err != nil {
		panic(err.Error())
	}

	segment := s.db.FindUser(newUser.ID)
	if segment.ID != 0 {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	s.db.CreateUser(newUser.ID)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func (s Server) getUserSegmentsByID(c *gin.Context) {
	var newID models.IDReq

	if err := c.BindJSON(&newID); err != nil {
		return
	}

	id := newID.ID
	user := s.db.FindUser(id)
	if user.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	segments := s.db.FindUserSegments(id)
	c.IndentedJSON(http.StatusOK, segments)
}

func (s Server) deleteSegmentByName(c *gin.Context) {
	var newName models.NameReq

	if err := c.BindJSON(&newName); err != nil {
		return
	}
	name := newName.Name

	s.db.DeleteSegment(name)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Deleted segment"})
}

func (s Server) deleteUserByID(c *gin.Context) {
	var newID models.IDReq

	if err := c.BindJSON(&newID); err != nil {
		return
	}
	id := newID.ID

	s.db.DeleteUser(id)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Deleted user"})
}

func (s Server) changeUserSegmentsByID(c *gin.Context) {
	idString := c.Param("id")
	id, err1 := strconv.Atoi(idString)

	if err1 != nil {
		panic(err1.Error())
	}

	var changeSegments models.ChangeReq

	if err := c.BindJSON(&changeSegments); err != nil {
		return
	}

	for _, seg := range changeSegments.Add {
		segment := s.db.FindSegment(seg)
		if segment.Name == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Segment not found"})
			return
		}
	}

	for _, seg := range changeSegments.Remove {
		segment := s.db.FindSegment(seg)
		if segment.Name == "" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Segment not found"})
			return
		}
	}

	user := s.db.FindUser(id)
	if user.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	s.db.ChangeUserSegments(id, changeSegments.Add, changeSegments.Remove)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Changed user segments"})
}
