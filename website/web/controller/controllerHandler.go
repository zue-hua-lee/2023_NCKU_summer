/**
  author: Jerry
 */
package controller

import (
	"errors"
	"net/http"
	"github.com/hyperledger/fabric/eep/service"
)

type Application struct {
	Fabric *service.ServiceSetup
}
func (app *Application) MyHomeView(w http.ResponseWriter, r *http.Request)  {
	showView(w, r, "myHome.html", nil)
}
func (app *Application) MyMatchView(w http.ResponseWriter, r *http.Request)  {
	showView(w, r, "myMatch.html", nil)
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	userName := r.FormValue("userName")
	capacity := r.FormValue("capacity")
	password := r.FormValue("password")

	var transactionID string
	_ , err := app.Fabric.GetUserIDbyCarID(carID)	
	if err == nil {
		err = errors.New("Register ERROR! CarID already exist!")
	}else{
		transactionID, err = app.Fabric.Register(carID, userName, capacity, password)
	}

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
		Flag3 bool
		Msg3 string
	}{
		Flag1:true,
		Msg1:"",
		Flag2:false,
		Msg2:"",
		Flag3:false,
		Msg3:"",
	}
	if err != nil {
		data.Msg1 = err.Error()
	}else{
		data.Msg1 = "Register success, Transaction ID: " + transactionID
	}
	showView(w, r, "myHome.html", data)
}
func (app *Application) ShowAllUser(w http.ResponseWriter, r *http.Request)  {
	msg, err := app.Fabric.ShowAllUser()

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
		Flag3 bool
		Msg3 string
	}{
		Flag1:true,
		Msg1:"",
		Flag2:false,
		Msg2:"",
		Flag3:false,
		Msg3:"",
	}
	if err != nil {
		data.Msg1 = err.Error()
	}else{
		data.Msg1 = "ShowAllUser success: " + msg
	}
	showView(w, r, "myHome.html", data)
}

var now_userID string = ""
func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	password := r.FormValue("password")

	var err error
	now_userID, err = app.Fabric.Login(carID, password)

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
		Flag3 bool
		Msg3 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:true,
		Msg2:"",
		Flag3:false,
		Msg3:"",
	}
	if err != nil {
		data.Msg2 = err.Error()
	}else{
		data.Msg2 = "Login success! " + now_userID + " login!"
		// now_userID, _ = app.Fabric.GetUserIDbyCarID(carID)
	}
	showView(w, r, "myHome.html", data)
}
func (app *Application) ShowNowUser(w http.ResponseWriter, r *http.Request)  {
	msg, err := app.Fabric.ShowUserbyID(now_userID)

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
		Flag3 bool
		Msg3 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:true,
		Msg2:"",
		Flag3:false,
		Msg3:"",
	}
	if err != nil {
		data.Msg2 = err.Error()
	}else{
		data.Msg2 = "ShowNowUser success: " + msg
	}
	showView(w, r, "myHome.html", data)
}

var now_offerID string = ""
func (app *Application) Offer(w http.ResponseWriter, r *http.Request)  {
	arrTime := r.FormValue("arrTime")
	depTime := r.FormValue("depTime")
	arrSoC := r.FormValue("arrSoC")
	depSoC := r.FormValue("depSoC")
	acdc := r.FormValue("acdc")
	origin := r.FormValue("origin")

	var err error
	now_offerID, err = app.Fabric.Offer(arrTime, depTime, arrSoC, depSoC, acdc, origin, now_userID)

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
		Flag3 bool
		Msg3 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:false,
		Msg2:"",
		Flag3:true,
		Msg3:"",
	}
	if err != nil {
		data.Msg3 = err.Error()
	}else{
		data.Msg3 = "Offer success, Transaction ID: " + now_offerID
	}
	showView(w, r, "myHome.html", data)
}
func (app *Application) ShowAllOffer(w http.ResponseWriter, r *http.Request)  {
	msg, err := app.Fabric.ShowAllOffer()

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
		Flag3 bool
		Msg3 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:false,
		Msg2:"",
		Flag3:true,
		Msg3:"",
	}
	if err != nil {
		data.Msg3 = err.Error()
	}else{
		data.Msg3 = "ShowAllOffer success: " + msg
	}
	showView(w, r, "myHome.html", data)
}

func (app *Application) Match(w http.ResponseWriter, r *http.Request)  {
	stationID := r.FormValue("stationID")
	maxSoC := r.FormValue("maxSoC")
	price := r.FormValue("price")

	transactionID, err := app.Fabric.Match(stationID, maxSoC, price, now_offerID)

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
	}{
		Flag1:true,
		Msg1:"",
		Flag2:false,
		Msg2:"",
	}
	if err != nil {
		data.Msg1 = err.Error()
	}else{
		data.Msg1 = "Match success, Transaction ID: " + transactionID
	}
	showView(w, r, "myMatch.html", data)
}
func (app *Application) ShowAllMatch(w http.ResponseWriter, r *http.Request)  {
	msg, err := app.Fabric.ShowAllMatch()

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
	}{
		Flag1:true,
		Msg1:"",
		Flag2:false,
		Msg2:"",
	}
	if err != nil {
		data.Msg1 = err.Error()
	}else{
		data.Msg1 = "ShowAllMatch success: " + msg
	}
	showView(w, r, "myMatch.html", data)
}
func (app *Application) Power(w http.ResponseWriter, r *http.Request)  {
	stationID := r.FormValue("stationID")
	chargerID := r.FormValue("chargerID")
	power := r.FormValue("power")
	state := r.FormValue("state")
	timestamp := r.FormValue("timestamp")

	transactionID, err := app.Fabric.Power(stationID, chargerID, power, state, timestamp)

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:true,
		Msg2:"",
	}
	if err != nil {
		data.Msg2 = err.Error()
	}else{
		data.Msg2 = "Power success, Transaction ID: " + transactionID
	}
	showView(w, r, "myMatch.html", data)
}
func (app *Application) ShowAllPower(w http.ResponseWriter, r *http.Request)  {
	msg, err := app.Fabric.ShowAllPower()

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:true,
		Msg2:"",
	}
	if err != nil {
		data.Msg2 = err.Error()
	}else{
		data.Msg2 = "ShowAllPower success: " + msg
	}
	showView(w, r, "myMatch.html", data)
}
func (app *Application) ShowPowerbyCharger(w http.ResponseWriter, r *http.Request)  {
	stationID := r.FormValue("stationID_search")
	chargerID := r.FormValue("chargerID_search")

	msg, err := app.Fabric.ShowPowerbyCharger(stationID, chargerID)

	data := &struct {
		Flag1 bool
		Msg1 string
		Flag2 bool
		Msg2 string
	}{
		Flag1:false,
		Msg1:"",
		Flag2:true,
		Msg2:"",
	}
	if err != nil {
		data.Msg2 = err.Error()
	}else{
		data.Msg2 = "ShowPowerbyCharger success: " + msg
	}
	showView(w, r, "myMatch.html", data)
}