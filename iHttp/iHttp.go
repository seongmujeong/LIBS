package iHttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

/*
	note
		HTTP 헤더에 새로운 Key-Value 추가
	param
		req				: HTTP Request 구조체 포인터
		Key				: 헤더에 추가할 데이터 Key값
		Value			: 헤더에 추가할 데이터 Value값
*/
func SetHeader(req *http.Request, Key string, Value string) {
	req.Header.Add(Key, Value)
}

/*
	note
		HTTP Basic 인증 헤더에 추가
	param
		req				: HTTP Request 구조체 포인터
		UserName		: 사용자명
		PassWord		: 비밀번호
*/
func SetBasicAuth(req *http.Request, UserName string, PassWord string) {
	req.SetBasicAuth(UserName, PassWord)
}

/*
	note
		신규 Request 생성
	param
		Method			: HTTP Method 정보
		Url				: 연결할 Url 정보
		Body			: POST시 전송 할 Body 데이터
	return
		*http.Client	: HTTP Client 구조체 포인터
		*http.Request	: HTTP Request 구조체 포인터
		error			: error 발생시 err 정보
*/
func CreateRequest(Method string, Url string, Body []byte, timeOut int32) (*http.Client, *http.Request, error) {

	buff := bytes.NewBuffer(Body)
	client := &http.Client{
		Timeout: time.Duration(timeOut) * time.Second,
	}
	req, err := http.NewRequest(Method, Url, buff)
	if err != nil {
		fmt.Println("Create Http Request Fail")
	}
	return client, req, err
}

/*
	note
		Request 정보 Send
	param
		client			: HTTP Client 구조체 포인터
		req				: HTTP Request 구조체 포인터
	return
		[]byte			: HTTP Request 구조체 포인터
		error			: error 발생시 err 정보
*/
func SendRequest(client *http.Client, req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Send Http Request Fail")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	return respBody, err
}

/*
	note
		URL에 Parameter 추가
	param
		req				: HTTP Request 구조체 포인터
		Key				: Parameter Key 값
		Value			: Parameter Value 값
*/
func AddQueryParams(req *http.Request, Key string, Value string) {
	query := req.URL.Query()
	query.Add(Key, Value)
	req.URL.RawQuery = query.Encode()
}

/*
	note
		HTTP 서버에 URI 및 콜백함수 등록
	param
		Uri				: 등록할 URI 정보
		reg				: 등록할 콜백 함수 포인터
*/
func RegistCallback(Uri string, reg func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(Uri, reg)
}

/*
	note
		HTTP 서버 스타트
	param
		Port			: HTTP 서버 포트 정보
*/
func StartServer(Port int32) {
	portStr := fmt.Sprintf(":%d", Port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

/*
	note
		WEBSOCK CONN 생성
	param
		Uri				: WebSock 연결 주소
		TImeOut			: WebSock TimeOut
	return
		*websocket.Conn : WebSock Connection 포인터
*/
func CreateWebSock(Uri string, TimeOut int32) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: time.Duration(TimeOut) * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
	conn, _, err := dialer.Dial(Uri, nil)
	if err != nil {
		fmt.Println("Error connecting to Websocket Server")
	}
	return conn, err
}

/*
	note
		WEBSOCK Message Send
	param
		Conn			: WebSock Connection 포인터
		Body			: 전송할 메시지 Body 데이터
*/
func SendWsMsg(Conn *websocket.Conn, Body string) error {
	err := Conn.WriteMessage(websocket.TextMessage, []byte(Body))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
	}
	return err
}

/*
	note
		WEBSOCK Message Recv
	param
		Conn			: WebSock Connection 포인터
	return
		string			: 응답받은 메시지 Body 데이터
*/
func RecvWsMsg(Conn *websocket.Conn) string {
	_, msg, err := Conn.ReadMessage()
	if err != nil {
		log.Println("Error in receive:", err)
	}
	return string(msg[:])
}
