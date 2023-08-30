/*
*

	author: Jerry
*/
package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/eep/service"
)
type Offer struct {
    OfferID string           `json:"offerID"`
    UserID string            `json:"userID"`
    ArrTime int              `json:"arrTime"`
    DepTime int              `json:"depTime"`
    ArrSoC int               `json:"arrSoC"`
    DepSoC int               `json:"depSoC"`
    Acdc int                 `json:"acdc"`
}
type Option struct {
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    MaxSoC int                 `json:"maxSoC"`
    Price int                  `json:"price"`
}
type Match struct {
	MatchID string             `json:"matchID"`
	OfferID string             `json:"offerID"`
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
	Date string                `json:"date"`
	ArrTime int                `json:"arrTime"`
	DepTime int                `json:"depTime"`
	ArrSoC int             	   `json:"arrSoC"`
	MaxSoC int                 `json:"maxSoC"`
	Price int                  `json:"price"`
}
type Choice struct {
	StationID string           `json:"stationID"`
	ChargerID string           `json:"chargerID"`
	Acdc string                `json:"acdc"`
	ArrTime string             `json:"arrTime"`
	DepTime string             `json:"depTime"`
	MaxSoC string              `json:"maxSoC"`
	Price string               `json:"price"`
}
type Power struct {
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    Power int                  `json:"power"`
    State int                  `json:"state"`
    TimeStamp int              `json:"timeStamp"`
}
type Power2 struct {
	PowerID string             `json:"powerID"`
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    Power int                  `json:"power"`
    State int                  `json:"state"`
    TimeStamp int              `json:"timeStamp"`
}

type Application struct {
	Fabric *service.ServiceSetup
}

func setCookie(w http.ResponseWriter, name, value string) {
    cookie := http.Cookie{
        Name:    name,
        Value:   value,
        Expires: time.Now().Add(3 * time.Hour),
    }
    http.SetCookie(w, &cookie)
}
func getCookieValue(r *http.Request, name string) (string, error) {
    cookie, err := r.Cookie(name)
    if err != nil {
        return "", err
    }
    return cookie.Value, nil
}
func setDataInCookies(w http.ResponseWriter, name string, data interface{}) {
    dataAsBytes, _ := json.Marshal(data)
    encodedData := url.QueryEscape(string(dataAsBytes))
    setCookie(w, name, encodedData)
}
func getDataFromCookies(r *http.Request, name string, data interface{}) error {
    cookie, err := r.Cookie(name)
    if err != nil {
        return err
    }
    encodedValue := cookie.Value
    decodedValue, err := url.QueryUnescape(encodedValue)
    if err != nil {
        return err
    }
    if err := json.Unmarshal([]byte(decodedValue), &data);
	err != nil {
        return err
    }
    return nil
}
func clearCookies(w http.ResponseWriter, names ...string) {
    for _, name := range names {
        cookie := http.Cookie{
            Name:    name,
            Value:   "",
            Expires: time.Now().Add(-time.Hour), // 設置過期時間為過去的時間
        }
        http.SetCookie(w, &cookie)
    }
}
func cookiesExist(r *http.Request, name string) bool {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return true
		}
	}
	return false
}
// 將時間轉為區間段
func timeToInt(Time string) (string, error) {
    layout := "15:04"
    t, err := time.Parse(layout, Time)
    if err != nil {
        return "", err
    }
    TimeInSeconds := t.Hour()*3600 + t.Minute()*60 + t.Second()
    Time2 := strconv.Itoa((TimeInSeconds / 300) + 1)
    return Time2, nil
}
func intToTime(interval int) (string) {
    totalSeconds := interval * 300
    hours := totalSeconds / 3600
    totalSeconds %= 3600
    minutes := totalSeconds / 60
    seconds := totalSeconds % 60
    
    formattedTime := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
    return formattedTime
}

func (app *Application) CreateAccountView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "createAccount.html", nil)
}
func (app *Application) LoginView(w http.ResponseWriter, r *http.Request) {
	if cookiesExist(r, "now_userID") {
		now_userID, _ := getCookieValue(r, "now_userID")
		fmt.Printf("[使用者登出] %s\n", now_userID)
	}
	clearCookies(w, "now_userID", "now_offerID", "now_matchID", "choice", "optionA", "optionB", "optionC")
	showView(w, r, "login.html", nil)
}
func (app *Application) MainPageView(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		NowCarID string
	}{
		NowCarID:"",
	}
	now_userID, _ := getCookieValue(r, "now_userID")
	data.NowCarID, _ = app.Fabric.ShowCarbyID(now_userID)
	showView(w, r, "mainPage.html", data)
}
func (app *Application) RequestView(w http.ResponseWriter, r *http.Reques裝) {
	if cookiesExist(r, "now_offerID") {
		app.Request4View(w, r)
	} else {
		app.Request1View(w, r)
	}
}
func (app *Application) Request1View(w http.ResponseWriter, r *http.Request) {
	clearCookies(w, "now_offerID", "now_matchID", "choice", "optionA", "optionB", "optionC")
	showView(w, r, "request1.html", nil)
}
type request2_data struct {
	FlagA bool
	MsgA1 string
	MsgA2 string
	MsgA3 string
	FlagB bool
	MsgB1 string
	MsgB2 string
	MsgB3 string
	FlagC bool
	MsgC1 string
	MsgC2 string
	MsgC3 string
	Flag bool
	Msg string
} 
func showOptions(r *http.Request) (data request2_data) {
	if cookiesExist(r, "optionA") {
		var optionA Option
		getDataFromCookies(r, "optionA", &optionA)
		data.FlagA = true
		data.MsgA1 = strconv.Itoa(optionA.MaxSoC)
		data.MsgA2 = strconv.Itoa(optionA.Price)
		data.MsgA3 = strconv.Itoa(optionA.Price)
	}

	if cookiesExist(r, "optionB") {
		var optionB Option
		getDataFromCookies(r, "optionB", &optionB)
		data.FlagB = true
		data.MsgB1 = strconv.Itoa(optionB.MaxSoC)
		data.MsgB2 = strconv.Itoa(optionB.Price)
		data.MsgB3 = strconv.Itoa(optionB.Price)
	}

	if cookiesExist(r, "optionC") {
		var optionC Option
		getDataFromCookies(r, "optionC", &optionC)
		data.FlagC = true
		data.MsgC1 = strconv.Itoa(optionC.MaxSoC)
		data.MsgC2 = strconv.Itoa(optionC.Price)
		data.MsgC3 = strconv.Itoa(optionC.Price)
	}
	return data
}
func (app *Application) Request2View(w http.ResponseWriter, r *http.Request) {
	clearCookies(w, "choice")
	data := showOptions(r)
	showView(w, r, "request2.html", data)
}
func (app *Application) Request3View(w http.ResponseWriter, r *http.Request) {
	var offer Offer
	now_offerID, _ := getCookieValue(r, "now_offerID")
	offerAsBytes, _ := app.Fabric.ShowOfferbyID(now_offerID)
	json.Unmarshal([]byte(offerAsBytes), &offer)

	var option Option
	getDataFromCookies(r, "choice", &option)

	var data Choice
	switch option.StationID {
	case "A":
		data.StationID = "甲地"
	case "B":
		data.StationID = "乙地"
	case "C":
		data.StationID = "丙地"
	}
	data.ChargerID = strconv.Itoa(option.ChargerID)
	if offer.Acdc == 1 {
		data.Acdc = "慢充"
	} else {
		data.Acdc = "快充"
	}
	data.ArrTime = intToTime(offer.ArrTime)
	data.DepTime = intToTime(offer.DepTime)
	data.MaxSoC = strconv.Itoa(option.MaxSoC)
	data.Price = strconv.Itoa(option.Price)

	showView(w, r, "request3.html", data)
}
func (app *Application) Request4View(w http.ResponseWriter, r *http.Request) {
	var offer Offer
	now_offerID, _ := getCookieValue(r, "now_offerID")
	offerAsBytes, _ := app.Fabric.ShowOfferbyID(now_offerID)
	json.Unmarshal([]byte(offerAsBytes), &offer)

	var option Option
	getDataFromCookies(r, "choice", &option)

	var data Choice
	switch option.StationID {
	case "A":
		data.StationID = "甲地"
	case "B":
		data.StationID = "乙地"
	case "C":
		data.StationID = "丙地"
	}
	data.ChargerID = strconv.Itoa(option.ChargerID)
	if offer.Acdc == 1 {
		data.Acdc = "慢充"
	} else {
		data.Acdc = "快充"
	}
	data.ArrTime = intToTime(offer.ArrTime)
	data.DepTime = intToTime(offer.DepTime)
	data.MaxSoC = strconv.Itoa(option.MaxSoC)
	data.Price = strconv.Itoa(option.Price)
	showView(w, r, "request4.html", data)
}
func (app *Application) TrackView(w http.ResponseWriter, r *http.Request) {
	if cookiesExist(r, "now_offerID") {
		app.TrackYesView(w, r)
	} else {
		app.TrackNoView(w, r)
	}
}
func (app *Application) TrackNoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "trackNo.html", nil)
}
func (app *Application) TrackYesView(w http.ResponseWriter, r *http.Request) {
	var offer Offer
	now_offerID, _ := getCookieValue(r, "now_offerID")
	offerAsBytes, _ := app.Fabric.ShowOfferbyID(now_offerID)
	json.Unmarshal([]byte(offerAsBytes), &offer)

<<<<<<< HEAD
var index int = 1
func (app *Application) Schedule() {
	// 設定每5分鐘執行
	ticker := time.NewTicker(5 * time.Minute)
	time.Sleep(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			// 顯示現在時間
			fmt.Println("Current time:", time.Now())
			// 將各個充電樁的功率上鏈
			fmt.Printf("第%d區間上鍊開始\n",index)
			var powers []Power
			for j := 1; j <= 12; j++{
				powers = append(powers, Power{StationID: "A", ChargerID: j, Power: 0, State: 0, TimeStamp: index})
			}
			for j := 1; j <= 6; j++{
				powers = append(powers, Power{StationID: "B", ChargerID: j, Power: 30, State: 1, TimeStamp: index})
			}
			for j := 1; j <= 20; j++{
				powers = append(powers, Power{StationID: "C", ChargerID: j, Power: 40, State: 1, TimeStamp: index})
			}
			app.Power(powers)
			fmt.Printf("第%d區間上鍊結束\n",index)
			index++
=======
	var option Option
	getDataFromCookies(r, "choice", &option)

	now_matchID, _ := getCookieValue(r, "now_matchID")

	type AllPower struct {
        Powers []Power2 `json:"powers"`
    }
    var allPower AllPower
	
	// msg, err := app.Fabric.ShowAllPower()
	// if err != nil {
	// 	log.Fatalln(err)
	// }else{
	// 	fmt.Println(msg)
	// }
	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04")
	timeAsInt, err := timeToInt(formattedTime)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	checktimeasint, _ := strconv.Atoi(timeAsInt)
	fmt.Println("Current time:", formattedTime)
	fmt.Printf("第%d區間\n",checktimeasint)
	fmt.Printf("ArrTime: %d\n", offer.ArrTime)
	fmt.Printf("DepTime: %d\n", offer.DepTime)

    AllPowerAsBytes, err := app.Fabric.ShowPowerbyMatch(now_matchID)
	if err != nil {
		fmt.Println(err)
	}
    json.Unmarshal([]byte(AllPowerAsBytes), &allPower)
    powers := allPower.Powers
	
	fmt.Println(AllPowerAsBytes)

	var timeChargePairs []struct {
		Time   string
		Charge int
	}
    for _, power := range powers {
        data := struct {
			Time   string
			Charge int
        }{
            Time: intToTime(power.TimeStamp),
            Charge: power.Power,
        }
        timeChargePairs = append(timeChargePairs, data)
    }

	data := struct {
		Msg1            string
		Msg2            string
		Msg3            string
		TimeChargeArray []struct {
			Time   string
			Charge int
>>>>>>> 2a1fad4c8c01e2b3ab234e7c0bd510265ea88d5b
		}
	}{
		Msg1: "",
		Msg2: "",
		Msg3: "",
		TimeChargeArray: timeChargePairs,
	}

	switch option.StationID {
	case "A":
		data.Msg1 = "甲地"
	case "B":
		data.Msg1 = "乙地"
	case "C":
		data.Msg1 = "丙地"
	}
	data.Msg2 = strconv.Itoa(option.ChargerID)
	if offer.Acdc == 1 {
		data.Msg3 = "慢充"
	} else {
		data.Msg3 = "快充"
	}
	showView(w, r, "trackYes.html", data)
}
func (app *Application) HistoryView(w http.ResponseWriter, r *http.Request) {
	now_userID, _ := getCookieValue(r, "now_userID")
	_, err := app.Fabric.ShowMatchbyUser(now_userID)
	if err == nil {
		app.HistoryListView(w, r)
	} else {
		app.HistoryNoView(w, r)
	}
}
func (app *Application) HistoryListView(w http.ResponseWriter, r *http.Request) {
    type AllMatch struct {
        Matchs []Match `json:"matchs"`
    }

    var allMatch AllMatch
	now_userID, _ := getCookieValue(r, "now_userID")
    AllMatchAsBytes, _ := app.Fabric.ShowMatchbyUser(now_userID)
    json.Unmarshal([]byte(AllMatchAsBytes), &allMatch)
    matchs := allMatch.Matchs

    var allMatchData []struct {
		Msg1 string
		Msg2 string
		Msg3 string
		Msg4 string
		Msg5 string
    }

    for _, match := range matchs {
        data := struct {
            Msg1 string
            Msg2 string
            Msg3 string
            Msg4 string
            Msg5 string
        }{
            Msg1: match.Date,
            Msg2: intToTime(match.ArrTime),
            Msg3: intToTime(match.DepTime),
            Msg4: strconv.Itoa(match.MaxSoC),
            Msg5: strconv.Itoa(match.Price),
        }
        allMatchData = append(allMatchData, data)
    }

    showView(w, r, "historyList.html", allMatchData)
}
func (app *Application) HistoryNoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "historyNo.html", nil)
}

func (app *Application) Schedule() {
	for {
		currentTime := time.Now()
		// 檢查現在時間是否是5的倍數
		if currentTime.Minute() % 5 == 0 {
			// 顯示現在時間
			formattedTime := currentTime.Format("15:04")
	
			timeAsInt, err := timeToInt(formattedTime)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			checktimeasint, _ := strconv.Atoi(timeAsInt)
	
			fmt.Println("Current time:", formattedTime)
			// 將各個充電裝的功率上鏈
			fmt.Printf("第%d區間上鍊開始\n",checktimeasint)
			var powers []Power
			for j := 1; j <= 12; j++{
				powers = append(powers, Power{StationID: "A", ChargerID: j, Power: 0, State: 0, TimeStamp: checktimeasint})
			}
			for j := 1; j <= 6; j++{
				powers = append(powers, Power{StationID: "B", ChargerID: j, Power: 30, State: 1, TimeStamp: checktimeasint})
			}
			for j := 1; j <= 20; j++{
				powers = append(powers, Power{StationID: "C", ChargerID: j, Power: 40, State: 1, TimeStamp: checktimeasint})
			}
			app.Power(powers)
			fmt.Printf("第%d區間上鍊結束\n",checktimeasint)
		}
		// 等待一分鐘後再次檢查
		time.Sleep(time.Minute)
	}
}


func (app *Application) Register(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	userName := r.FormValue("userName")
	capacity := r.FormValue("capacity")
	password := r.FormValue("password")

	// 查看車牌是否已註冊過
	_ , err := app.Fabric.GetUserIDbyCarID(carID)
	if err == nil {
		err = errors.New("Description: 車牌曾經註冊!請登入!")
	}else{
		_, err = app.Fabric.Register(carID, userName, capacity, password)
	}

	if err != nil {
		// 註冊失敗，要求重新輸入
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		errMessage := strings.Split(err.Error(), "Description: ")
		data.Msg = errMessage[1]
		showView(w, r, "createAccount.html", data)
	}else{
		// 註冊成功則登入
		now_userID, _ := app.Fabric.Login(carID, password)
		setCookie(w, "now_userID", now_userID)
		fmt.Printf("[新使用者註冊並登入] %s\n", now_userID)
		http.Redirect(w, r, "/mainPage.html", http.StatusSeeOther)
	}
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	password := r.FormValue("password")

	now_userID, err := app.Fabric.Login(carID, password)

	if err != nil {
		// 登入失敗，要求重新輸入
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		errMessage := strings.Split(err.Error(), "Description: ")
		data.Msg = errMessage[1]
		showView(w, r, "login.html", data)
	}else{
		// 登入成功
		setCookie(w, "now_userID", now_userID)
		fmt.Printf("[使用者登入] %s\n", now_userID)
		http.Redirect(w, r, "/mainPage.html", http.StatusSeeOther)
	}
}

func (app *Application) Offer(w http.ResponseWriter, r *http.Request)  {
	// 將進場時間轉換為區間段
	arrTime := r.FormValue("arrTime")
	arrTime2, err := timeToInt(arrTime)
	if err != nil {
		http.Error(w, "Invalid arrtime format", http.StatusBadRequest)
	}

	// 將離場時間轉換為區間段
	depTime := r.FormValue("depTime")
	depTime2, err := timeToInt(depTime)
	if err != nil {
		http.Error(w, "Invalid depTime format", http.StatusBadRequest)
	}

	arrSoC := r.FormValue("arrSoC")
	depSoC := r.FormValue("depSoC")
	acdc := r.FormValue("acdc")
	now_userID, _ := getCookieValue(r, "now_userID")
	now_offerID, err := app.Fabric.Offer(arrTime2, depTime2, arrSoC, depSoC, acdc, now_userID)

	if err != nil {
		// 申請充電失敗
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		errMessage := strings.Split(err.Error(), "Description: ")
		data.Msg = errMessage[1]
		showView(w, r, "request1.html", data)
	}else{
		// 申請充電成功
		setCookie(w, "now_offerID", now_offerID)
		fmt.Printf("[發送新請求] %s: %s\n", now_userID, now_offerID)
		
		// 呼叫最佳化函式1
		optionA := Option{StationID: "A", ChargerID: 1, MaxSoC: 100, Price: 100}
		optionB := Option{StationID: "B", ChargerID: 2, MaxSoC: 100, Price: 50}
		optionC := Option{StationID: "C", ChargerID: 3, MaxSoC: 100, Price: 30}
		setDataInCookies(w, "optionA", optionA)
		setDataInCookies(w, "optionB", optionB)
		setDataInCookies(w, "optionC", optionC)
		http.Redirect(w, r, "/request2.html", http.StatusSeeOther)
	}
}
func (app *Application) Choice(w http.ResponseWriter, r *http.Request) {
	choice := r.FormValue("choice")
	var option Option
	switch choice {
	case "place_1":
		getDataFromCookies(r, "optionA", &option)
	case "place_2":
		getDataFromCookies(r, "optionB", &option)
	case "place_3":
		getDataFromCookies(r, "optionC", &option)
	default:
		// 使用者沒有做出選擇
		data := showOptions(r)
		data.Flag = true
		data.Msg = "請做出選擇後重新送出"
		showView(w, r, "request2.html", data)
		return
	}
	// 使用者做出選擇
	setDataInCookies(w, "choice", option)
	http.Redirect(w, r, "/request3.html", http.StatusSeeOther)
}

func (app *Application) Match(w http.ResponseWriter, r *http.Request)  {
	clearCookies(w, "optionA", "optionB", "optionC")

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	
	var option Option
	getDataFromCookies(r, "choice", &option)

	var offer Offer
	now_offerID, _ := getCookieValue(r, "now_offerID")
	offerAsBytes, _ := app.Fabric.ShowOfferbyID(now_offerID)
	json.Unmarshal([]byte(offerAsBytes), &offer)

	now_matchID, _ := app.Fabric.Match(option.StationID, strconv.Itoa(option.ChargerID), date, strconv.Itoa(offer.ArrTime), strconv.Itoa(offer.DepTime), strconv.Itoa(offer.ArrSoC), strconv.Itoa(option.MaxSoC), strconv.Itoa(option.Price), now_offerID)
	fmt.Printf("[發送新配對] %s: %s\n", now_offerID, now_matchID)
	setCookie(w, "now_matchID", now_matchID)

	app.Request4View(w, r)
}

// func (app *Application) ShowAllMatch(w http.ResponseWriter, r *http.Request)  {
// 	msg, err := app.Fabric.ShowAllMatch()

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:true,
// 		Msg1:"",
// 		Flag2:false,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg1 = errMessage[1]
// 	}else{
// 		data.Msg1 = "ShowAllMatch success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }

func (app *Application) Power(powers []Power)  {
    powersAsBytes, err := json.Marshal(powers)
    if err != nil {
        log.Fatalln(err)
    }
	_, err = app.Fabric.Power(powersAsBytes)
	if err != nil {
		log.Fatalln(err)
	}
}
// func (app *Application) ShowAllPower(w http.ResponseWriter, r *http.Request)  {
// 	msg, err := app.Fabric.ShowAllPower()

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:false,
// 		Msg1:"",
// 		Flag2:true,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg2 = errMessage[1]
// 	}else{
// 		data.Msg2 = "ShowAllPower success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }
// func (app *Application) ShowPowerbyCharger(w http.ResponseWriter, r *http.Request)  {
// 	stationID := r.FormValue("stationID_search")
// 	chargerID := r.FormValue("chargerID_search")

// 	msg, err := app.Fabric.ShowPowerbyCharger(stationID, chargerID)

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:false,
// 		Msg1:"",
// 		Flag2:true,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg2 = errMessage[1]
// 	}else{
// 		data.Msg2 = "ShowPowerbyCharger success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }