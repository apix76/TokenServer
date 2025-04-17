package conf

import (
	"encoding/json"
	"log"
	"os"
)

type Conf struct {
	CertFile       string
	Keyfile        string
	HttpPort       string
	HttpsPort      string
	GetPath        string
	RefreshPath    string
	PgsqlNameServe string
}

func NewConf() Conf {
	con := Conf{}

	FileCon, err := os.Open("config.cfg")
	if err != nil {
		log.Fatal(err)
	}
	defer FileCon.Close()

	err = json.NewDecoder(FileCon).Decode(&con)
	if err != nil {
		log.Fatal(err)
	}

	return con
}
