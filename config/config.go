package config

//Conf global configure instance
var Conf config

type config struct {
	Path      path
	Services  map[string]service
	Databases map[string]Database
}

type path struct {
	LogPath string
}

type service struct {
	Serviceaddrs    string
	Serverport      string
	ServiceEndpoint string
}

//Database data base type with its connection parameters
type Database struct {
	Host     string
	Port     string
	DBname   string
	Username string
	Password string
}
