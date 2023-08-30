package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"time"
	"user-segmentation-api/models"
)

type Database struct {
	driverName   string
	dbName       string
	dbConnection *sql.DB
}

func NewDatabase(driverName string) Database {
	var d Database
	d.driverName = driverName
	d.dbName = "sys"

	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS,
		HOST, DBNAME)

	db, err := sql.Open("mysql", URL)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE if not exists users (" +
		"id integer NOT NULL AUTO_INCREMENT," +
		"ttl varchar(50)," +
		"PRIMARY KEY (id) )")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE if not exists segments ( " +
		"name varchar(50) NOT NULL," +
		"percentage integer," +
		"PRIMARY KEY(name) )")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE if not exists segments_users ( " +
		"user_id integer NOT NULL," +
		"segment_name varchar(50) NOT NULL," +
		"PRIMARY KEY(user_id,segment_name) )")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE if not exists stats ( " +
		"user_id integer NOT NULL," +
		"segment_name varchar(50) NOT NULL," +
		"operation varchar(50) NOT NULL," +
		"timestamp datetime NOT NULL," +
		"PRIMARY KEY(user_id,segment_name,operation,timestamp) )")
	if err != nil {
		panic(err.Error())
	}

	d.dbConnection = db
	return d
}

func (d Database) Close() {
	err := d.dbConnection.Close()
	if err != nil {
		panic(err.Error())
	}
}

func (d Database) Vers() {

	var version string

	err2 := d.dbConnection.QueryRow("SELECT VERSION()").Scan(&version)

	if err2 != nil {
		panic(err2.Error())
	}

	fmt.Println(version)
}

func (d Database) FindSegment(s string) models.Segment {
	var segment models.Segment
	var intBuffer sql.NullInt32

	err := d.dbConnection.QueryRow("select * from segments where name = ?", s).Scan(&segment.Name, &intBuffer)
	if err == sql.ErrNoRows {
		return segment
	} else if err != nil {
		panic(err.Error())
	}

	segment.Percentage = int(intBuffer.Int32)

	return segment
}

func (d Database) CreateSegment(s models.Segment) {
	insert, err := d.dbConnection.Query("INSERT INTO segments VALUES (?,?)", s.Name, s.Percentage)

	if err != nil {
		panic(err.Error())
	}

	defer func(insert *sql.Rows) {
		err := insert.Close()
		if err != nil {
			panic(err.Error())
		}
	}(insert)

}

func (d Database) DeleteSegment(s string) {

	clear, err := d.dbConnection.Query("DELETE FROM segments_users WHERE segment_name=?", s)

	if err != nil {
		panic(err.Error())
	}

	defer func(clear *sql.Rows) {
		err := clear.Close()
		if err != nil {
			panic(err.Error())
		}
	}(clear)

	remove, err1 := d.dbConnection.Query("DELETE FROM segments WHERE name=?", s)

	if err1 != nil {
		panic(err1.Error())
	}

	defer func(remove *sql.Rows) {
		err1 := remove.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(remove)
}

func (d Database) GetAllSegments() []models.Segment {
	var segmentList []models.Segment

	segments, err := d.dbConnection.Query("SELECT * FROM segments")

	if err != nil {
		panic(err.Error())
	}

	defer func(segments *sql.Rows) {
		err1 := segments.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(segments)

	for segments.Next() {
		var segment models.Segment
		var intBuffer sql.NullInt32
		err2 := segments.Scan(&segment.Name, &intBuffer)
		if err2 != nil {
			panic(err2.Error())
		}
		segment.Percentage = int(intBuffer.Int32)
		segmentList = append(segmentList, segment)
	}

	return segmentList
}

func (d Database) FindUser(s int) models.User {
	var user models.User

	err := d.dbConnection.QueryRow("select id from users where id = ?", s).Scan(&user.ID)

	if err == sql.ErrNoRows {
		return user
	} else if err != nil {
		panic(err.Error())
	}

	d.RemoveTimeoutUsers()

	return user
}

func (d Database) RemoveTimeoutUsers() {

	users, err := d.dbConnection.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}

	defer func(users *sql.Rows) {
		err1 := users.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(users)

	for users.Next() {
		var user models.User
		var stringBuffer sql.NullString
		err2 := users.Scan(&user.ID, &stringBuffer)
		if err2 != nil {
			panic(err2.Error())
		}
		ttl := stringBuffer.String
		if ttl != "" {
			ttlList := strings.Split(ttl, "-")
			now := time.Now()
			ttlYear, _ := strconv.Atoi(ttlList[0])
			ttlMonth, _ := strconv.Atoi(ttlList[1])
			ttlDay, _ := strconv.Atoi(ttlList[2])
			if ttlYear == now.Year() && ttlMonth == int(now.Month()) && ttlDay == now.Day() {
				d.DeleteUser(user.ID)
			}
		}

	}

}

func (d Database) FindUserSegments(s int) []models.Segment {
	var segmentList []models.Segment

	segments, err := d.dbConnection.Query("SELECT segments.* "+
		"FROM (SELECT * FROM users WHERE id=?) u "+
		"INNER JOIN segments_users "+
		"ON segments_users.user_id = u.id "+
		"INNER JOIN segments "+
		"ON segments_users.segment_name = segments.name", s)

	if err != nil {
		panic(err.Error())
	}

	defer func(segments *sql.Rows) {
		err1 := segments.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(segments)

	for segments.Next() {
		var segment models.Segment
		var intBuffer sql.NullInt32
		err2 := segments.Scan(&segment.Name, &intBuffer)
		if err2 != nil {
			panic(err2.Error())
		}
		segment.Percentage = int(intBuffer.Int32)
		segmentList = append(segmentList, segment)
	}

	return segmentList
}

func (d Database) ChangeUserSegments(id int, addList []string, removeList []string) {

	for _, s := range removeList {
		remove, err := d.dbConnection.Query("DELETE FROM segments_users WHERE user_id=? AND segment_name=?", id, s)

		if err != nil {
			panic(err.Error())
		}

		err1 := remove.Close()

		if err1 != nil {
			panic(err1.Error())
		}

		removeStats, err2 := d.dbConnection.Query("INSERT INTO stats VALUES(?,?,?,?)", id, s, "remove", time.Now())

		if err2 != nil {
			panic(err2.Error())
		}

		err3 := removeStats.Close()

		if err3 != nil {
			panic(err3.Error())
		}
	}

	for _, s := range addList {
		add, err := d.dbConnection.Query("INSERT INTO segments_users VALUES(? , ?)", id, s)

		if err != nil {
			panic(err.Error())
		}

		err1 := add.Close()

		if err1 != nil {
			panic(err1.Error())
		}

		addStats, err2 := d.dbConnection.Query("INSERT INTO stats VALUES(?,?,?,?)", id, s, "add", time.Now())

		if err2 != nil {
			panic(err2.Error())
		}

		err3 := addStats.Close()

		if err3 != nil {
			panic(err3.Error())
		}
	}

}

func (d Database) CreateUser(user models.User) {
	var insert *sql.Rows
	var err error

	if user.TTL == 0 {
		insert, err = d.dbConnection.Query("INSERT INTO users(id) VALUES (?)", user.ID)
	} else {
		insert, err = d.dbConnection.Query("INSERT INTO users(id,ttl) VALUES (?,?)", user.ID, time.Now().AddDate(0, 0, user.TTL).Format(time.DateOnly))
	}

	if err != nil {
		panic(err.Error())
	}

	defer func(insert *sql.Rows) {
		err1 := insert.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(insert)

}

func (d Database) DeleteUser(s int) {

	clear, err := d.dbConnection.Query("DELETE FROM segments_users WHERE user_id=?", s)

	if err != nil {
		panic(err.Error())
	}

	defer func(clear *sql.Rows) {
		err := clear.Close()
		if err != nil {
			panic(err.Error())
		}
	}(clear)

	remove, err1 := d.dbConnection.Query("DELETE FROM users WHERE id=?", s)

	if err1 != nil {
		panic(err1.Error())
	}

	defer func(remove *sql.Rows) {
		err1 := remove.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(remove)
}

func (d Database) GetStats(month int, year int) []models.StatString {
	stats, err := d.dbConnection.Query("SELECT * FROM stats")

	if err != nil {
		panic(err.Error())
	}

	defer func(stats *sql.Rows) {
		err1 := stats.Close()
		if err1 != nil {
			panic(err1.Error())
		}
	}(stats)

	var allStats []models.StatString

	for stats.Next() {
		var stat models.StatString
		err2 := stats.Scan(&stat.UserID, &stat.SegmentName, &stat.Operation, &stat.Timestamp)
		if err2 != nil {
			panic(err2.Error())
		}
		timestamp := stat.Timestamp
		if int(timestamp.Month()) == month && timestamp.Year() == year {
			allStats = append(allStats, stat)
		}
	}
	return allStats
}
