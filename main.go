package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pairProgramming/helper"
)

type ResponsePost struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	res, err := getPost()
	helper.PanicIfNeed(err)

	for _, post := range res {
		helper.Print(post.Title)
	}
}

func getPost() (posts []ResponsePost, err error) {
	var result []byte

	go func() {
		client, err := http.Get("https://jsonplaceholder.typicode.com/posts")
		helper.PanicIfNeed(err)
		defer client.Body.Close()
		result, err = ioutil.ReadAll(client.Body)
		helper.PanicIfNeed(err)
	}()

	var data []ResponsePost

	json.Unmarshal(result, &data)

	return data, nil

}
