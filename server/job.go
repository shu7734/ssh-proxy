package main

import (
	"fmt"
	"strings"
)

type Job struct {
	Id    string
	Token string
	Host  string
}

func JobFromAuthString(authUser string) (*Job, error) {
	splitted := strings.Split(authUser, "+")

	if len(splitted) != 2 {
		return nil, fmt.Errorf("Wrong username format")
	}

	jobId, token := splitted[0], splitted[1]

	return &Job{Id: jobId, Token: token}, nil
}
