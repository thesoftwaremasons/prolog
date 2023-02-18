package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type httpServer struct {
	log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{log: NewLog()}
}

type ProduceRequest struct {
	record Record `json:"record"`
}
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	offset uint64 `json:"offset"`
}
type ConsumeResponse struct {
	Record Record `json:"record"`
}

func NewHTTPServer(addr string) *http.Server {
	httpSvr := newHTTPServer()
	router := mux.NewRouter()
	router.HandleFunc("/6", httpSvr.handleProduce).Methods("POST")
	router.HandleFunc("/", httpSvr.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func (c *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {

	//decoding data
	var produceRequest ProduceRequest

	err := json.NewDecoder(r.Body).Decode(&produceRequest)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := c.log.Append(produceRequest.record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{
		offset,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (c *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var consumeRequest ConsumeRequest

	err := json.NewDecoder(r.Body).Decode(&consumeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	record, err := c.log.Read(consumeRequest.offset)
	if err == ErrOfsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
