package main

import (
	"central-dummy-master/delivery"
	"central-dummy-master/delivery/router"
	"central-dummy-master/entities"
	"central-dummy-master/repository/databasefactory"
	"central-dummy-master/usecase"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kpango/glg"
)

func main() {
	go delivery.PrintoutObserver()

	//* CA Requirement
	if entities.UseCentralAuthority {
		er := usecase.OpenCentralSession()
		if er != nil {
			_ = glg.Error(er.Error())
			os.Exit(1)
		}
	}

	router.Start()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	LoadConfiguration(false)
	if os.Getenv("app.database_driver") != "" {
		PrepareDatabase()
	}

	go ReloadObserver()
}

func LoadConfiguration(isReload bool) {
	var er error
	if isReload {
		_ = glg.Log("Reloading configuration file...")
		er = godotenv.Overload(".env")
	} else {
		_ = glg.Log("Loading configuration file...")
		er = godotenv.Load(".env")
	}

	if er != nil {
		_ = glg.Error("Configuration file not found...")
		os.Exit(1)
	}

	//Opsi agar log utk level LOG, DEBUG, INFO dicatat atau tidak
	//Jika menggunakan docker atau dibuatkan service, log sudah dibuatkan, sehingga direkomendasikan
	//app log di set false
	appLog := os.Getenv("app.log")
	if appLog == "true" {
		log := glg.FileWriter("log/application.log", 0666)
		glg.Get().
			SetMode(glg.BOTH).
			AddLevelWriter(glg.LOG, log).
			AddLevelWriter(glg.DEBG, log).
			AddLevelWriter(glg.INFO, log)
	}

	//Untuk error akan selalu dicatat dalam file
	logEr := glg.FileWriter("log/application.err", 0666)
	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, logEr).
		AddLevelWriter(glg.WARN, logEr)
}

func PrepareDatabase() {
	var er error
	databasefactory.AppDb, er = databasefactory.GetDatabase()
	if er != nil {
		glg.Fatal(er.Error())
	}

	_ = glg.Log("Connecting to database...")
	if er = databasefactory.AppDb.Connect(); er != nil {
		_ = glg.Error("Connection to database failed : ", er.Error())
		os.Exit(1)
	}

	if er = databasefactory.AppDb.Ping(); er != nil {
		_ = glg.Error("Cannot ping database : ", er.Error())
		os.Exit(1)
	}

	_ = glg.Log("Database Connected")
}

func ReloadObserver() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP)

	func() {
		for {
			<-sign
			LoadConfiguration(true)
		}
	}()
}
