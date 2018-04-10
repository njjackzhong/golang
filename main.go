//package main
package main

import (
	"encoding/json"
	"net/http"

	"github.com/donnie4w/go-logger/logger"
	"github.com/elgs/jsonql"
	"github.com/gorilla/mux"
)

type Query struct {
	Where    string `json:"where,omitempty"`
	Messages string `json:"messages,omitempty"`
	Err      error  `json:"error,omitempty"`
}

type Judge struct {
	Where   string      `json:"where,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Err     error       `json:"error,omitempty"`
}

type JudgeResponse struct {
	Message string `json:"message"` 
	Code    int    `json:"code"`
	Valid   bool   `json:"valid"` 
}

func SelectMessage(w http.ResponseWriter, req *http.Request) {

	var query Query
	_ = json.NewDecoder(req.Body).Decode(&query)

	//messages := req.URL.Query().Get("messages")
	//where := req.URL.Query().Get("where")

	//logger.Info("请求参数(消息组)  messages: " + query.Messages)
	//logger.Info("请求参数(过滤条件)where: " + query.Where)

	parser, err := jsonql.NewStringQuery(query.Messages)

	if err == nil {
		selectedMessages, err := parser.Query(query.Where)

		//logger.Info("返回数据:",selectedMessages,err)

		//if err == nil {
		//json.NewEncoder(w).Encode(true)

		query.Err = err

		json.NewEncoder(w).Encode(selectedMessages)
		//}
	}
}

func JudgeMessage(w http.ResponseWriter, req *http.Request) {
	var judge Judge

	_ = json.NewDecoder(req.Body).Decode(&judge)

	logger.Info("请求参数: ", judge)

	parser := jsonql.NewQuery(judge.Message)

	message, err := parser.Query(judge.Where)

	var judgeResponse JudgeResponse
	if message == nil {
		if err == nil {
			judgeResponse.Message = ""
		} else {
			judgeResponse.Message = err.Error()
		}
		judgeResponse.Code = -1
		judgeResponse.Valid = false
	} else {
		judgeResponse.Message = ""
		judgeResponse.Code = 0
		judgeResponse.Valid = true
	}

	logger.Info("返回值 : ", judgeResponse)

	//w.Header.set("", "")

	//http.Header.Add()
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(judgeResponse)

}

func main() {
	router := mux.NewRouter()

	logger.SetRollingDaily("d://logs//JudgeService", "Judge.txt")
	logger.SetLevel(logger.DEBUG)

	router.HandleFunc("/select", SelectMessage).Methods("POST")
	router.HandleFunc("/judge", JudgeMessage).Methods("POST")

	logger.Info("Starting JudgeService, Listen at port :12345"," Version:v1.0-20171207-snapshot")

	http.ListenAndServe(":12345", router)

}
