package rest

import (
	"TokenServer/conf"
	"TokenServer/usecase"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HTTPHandler struct {
	usecase.UseCase
}

func Server(conf conf.Conf, useCase usecase.UseCase) {

	handler := HTTPHandler{useCase}

	mux := http.NewServeMux()
	mux.HandleFunc(conf.GetPath, handler.ServeGet)
	mux.HandleFunc(conf.RefreshPath, handler.ServeRefresh)
	if conf.CertFile != "" && conf.Keyfile != "" {
		fmt.Println("Starting HTTPS server")
		go http.ListenAndServeTLS(conf.HttpsPort, conf.CertFile, conf.Keyfile, mux)
	}
	fmt.Println("Starting HTTP server")
	http.ListenAndServe(conf.HttpPort, mux)
}

func (g HTTPHandler) ServeGet(w http.ResponseWriter, req *http.Request) {
	type Request struct {
		Guid string `json:"guid"`
	}

	type Response struct {
		RefreshToken string
		AccessToken  string
	}

	set := Request{}
	res := Response{}

	err := json.NewDecoder(req.Body).Decode(&set)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	if set.Guid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		g.UseCase.UserIp = strings.Split(ip, ",")[0]
	}

	if ip := req.Header.Get("X-Real-IP"); ip != "" && g.UseCase.UserIp == "" {
		g.UseCase.UserIp = strings.Split(ip, ",")[0]
	}

	if g.UseCase.UserIp == "" {
		g.UseCase.UserIp = strings.Split(req.RemoteAddr, ":")[0]
	}

	res.AccessToken, res.RefreshToken, err = g.CreateSession(set.Guid)
	if err != nil {
		//TODO: Обработать ошибку
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatal(err)
	}
}

func (r HTTPHandler) ServeRefresh(w http.ResponseWriter, req *http.Request) {
	type Request struct {
		RefreshToken string `json:"RefreshToken"`
		AccessToken  string `json:"AccessToken"`
		Ctx          context.Context
	}
	type Response struct {
		RefreshToken string
		AccessToken  string
	}

	set := Request{}
	res := Response{}

	err := json.NewDecoder(req.Body).Decode(&set)
	if err != nil {
		// FIXME: Если я передам невалидные данные - сервер скрашится?
		log.Fatal(err)
	}
	defer req.Body.Close()

	// FIXME: Куда? ctx принято передавать первым аргументам в функции, где он используется.
	// FIXME: Можно использовать req.Context() для большей гибкости
	set.Ctx = context.Background()

	// FIXME: Дублированный код с ServeGet - можно вынести в middleware либо
	// хотя бы в функцию.
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		// FIXME: Не мешать данные с конфигурацией - userIP передавать явно там, где он нужен.
		r.UseCase.UserIp = strings.Split(ip, ",")[0]
	}

	if ip := req.Header.Get("X-Real-IP"); ip != "" && r.UseCase.UserIp == "" {
		r.UseCase.UserIp = strings.Split(ip, ",")[0]
	}

	if r.UseCase.UserIp == "" {
		r.UseCase.UserIp = strings.Split(req.RemoteAddr, ":")[0]
	}

	res.AccessToken, res.RefreshToken, err = r.UseCase.RefreshSession(set.AccessToken, set.RefreshToken)
	if err != nil {
		// FIXME: Какой статус ошибки?
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Fatal(err)
		}
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		// FIXME: Точно 400?
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
