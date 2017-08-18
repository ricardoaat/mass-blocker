package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ricardoaat/mass-blocker/blocker"
	"github.com/ricardoaat/mass-blocker/config"
	"github.com/ricardoaat/mass-blocker/db"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func logInit() {

	t := time.Now()
	logfile := config.Conf.Path.LogPath + fmt.Sprintf("massblocker_%s.log", t.Format("20060102T150405"))
	//logfile := config.Conf.Path.LogPath + "notif.log"
	fmt.Println("Loging to " + logfile)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(new(prefixed.TextFormatter))
	log.AddHook(lfshook.NewHook(lfshook.PathMap{
		log.DebugLevel: logfile,
		log.InfoLevel:  logfile,
		log.ErrorLevel: logfile,
		log.WarnLevel:  logfile,
		log.PanicLevel: logfile,
	}))

}

var (
	version   string
	buildDate string
)

func main() {

	err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal("Failed to load config.toml ", err)
		fmt.Println(err)
	}
	logInit()
	log.Info(fmt.Sprintf("Version: %s Build Date: %s", version, buildDate))
	log.Info("--------------Init program--------------")
	log.Debug("Loaded configuration " + fmt.Sprint(config.Conf))
	v := flag.Bool("v", false, "Returns the binary version and built date info")
	flag.Parse()
	if !*v {
		db.InitDatabases()
		blocker.Start()
		db.CloseDatabases()
	}
}
