package test

import (
	"console-application-service-age/internal/handlers"
	"console-application-service-age/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	urlTest1 = "/test1"
	urlTest2 = "/test2"
	urlTest3 = "/test3"
	network  = "tcp"
	address  = "127.0.0.1:8081"
)

type handler struct {
	repository Repository
}

func NewHandler(repository Repository) handlers.Handler {
	return &handler{repository: repository}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, urlTest1, h.IncrementVal)
	router.HandlerFunc(http.MethodPost, urlTest2, h.Signature)
	router.HandlerFunc(http.MethodPost, urlTest3, h.Multiplication)
}

func (h *handler) IncrementVal(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var test1 test1
	err = json.Unmarshal(requestBody, &test1)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	newValue, err := h.repository.IncrementByValue(context.TODO(), test1.Key, test1.Val)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]int64{test1.Key: newValue})
}

func (h *handler) Signature(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var test2 test2
	err = json.Unmarshal(requestBody, &test2)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	hmac := utils.GenerateHMAC(test2.SLine, test2.Key)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{"result": hmac})
}

func (h *handler) Multiplication(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var test []test3
	err = json.Unmarshal(requestBody, &test)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	line := convToString(test)
	res, err := connectingToApp(line)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	toMap, err := convToMap(test, res)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(toMap)
}

func connectingToApp(request string) (string, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer conn.Close()

	n, err := conn.Write([]byte(request))
	if err != nil || n == 0 {
		log.Println(err)
		return "", err
	}

	response := make([]byte, 1024*4)
	n, err = conn.Read(response)
	if err != nil || n == 0 {
		log.Println(err)
		return "", err
	}

	return string(response[0:n]), nil
}

func convToMap(test []test3, line string) (map[string]int, error) {
	res := make(map[string]int)
	fields := strings.Fields(line)
	if len(test) != len(fields) {
		return nil, fmt.Errorf("model length(%d) != length of the result(%d)", len(test), len(fields))
	}
	for i := 0; i < len(test); i++ {
		num, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, err
		}
		res[test[i].Key] = num
	}

	return res, nil
}

func convToString(test []test3) string {
	var res string
	for _, t := range test {
		res += fmt.Sprintf("%s,%s\r\n", t.ALine, t.BLine)
	}
	res += fmt.Sprintf("\r\n")
	return res
}
