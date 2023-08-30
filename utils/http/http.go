package http

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
	"user-segmentation-api/models"
	"user-segmentation-api/utils/database"
)

type Server struct {
	db          database.Database
	address     string
	segmentList []models.SegmentCounter
}

func NewServer(c string) Server {
	var s Server
	s.address = c
	return s
}

func (s Server) StartServer() {
	s.db = database.NewDatabase("mysql")
	defer s.db.Close()

	s.db.Vers()

	segments := s.db.GetAllSegments()

	for _, segment := range segments {
		newCounter := models.NewCounter(segment)
		s.segmentList = append(s.segmentList, newCounter)
		print(newCounter.Segment.Name, newCounter.Proportion)
	}

	router := gin.Default()
	router.POST("/segments", s.postSegments)
	router.POST("/users", s.postUsers)
	router.GET("/segments", s.getUserSegmentsByID)
	router.DELETE("/segments", s.deleteSegmentByName)
	router.DELETE("/users", s.deleteUserByID)
	router.POST("/users/:id", s.changeUserSegmentsByID)
	router.GET("/stats", CSVExport(s))

	err := router.Run(s.address)
	if err != nil {
		panic(err.Error())
	}
}

func getStats(stats []models.StatString) (*os.File, error) {
	var records [][]string

	for _, stat := range stats {
		statList := []string{strconv.Itoa(stat.UserID), stat.SegmentName, stat.Operation, stat.Timestamp.Format(time.DateTime)}
		records = append(records, statList)
	}

	file, err := os.Create("stats.csv")
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err2 := file.Close()
		if err2 != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CSVExport(s Server) gin.HandlerFunc {
	return func(c *gin.Context) {

		var newStats models.StatsReq

		if err := c.BindJSON(&newStats); err != nil {
			panic(err.Error())
		}

		statsList := s.db.GetStats(newStats.Month, newStats.Year)

		_, err := getStats(statsList)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stat"})
			return //stop it on error
		}
		c.FileAttachment("./stats.csv", "stats.csv")
		c.Writer.Header().Set("attachment", "filename=stats.csv")
	}
}

func (s Server) postSegments(c *gin.Context) {
	var newSegment models.Segment

	if err := c.BindJSON(&newSegment); err != nil {
		panic(err.Error())
	}

	segment := s.db.FindSegment(newSegment.Name)
	if segment.Name != "" {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Segment already exists"})
		return
	}

	s.db.CreateSegment(newSegment)

	counter := models.NewCounter(newSegment)
	s.segmentList = append(s.segmentList, counter)

	c.IndentedJSON(http.StatusCreated, newSegment)
}

func (s Server) postUsers(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		panic(err.Error())
	}

	segment := s.db.FindUser(newUser.ID)
	if segment.ID != 0 {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	}

	s.db.CreateUser(newUser)

	segments := s.db.GetAllSegments()

	s.segmentList = UpdateCounters(s.segmentList, segments)

	for i, counter := range s.segmentList {
		print(counter.Count, counter.Proportion)
		s.segmentList[i].Count = counter.Count + 1
		if s.segmentList[i].Count <= counter.Proportion {
			var singleList []string
			singleList = append(singleList, counter.Segment.Name)
			var emptyList []string
			s.db.ChangeUserSegments(newUser.ID, singleList, emptyList)
		}
		if s.segmentList[i].Count == 100 {
			s.segmentList[i].Count = 0
		}
	}

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

func UpdateCounters(list1 []models.SegmentCounter, list2 []models.Segment) []models.SegmentCounter {
	for _, counter := range list2 {
		found := false
		for _, oldCounter := range list1 {
			if counter.Name == oldCounter.Segment.Name {
				found = true
			}
		}
		if !found {
			list1 = append(list1, models.NewCounter(counter))
		}
	}
	return list1
}
