package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/go-zoo/bone"
	"github.com/satori/go.uuid"
)

const (
	StatusFail       = "failure"
	StatusInProgress = "in-progress"
	StatusSuccess    = "success"
)

func CreateRoute(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendErr(res, http.StatusBadRequest, err)
		return
	}

	path := Path{}
	token := uuid.NewV4().String()

	err = json.Unmarshal(data, &path)
	if err != nil {
		sendErr(res, http.StatusBadRequest, err)
		return
	}

	record := Record{
		Status: StatusInProgress,
		Path:   path,
	}

	// store the path in db
	err = db.Set(token, record, 0).Err()
	if err != nil {
		sendErr(res, http.StatusInternalServerError, err)
		return
	}

	// async call Google Map service and update the record
	go CalcRoute(token, record)

	send(res, http.StatusOK, map[string]string{
		"token": token,
	})
}

func GetRoute(res http.ResponseWriter, req *http.Request) {
	id := bone.GetValue(req, "token")

	val, err := db.Get(id).Result()

	if err == redis.Nil {
		sendErr(res, http.StatusNotFound, errors.New("No such record"))
	} else if err != nil {
		sendErr(res, http.StatusInternalServerError, err)
	} else {
		record := Record{}

		// unmarshal the string data from redis and make it become a struct (object)
		err = json.Unmarshal([]byte(val), &record)
		if err != nil {
			sendErr(res, http.StatusInternalServerError, err)
			return
		}

		if record.Status == StatusInProgress {
			send(res, http.StatusOK, map[string]string{
				"status": StatusInProgress,
			})
		} else if record.Status == StatusFail {
			sendErr(res, http.StatusInternalServerError, err)
		} else {
			send(res, http.StatusOK, record)
		}
	}

}

func send(res http.ResponseWriter, code int, data interface{}) {
	rd.JSON(res, code, data)
}

func sendErr(res http.ResponseWriter, code int, err error) {
	send(res, code, map[string]string{
		"status": StatusFail,
		"error":  err.Error(),
	})
}
