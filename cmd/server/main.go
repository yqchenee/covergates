package main

import (
	"os"

	"github.com/code-devel-cover/CodeCover/config"
	"github.com/code-devel-cover/CodeCover/core"
	"github.com/code-devel-cover/CodeCover/routers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	// load sqlite driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func connectDatabase() *gorm.DB {
	x, _ := gorm.Open("sqlite3", "core.db")
	return x

}

func Run(c *cli.Context) error {
	config := &config.Config{
		Server: config.Server{
			Addr: "http://34.70.126.184:3000",
			Base: "",
		},
		Github: config.Github{
			ClientID:     "a150e893154bafde8a00",
			ClientSecret: "59a3f97b6e7569d0b6898bc5fb2e84f93e64113d",
			Server:       "https://github.com",
			APIServer:    "https://api.github.com",
		},
		Gitea: config.Gitea{
			ClientID:     "c8c6a2cc-f948-475c-8663-f420c8fc15ab",
			ClientSecret: "J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE=",
			Server:       "http://localhost:3000",
			SkipVerity:   true,
		},
	}
	db := connectDatabase()
	app, err := InitializeApplication(config, db)
	if err != nil {
		return err
	}
	r := gin.Default()
	app.routers.RegisterRoutes(r)
	app.db.Migrate()
	r.Run(":3000")
	return nil
}

func main() {
	log.SetReportCaller(true)
	app := &cli.App{
		Name:   "codecover",
		Action: Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type application struct {
	routers *routers.Routers
	db      core.DatabaseService
}

func newApplication(
	routers *routers.Routers,
	db core.DatabaseService,
) application {
	return application{
		routers: routers,
		db:      db,
	}
}
