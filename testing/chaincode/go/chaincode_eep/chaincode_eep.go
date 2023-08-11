package main

import (
    "fmt"
    "encoding/json"
    "strconv"
    "strings"
    "unicode"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type MyUser struct {
    UserID string             `json:"userID"`
    CarID string              `json:"carID"`
    UserName string           `json:"userName"`
    Password string           `json:"password"`
    Capacity int              `json:"capacity"`
}
type MyOffer struct {
    OfferID string           `json:"offerID"`
    UserID string            `json:"userID"`
    ArrTime int              `json:"arrTime"`
    DepTime int              `json:"depTime"`
    ArrSoC int               `json:"arrSoC"`
    DepSoC int               `json:"depSoC"`
    Origin string            `json:"origin"`
    Acdc int                 `json:"acdc"`
}
type MyMatch struct {
    MatchID string             `json:"matchID"`
    OfferID string             `json:"offerID"`
	StationID string           `json:"stationID"`
    MaxSoC int                 `json:"maxSoC"`
    Price int                  `json:"price"`
}
type Power struct {
	PowerID string             `json:"powerID"`
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    Power int                  `json:"power"`
    State int                  `json:"state"`
    TimeStamp int              `json:"timeStamp"`
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("eep Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "init" {
        return t.Init(stub)
	} else if function == "register" {
		return t.register(stub, args)
	} else if function == "showAllUser" {
		return t.showAllUser(stub)
	} else if function == "login" {
		return t.login(stub, args)
    } else if function == "getUserIDbyCarID" {
        return t.getUserIDbyCarID(stub, args)
	} else if function == "showUserbyID" {
		return t.showUserbyID(stub, args)
	} else if function == "offer" {
		return t.offer(stub, args)
	} else if function == "showAllOffer" {
		return t.showAllOffer(stub)
	} else if function == "match" {
		return t.match(stub, args)
	} else if function == "showAllMatch" {
		return t.showAllMatch(stub)
	} else if function == "power" {
		return t.power(stub, args)
	} else if function == "showAllPower" {
		return t.showAllPower(stub)
	} else if function == "showPowerbyCharger" {
		return t.showPowerbyCharger(stub, args)
	}
	return shim.Error("Invalid invoke function name.")
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
        fmt.Println("Energy Exchange Platform is coming up")
        fmt.Println("Designed by JerryisaFish in TW")
        funcName, args := stub.GetFunctionAndParameters()
        var number int
        var err error
        txId := stub.GetTxID()
        fmt.Println("eep Init() is running")
        fmt.Println("Transaction ID:", txId)
        fmt.Println("  GetFunctionAndParameters() function:", funcName)
        fmt.Println("  GetFunctionAndParameters() args count:", len(args))
        fmt.Println("  GetFunctionAndParameters() args found:", args)

        if len(args) == 1 {
            fmt.Println("  GetFunctionAndParameters() arg[0] length:", len(args[0]))
            if len(args[0]) == 0 {
                fmt.Println("  Uh oh, args[0] is empty...")
            } else {
                fmt.Println("  Great news everyone, args[0] is not empty")
                number, err = strconv.Atoi(args[0])
                if err != nil {
                    return shim.Error("Expecting a numeric string argument to Init() for instantiate")
                }
                err = stub.PutState("selftest", []byte(strconv.Itoa(number)))
                if err != nil {
                    return shim.Error(err.Error())
                }
            }
        }
        // showing the alternative argument shim function
        alt := stub.GetStringArgs()
        fmt.Println("  GetStringArgs() args count:", len(alt))
        fmt.Println("  GetStringArgs() args found:", alt)
        // store compatible marbles application version
        err = stub.PutState("eep_ver", []byte("1.0.0"))
        if err != nil {
                return shim.Error(err.Error())
        }
        fmt.Println("Ready for Action! let's Fa~ Da~ Tsai~")
        return shim.Success(nil)
}

var user_count int = 0
func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

    var user MyUser
    checkcarID := strings.Split(args[0],"-")
    if len(checkcarID) != 2{
        return shim.Error("車牌號碼格式錯誤!")
    }
    checkcarID1 := strings.Split(checkcarID[0],"")
    if len(checkcarID1) != 3 || checkcarID[0][0] != 'E' || !unicode.IsUpper(rune(checkcarID[0][1])) || !unicode.IsUpper(rune(checkcarID[0][2])) {
        return shim.Error("車牌號碼格式錯誤!")
    }
    checkcarID2, err := strconv.Atoi(checkcarID[1])
    if err != nil || checkcarID2 < 1000 || checkcarID2 > 9999{
        return shim.Error("車牌號碼格式錯誤!")
    }
    user.CarID = args[0]
    if args[1] == "" {
        return shim.Error("尚未填寫車主名稱!")
    }
    user.UserName = args[1]
    checkcapacity, err := strconv.Atoi(args[2])
    if err != nil || checkcapacity <= 0{
        return shim.Error("電池容量格式錯誤!")
    }
    user.Capacity = checkcapacity
    if args[3] == "" {
        return shim.Error("尚未填寫密碼!")
    }
    user.Password = args[3]
    
    userID := "user" + strconv.Itoa(user_count)
    user.UserID = userID
    user_count++

    UserAsBytes, _ := json.Marshal(user)
    err = stub.PutState(userID, UserAsBytes)
	if err != nil {
		return shim.Error(err.Error())
    }
    return shim.Success(nil)
}
func (t *SimpleChaincode) showAllUser(stub shim.ChaincodeStubInterface) pb.Response {
    type myAlluser struct {
        Users []MyUser `json:"users"`
    }
    var alluser myAlluser
    var err error

    userIterator, err := stub.GetStateByRange("user0", "user999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer userIterator.Close()

    for userIterator.HasNext() {
        data, err := userIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var user MyUser
        json.Unmarshal(ValAsBytes, &user)
		alluser.Users = append(alluser.Users, user)
    }

	alluserAsBytes, _ := json.Marshal(alluser)
	return shim.Success(alluserAsBytes)
}

func (t *SimpleChaincode) login(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
    carID := args[0]
    password := args[1]

    userIterator, err := stub.GetStateByRange("user0", "user999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer userIterator.Close()

    for userIterator.HasNext() {
        data, err := userIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var user MyUser
        json.Unmarshal(ValAsBytes, &user)
        if carID == user.CarID {
            if password == user.Password {
                return shim.Success([]byte(user.UserID))
            } else {
                return shim.Error("Password ERROR!")
            }
        }
    }
	return shim.Error("CarID ERROR!")
}
func (t *SimpleChaincode) getUserIDbyCarID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    carID := args[0]

    userIterator, err := stub.GetStateByRange("user0", "user999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer userIterator.Close()

    for userIterator.HasNext() {
        data, err := userIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var user MyUser
        json.Unmarshal(ValAsBytes, &user)
        if carID == user.CarID {
            return shim.Success([]byte(user.UserID))
        }
    }
	return shim.Error("getUserIDbyCarID ERROR!")
}
func (t *SimpleChaincode) showUserbyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
    var user MyUser
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    if args[0] == "" {
		return shim.Error("尚無使用者登入!")
    }
    userID := args[0]
    
    userAsBytes, err := stub.GetState(userID)
    if err != nil {
            return shim.Error("Failed to find user - " + userID)
    }
    json.Unmarshal(userAsBytes, &user)
    if user.UserID != userID {
            return shim.Error("User does not exist - " + userID)
    }
    return shim.Success(userAsBytes)
}

var offer_count int = 0
func (t *SimpleChaincode) offer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

    var offer MyOffer
    checkarrtime, err := strconv.Atoi(args[0])
    if err != nil || checkarrtime < 1 || checkarrtime > 288{
        return shim.Error("開始充電時間格式錯誤!")
    }
    offer.ArrTime = checkarrtime
    checkdeptime, err := strconv.Atoi(args[1])
    if err != nil || checkdeptime < 1 || checkdeptime > 288 || checkdeptime <= offer.ArrTime {
        return shim.Error("結束充電時間格式錯誤!")
    }
    offer.DepTime = checkdeptime
    checkarrSoC, err := strconv.Atoi(args[2])
    if err != nil || checkarrSoC < 0 || checkarrSoC > 100{
        return shim.Error("開始電池狀態格式錯誤!")
    }
    offer.ArrSoC = checkarrSoC
    checkdepSoC, err := strconv.Atoi(args[3])
    if err != nil || checkdepSoC < 0 || checkdepSoC > 100 || checkdepSoC <= offer.ArrSoC {
        return shim.Error("結束電池狀態格式錯誤!")
    }
    offer.DepSoC = checkdepSoC
    checkacdc, err := strconv.Atoi(args[4])
    if err != nil || !(checkacdc == 0 || checkacdc == 1) {
        return shim.Error("快/慢充格式錯誤!")
    }
    offer.Acdc = checkacdc
    if !(args[5] == "A" || args[5] == "B" || args[5] == "C") {
        return shim.Error("起點位置格式錯誤!")
    }
    offer.Origin = args[5]
    if args[6] == "" {
		return shim.Error("尚無使用者登入!")
    }
    offer.UserID = args[6]

    offerID := "offer" + strconv.Itoa(offer_count)
    offer.OfferID = offerID
    offer_count++

    OfferAsBytes, _ := json.Marshal(offer)
    err = stub.PutState(offerID, OfferAsBytes)
	if err != nil {
		return shim.Error(err.Error())
    }
    return shim.Success([]byte(offerID))
}
func (t *SimpleChaincode) showAllOffer(stub shim.ChaincodeStubInterface) pb.Response {
    type myAlloffer struct {
        Offers []MyOffer `json:"offers"`
    }
    var alloffer myAlloffer
    var err error

    offerIterator, err := stub.GetStateByRange("offer0", "offer999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer offerIterator.Close()

    for offerIterator.HasNext() {
        data, err := offerIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var offer MyOffer
        json.Unmarshal(ValAsBytes, &offer)
		alloffer.Offers = append(alloffer.Offers, offer)
    }

	allofferAsBytes, _ := json.Marshal(alloffer)
	return shim.Success(allofferAsBytes)
}

var match_count int = 0
func (t *SimpleChaincode) match(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

    var match MyMatch
    if !(args[0] == "A" || args[0] == "B" || args[0] == "C") {
        return shim.Error("充電站格式錯誤!")
    }
    match.StationID = args[0]
    checkmaxSoC, err := strconv.Atoi(args[1])
    if err != nil || checkmaxSoC < 0 || checkmaxSoC > 100 {
        return shim.Error("最高可充電池狀態格式錯誤!")
    }
    match.MaxSoC = checkmaxSoC
    checkprice, err := strconv.Atoi(args[2])
    if err != nil || checkprice < 0 {
        return shim.Error("電價格式錯誤!")
    }
    match.Price = checkprice
    if args[3] == "" {
		return shim.Error("尚未提出充電申請!")
    }
    match.OfferID = args[3]
    matchID := "match" + strconv.Itoa(match_count)
    match.MatchID = matchID
    match_count++

    MatchAsBytes, _ := json.Marshal(match)
    err = stub.PutState(matchID, MatchAsBytes)

	if err != nil {
		return shim.Error(err.Error())
    }
    return shim.Success(nil)
}
func (t *SimpleChaincode) showAllMatch(stub shim.ChaincodeStubInterface) pb.Response {
    type myAllmatch struct {
        Matchs []MyMatch `json:"matchs"`
    }
    var allmatch myAllmatch
    var err error

    matchIterator, err := stub.GetStateByRange("match0", "match999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer matchIterator.Close()

    for matchIterator.HasNext() {
        data, err := matchIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var match MyMatch
        json.Unmarshal(ValAsBytes, &match)
		allmatch.Matchs = append(allmatch.Matchs, match)
    }
	allmatchAsBytes, _ := json.Marshal(allmatch)
	return shim.Success(allmatchAsBytes)
}

var power_count int = 0
func (t *SimpleChaincode) power(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

    var power Power
    if !(args[0] == "A" || args[0] == "B" || args[0] == "C") {
        return shim.Error("充電站格式錯誤!")
    }
    power.StationID = args[0]
    checkchargerID, err := strconv.Atoi(args[1])
    if err != nil || checkchargerID < 1 || checkchargerID > 30 {
        return shim.Error("充電樁格式錯誤!")
    }
    power.ChargerID = checkchargerID
    checkpower, err := strconv.Atoi(args[2])
    if err != nil || checkpower < 0 {
        return shim.Error("充電功率格式錯誤!")
    }
    power.Power = checkpower
    checkstate, err := strconv.Atoi(args[3])
    if err != nil || !(checkstate == 0 || checkstate == 1) {
        return shim.Error("快/慢充格式錯誤!")
    }
    power.State = checkstate
    checktimestamp, err := strconv.Atoi(args[4])
    if err != nil || checktimestamp < 1 || checktimestamp > 288 {
        return shim.Error("時間點格式錯誤!")
    }
    power.TimeStamp = checktimestamp
    powerID := "power" + strconv.Itoa(power_count)
    power.PowerID = powerID
    power_count++

    PowerAsBytes, _ := json.Marshal(power)
    err = stub.PutState(powerID, PowerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
    }
    return shim.Success(nil)
}
func (t *SimpleChaincode) showAllPower(stub shim.ChaincodeStubInterface) pb.Response {
    type myAllpower struct {
        Powers []Power `json:"powers"`
    }
    var allpower myAllpower
    var err error

    powerIterator, err := stub.GetStateByRange("power0", "power999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer powerIterator.Close()

    for powerIterator.HasNext() {
        data, err := powerIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var power Power
        json.Unmarshal(ValAsBytes, &power)
		allpower.Powers = append(allpower.Powers, power)
    }
	allpowerAsBytes, _ := json.Marshal(allpower)
	return shim.Success(allpowerAsBytes)
}
func (t *SimpleChaincode) showPowerbyCharger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
    if !(args[0] == "A" || args[0] == "B" || args[0] == "C") {
        return shim.Error("充電站格式錯誤!")
    }
    stationID := args[0]
    checkchargerID, err := strconv.Atoi(args[1])
    if err != nil || checkchargerID < 1 || checkchargerID > 30 {
        return shim.Error("充電樁格式錯誤!")
    }
    chargerID := checkchargerID

    powerIterator, err := stub.GetStateByRange("power0", "power999999999")
    if err != nil {
		return shim.Error(err.Error())
    }
    defer powerIterator.Close()

    for powerIterator.HasNext() {
        data, err := powerIterator.Next()
		if err != nil {
		    return shim.Error(err.Error())
		}
		ValAsBytes := data.Value

        var power Power
        json.Unmarshal(ValAsBytes, &power)
        if power.StationID == stationID && power.ChargerID == chargerID {
            return shim.Success(ValAsBytes)
        }
    }
    return shim.Error("No data for chargerID: " + stationID + strconv.Itoa(chargerID))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting EEP Chaincode: %s", err)
	}
}
