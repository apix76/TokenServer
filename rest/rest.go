package rest

import (
	"TokenServer/conf"
	"TokenServer/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type HTTPHandler struct {
	usecase.UseCase
}

func UseRealIP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var UserIp string
		if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			UserIp = strings.Split(ip, ",")[0]
		}
		if ip := r.Header.Get("X-Real-IP"); ip != "" && UserIp == "" {
			UserIp = strings.Split(ip, ",")[0]
		}

		if UserIp != "" {
			r.RemoteAddr = UserIp + ":0"
		}

		h.ServeHTTP(w, r)
	})
}

func Server(conf conf.Conf, useCase usecase.UseCase) {

	handler := HTTPHandler{useCase}

	mux := http.NewServeMux()
	mux.Handle(conf.GetPath, UseRealIP(http.HandlerFunc(handler.ServeGet)))
	mux.Handle(conf.RefreshPath, UseRealIP(http.HandlerFunc(handler.ServeRefresh)))
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
	httpStatus := 0
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

	UserIp := strings.Split(req.RemoteAddr, ":")[0]

	res.AccessToken, res.RefreshToken, err, httpStatus = g.CreateSession(set.Guid, UserIp)
	if err != nil {
		w.WriteHeader(httpStatus)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r HTTPHandler) ServeRefresh(w http.ResponseWriter, req *http.Request) {
	type Request struct {
		RefreshToken string `json:"RefreshToken"`
	}
	type Response struct {
		RefreshToken string
		AccessToken  string
	}
	httpStatus := 0
	set := Request{}
	res := Response{}

	err := json.NewDecoder(req.Body).Decode(&set)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	UserIp := strings.Split(req.RemoteAddr, ":")[0]

	res.AccessToken, res.RefreshToken, err, httpStatus = r.UseCase.RefreshSession(set.RefreshToken, UserIp)
	if err != nil {
		w.WriteHeader(httpStatus)
		w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
