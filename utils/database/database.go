package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"user-segmentation-api/utils/models"
)

type Database struct {
	driverName   string
	dbName       string
	dbConnection *sql.DB
}

func New(driverName string) Database {
	var d Database
	d.driverName = driverName
	d.dbName = "sys"
	db, err := sql.Open("mysql", os.Getenv("DATA_SOURCE"))
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE if not exists users (" +
		"id integer NOT NULL AUTO_INCREMENT," +
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

func (d Database) CreateSegment(s string) {
	insert, err := d.dbConnection.Query("INSERT INTO segments(name) VALUES (?)", s)

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

func (d Database) FindUser(s int) models.User {
	var user models.User

	err := d.dbConnection.QueryRow("select * from users where id = ?", s).Scan(&user.ID)

	if err == sql.ErrNoRows {
		return user
	} else if err != nil {
		panic(err.Error())
	}

	return user
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
	}

	for _, s := range addList {
		add, err := d.dbConnection.Query("INSERT INTO segments_users VALUES(?,?)", id, s)

		if err != nil {
			panic(err.Error())
		}

		err1 := add.Close()

		if err1 != nil {
			panic(err1.Error())
		}
	}
}

func (d Database) CreateUser(s int) {
	insert, err := d.dbConnection.Query("INSERT INTO users VALUES (?)", s)

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
