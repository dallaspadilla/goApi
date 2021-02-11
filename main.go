package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
	"net/http"
	"time"
)

type User struct {
	Id     int64
	Name   string `json:"content" binding:"required"`
	Emails []string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.Id, u.Name, u.Emails)
}

type Bulletins struct {
	Author    string    `json:"author" binding: "required`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	//Id     int64
	//AuthorId int64
	//Author   *User `pg:"rel:has-one"`
	//Content string `json:"content" binding:"required"`
	//CreatedAt time.Time `json:"created_at"`
}

//
//func (s Bulletins) String() string {
//	return fmt.Sprintf("Story<%d %s %s>", s.Id, s.Content, s.Author)
//}

//type Bulletin struct {
//	Author string `json:"author" binding: "required`
//	Content string `json:"content" binding:"required"`
//	CreatedAt time.Time `json:"created_at"`
//}

var db *sql.DB

//func GetBulletins() ([]Bulletins, error) {
//
//	const q = `SELECT author, content, created_at FROM bulletins ORDER BY created_at DESC LIMIT 100`
//
//	rows, err := db.Query(q)
//	if err != nil {
//		return nil, err
//	}
//
//	results := make([]Bulletins, 0)
//
//	for rows.Next() {
//		var author string
//		var content string
//		var createAt time.Time
//		err = rows.Scan(&author, &content, &createAt)
//		if err != nil {
//			return nil, err
//		}
//		results = append(results, Bulletin{ Author: author, Content: content, CreatedAt: createAt})
//	}
//
//	return results, nil
//}

func GetStorys() ([]Bulletins, error) {

	const q = `SELECT author, content, createdAt FROM bulletins ORDER BY createdAt DESC LIMIT 100`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]Bulletins, 0)

	for rows.Next() {
		var author string
		var content string
		var createAt time.Time
		err = rows.Scan(&author, &content, &createAt)
		if err != nil {
			return nil, err
		}
		results = append(results, Bulletins{Author: author, Content: content, CreatedAt: createAt})
	}

	return results, nil
}

func AddBulletin(bulletin Bulletins) error {
	const q = `INSERT INTO bulletins(author, content, createdAt) Values ($1, $2, $3)`
	_, err := db.Exec(q, bulletin.Author, bulletin.Content, bulletin.CreatedAt)
	return err
}

func main() {

	r := gin.Default()
	r.GET("/board", func(context *gin.Context) {
		results, err := GetStorys()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
			return
		}
		context.JSON(http.StatusOK, results)
	})
	r.POST("/board", func(context *gin.Context) {
		var b Bulletins

		if context.Bind(&b) == nil {
			b.CreatedAt = time.Now()
			if err := AddBulletin(b); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
				return
			}
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
	// dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable\", DbHost, DbUser, DbPassword, DbName")
	// db, err = sql.Open("postgres", dbInfo)
	db := pg.Connect(&pg.Options{
		User:     "dallas",
		Password: "password",
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	log.Println("running..")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Bulletins)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
