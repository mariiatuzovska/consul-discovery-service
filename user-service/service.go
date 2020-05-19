package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
)

var (
	ServiceName       = "user-service"
	Version           = "0.0.1"
	DBFilePathSQLite3 = "./users.db"
)

func main() {
	app := cli.NewApp()
	app.Name = ServiceName
	app.Usage = "command line client"
	app.Description = ""
	app.Version = Version
	app.Authors = []cli.Author{cli.Author{Name: "Tuzovska Mariia"}}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "starting service via http",
			Action: func(c *cli.Context) error {
				srv, err := NewService()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(fmt.Sprintf("%s:%s", c.String("host"), c.String("port")))
				err = srv.Start(fmt.Sprintf("%s:%s", c.String("host"), c.String("port")))
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1",
				},
				&cli.StringFlag{
					Name:  "port",
					Value: "8080",
				},
			},
		},
	}
	app.Run(os.Args)
}

type Service struct {
	*gorm.DB
	*echo.Echo
}

type User struct {
	gorm.Model
	FirstName string `gorm:"not_null"`
	Name      string `gorm:"not_null"`
	Age       int    `gorm:"not_null"`
}

func NewService() (*Service, error) {
	if _, err := os.Open(DBFilePathSQLite3); err != nil {
		_, err = os.Create(DBFilePathSQLite3)
		if err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open("sqlite3", DBFilePathSQLite3)
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON;")
	db.AutoMigrate(&User{})

	srv := &Service{db, echo.New()}
	srv.HidePort = true
	srv.HideBanner = true

	srv.GET("/", srv.GetUsers)
	srv.POST("/", srv.CreateUser)
	srv.PATCH("/", srv.UpdateUser)
	srv.DELETE("/", srv.DeleteUser)

	return srv, nil
}

func (srv *Service) GetUsers(c echo.Context) error {
	query := new(User)
	if err := c.Bind(query); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	users := []User{}
	if srv.Find(&users, query).RecordNotFound() {
		c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, users)
}

func (srv *Service) CreateUser(c echo.Context) error {
	query := new(User)
	if err := c.Bind(query); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if query.FirstName == "" || query.Name == "" || query.Age < 1 {
		return c.NoContent(http.StatusBadRequest)
	}
	query.ID = 0
	user := new(User)
	srv.Model(User{}).Create(query).Last(user, query)
	return c.JSON(http.StatusOK, user)
}

func (srv *Service) UpdateUser(c echo.Context) error {
	query := new(User)
	if err := c.Bind(query); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if query.FirstName == "" || query.Name == "" || query.Age < 1 {
		return c.NoContent(http.StatusBadRequest)
	}
	user := new(User)
	if srv.Model(User{}).First(user, query.ID).RecordNotFound() {
		c.NoContent(http.StatusNotFound)
	}
	srv.Model(user).Update(query)
	srv.Model(User{}).First(user, query.ID)
	return c.JSON(http.StatusOK, user)
}

func (srv *Service) DeleteUser(c echo.Context) error {
	query := new(User)
	if err := c.Bind(query); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	user := new(User)
	if srv.Model(User{}).First(user, query.ID).RecordNotFound() {
		c.NoContent(http.StatusNotFound)
	}
	srv.Delete(user)
	return c.NoContent(http.StatusOK)
}
