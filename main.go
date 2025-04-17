package main

import (
	"TokenServer/conf"
	"TokenServer/rest"
	"TokenServer/service/db"
	"TokenServer/service/token"
	"TokenServer/usecase"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	con := conf.NewConf()
	useCase := usecase.UseCase{
		Token: token.Token{
			ExpTimeAccess:  con.ExpTimeAccess,
			ExpTimeRefresh: con.ExpTimeRefresh,
		},
	}

	var err error
	useCase.DB, err = db.New(con.PgsqlNameServe)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	rest.Server(con, useCase)

}
