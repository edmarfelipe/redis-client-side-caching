package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/rueian/rueidis"
)

type Value struct {
	Value string `json:"value"`
	Took  string `json:"took"`
}

func resp(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := flag.String("port", "3000", "http port server")
	flag.Parse()

	if port == nil {
		log.Println("need to specify port")
		return
	}

	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		command := c.B().
			Get().Key("key").Build()

		res := c.Do(r.Context(), command)
		if res.Error() != nil {
			resp(w, 500, err)
			return
		}

		value, err := res.ToString()
		if res.Error() != nil {
			resp(w, 500, err)
			return
		}

		resp(w, 200, Value{
			Value: value,
			Took:  time.Since(start).String(),
		})
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		value := time.Now().Format(time.DateTime)
		command := c.B().
			Set().
			Key("key").
			Value(value).
			Build()

		err = c.Do(r.Context(), command).Error()
		if err != nil {
			resp(w, 500, err)
			return
		}

		resp(w, 200, Value{
			Value: value,
			Took:  time.Since(start).String(),
		})
	})

	http.HandleFunc("/get-cached", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		command := c.B().
			Get().Key("key").Cache()

		res := c.DoCache(r.Context(), command, time.Second*60)
		if res.Error() != nil {
			resp(w, 500, err)
			return
		}

		value, err := res.ToString()
		if res.Error() != nil {
			resp(w, 500, err)
			return
		}

		resp(w, 200, Value{
			Value: value,
			Took:  time.Since(start).String(),
		})
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
