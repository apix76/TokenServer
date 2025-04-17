package main

import (
	"TokenServer/conf"
	"TokenServer/rest"
	"TokenServer/service/db"
	"TokenServer/usecase"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	con := conf.NewConf()
	useCase := usecase.UseCase{} //{
	//	Token: token.Token{
	//		ExpTimeAccess:  con.ExpTimeAccess,
	//		ExpTimeRefresh: con.ExpTimeRefresh,
	//	},
	//}

	var err error
	useCase.DB, err = db.New(con.PgsqlNameServe)
	if err != nil {
		log.Fatal(err)
	}

	//useCase.Token.PrivateKey, useCase.Token.PublicKey, err = token.TokenKey(con.PrivateKey, con.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	rest.Server(con, useCase)

}
