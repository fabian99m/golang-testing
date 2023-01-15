package handler

import (
	"errors"
	"net/http"
	"strconv"
)

var ErrNotFound = errors.New(strconv.Itoa(http.StatusNotFound))
var ErrBadRequest = errors.New(strconv.Itoa(http.StatusBadRequest))
