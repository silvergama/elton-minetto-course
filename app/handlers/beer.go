package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/silvergama/studying-golang/elton-minetto-course/core/beer"
	"github.com/urfave/negroni"
)

// a função recebe como terceiro parâmetro a interface ou seja, ela pode receber
// qualquer coisa que implemente a interface, isso é muito útil para escrevermos
// testes, ou podermos substituir toda a implementação da regra de negócios
func MakeBeerHandlers(r *mux.Router, n *negroni.Negroni, service beer.UseCase) {
	r.Handle("/v1/beer", n.With(
		negroni.Wrap(getAllBeer(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/beer/{id}", n.With(
		negroni.Wrap(getBeer(service)),
	)).Methods("GET", "OPTIONS")

	r.Handle("/v1/beer", n.With(
		negroni.Wrap(storeBeer(service)),
	)).Methods("POST", "OPTIONS")

	/*
		r.Handle("/v1/beer/{id}", n.With(
			negroni.Wrap(updateBeer(service)),
		)).Methods("PUT", "OPTIONS")

		r.Handle("/v1/beer/{id}", n.With(
			negorni.Wrap(removeBeer(service)),
		)).Methods("DELETE", "OPTIONS")
	*/
}

func getAllBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		all, err := service.GetAll()
		if err != nil {
			w.Write(formatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//vamos converter o resultado em JSON e gerar a resposta
		err = json.NewEncoder(w).Encode(all)
		if err != nil {
			w.Write(formatJSONError("Erro ao tentar converter para JSON"))
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func getBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// vamos pegar o ID da URL
		// na definição do protocolo http, os parâmetros são enviados no formato
		// de texto, por isso converter em int64
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.Write(formatJSONError(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := service.Get(id)
		if err != nil {
			w.Write(formatJSONError(err.Error()))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// vamos converter o resultado para JSON e gerar e resposa
		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			w.Write(formatJSONError("Erro ao tentar converter JSON"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func storeBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// vamos pegar os dados enviados pelo usuário via body
		var b beer.Beer
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			w.Write(formatJSONError(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Precisamos validar os dados antes de salvar na base de dados.
		// Pergunta: Como fazer isso?
		err = service.Store(&b)
		if err != nil {
			w.Write(formatJSONError(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}
