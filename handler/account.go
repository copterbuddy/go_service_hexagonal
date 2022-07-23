package handler

import (
	"bank/errs"
	"bank/logs"
	"bank/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type accountHandler struct {
	accSrv service.AccountService
}

func NewAccountHandler(accSrv service.AccountService) accountHandler {
	return accountHandler{accSrv: accSrv}
}

func (h accountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	customerID, err := strconv.Atoi(mux.Vars(r)["customerID"])
	if err != nil {
		logs.Error(err)
		handleError(w, errs.NewUnExpectedError())
		return
	}

	logs.Error("test convert: " + strconv.Itoa(customerID))
	if r.Header.Get("Content-Type") != "application/json" {
		handleError(w, errs.NewValidationError("request body incorrect format1"))
		return
	}

	request := service.NewAccountRequest{}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logs.Error(err)
		handleError(w, errs.NewValidationError("request body incorrect format2"))
		return
	}

	response, err := h.accSrv.NewAccount(customerID, request)
	if err != nil {
		logs.Error(err)
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h accountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	customerID, _ := strconv.Atoi(mux.Vars(r)["customerID"])

	responses, err := h.accSrv.GetAccount(customerID)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(responses)
}
