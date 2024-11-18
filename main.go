package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"

	"github.com/jorgerojas26/lazysql/app"
	"github.com/jorgerojas26/lazysql/components"
	"github.com/jorgerojas26/lazysql/drivers"
	"github.com/jorgerojas26/lazysql/helpers"
	"github.com/jorgerojas26/lazysql/helpers/logger"
	"github.com/jorgerojas26/lazysql/models"
)

var version = "dev"

func main() {
	rawLogLvl := flag.String("loglvl", "info", "Log level")
	logFile := flag.String("logfile", "", "Log file")
	flag.Parse()

	logLvl, parseError := logger.ParseLogLevel(*rawLogLvl)
	if parseError != nil {
		panic(parseError)
	}
	logger.SetLevel(logLvl)

	if *logFile != "" {
		fileError := logger.SetFile(*logFile)
		if fileError != nil {
			panic(fileError)
		}
	}

	logger.Info("Starting LazySQL...", nil)

	mysqlError := mysql.SetLogger(log.New(io.Discard, "", 0))
	if mysqlError != nil {
		panic(mysqlError)
	}

	// check if "version" arg is passed
	argsWithProg := os.Args
	if len(argsWithProg) > 1 {
		switch argsWithProg[1] {
		case "version":
			println("LazySQL version: ", version)
			os.Exit(0)
		default:
			connectionString := argsWithProg[1]
			parsed, err := helpers.ParseConnectionString(connectionString)
			if err != nil {
				fmt.Printf("Could not parse connection string: %s\n", err)
				os.Exit(1)
			}
			connection := models.Connection{
				Name:     connectionString,
				Provider: parsed.Driver,
				DBName:   connectionString,
				URL:      connectionString,
			}
			newDbDriver := &drivers.SQLite{}
			err = newDbDriver.Connect(connection.URL)
			if err != nil {
				fmt.Printf("Could not connect to database %s: %s\n", connectionString, err)
				os.Exit(1)
			}
			components.MainPages.AddPage(connection.URL, components.NewHomePage(connection, newDbDriver).Flex, true, true)
		}
	}

	if err := app.App.
		SetRoot(components.MainPages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}
