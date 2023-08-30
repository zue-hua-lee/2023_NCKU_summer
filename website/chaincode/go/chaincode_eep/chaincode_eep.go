package main

import (
	"encoding/json"
	"fmt"
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
	Date string              `json:"date"`
    ArrTime int              `json:"arrTime"`
    DepTime int              `json:"depTime"`
    ArrSoC int               `json:"arrSoC"`
    DepSoC int               `json:"depSoC"`
    Acdc int                 `json:"acdc"`
}
type MyMatch struct {
    MatchID string             `json:"matchID"`
    OfferID string             `json:"offerID"`
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
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
	} else if function == "login" {
		return t.login(stub, args)
    } else if function == "offer" {
		return t.offer(stub, args)
	} else if function == "match" {
		return t.match(stub, args)
	} else if function == "power" {
		return t.power(stub, args)
	} else if function == "getUserbyCar" {
        return t.getUserbyCar(stub, args)
	} else if function == "showCarbyUser" {
		return t.showCarbyUser(stub, args)
	} else if function == "showOfferbyID" {
		return t.showOfferbyID(stub, args)
	} else if function == "showMatchbyID" {
		return t.showMatchbyID(stub, args)
	} else if function == "showMatchbyUser" {
		return t.showMatchbyUser(stub, args)
	} else if function == "showPowerbyMatch" {
		return t.showPowerbyMatch(stub, args)
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
        return shim.Error("車牌號碼格式錯誤！車牌格式範例：EAB-1234")
    }
    checkcarID1 := strings.Split(checkcarID[0],"")
    if len(checkcarID1) != 3 || checkcarID[0][0] != 'E' || !unicode.IsUpper(rune(checkcarID[0][1])) || !unicode.IsUpper(rune(checkcarID[0][2])) {
        return shim.Error("車牌號碼格式錯誤！車牌號碼前三碼應為大寫英文字母，且第一字為E")
    }
    checkcarID2, err := strconv.Atoi(checkcarID[1])
    if err != nil || checkcarID2 < 1000 || checkcarID2 > 9999{
        return shim.Error("車牌號碼格式錯誤！車牌號碼後三碼應為四位數字。請重新填寫車排號碼")
    }
    user.CarID = args[0]
    if args[1] == "" {
        return shim.Error("尚未填寫車主名稱！請重新填寫車主名稱")
    }
    user.UserName = args[1]
    checkcapacity, err := strconv.Atoi(args[2])
    if err != nil || checkcapacity <= 0{
        return shim.Error("電池容量格式錯誤！請重新填寫電池容量")
    }
    user.Capacity = checkcapacity
    if args[3] == "" {
        return shim.Error("尚未填寫密碼！請重新填寫密碼")
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
                return shim.Error("密碼錯誤！請重新填寫登入資訊")
            }
        }
    }
	return shim.Error("車牌號碼錯誤！請重新填寫登入資訊")
}

var offer_count int = 0
func (t *SimpleChaincode) offer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

    var offer MyOffer
    offer.Date = args[0]
    checkarrtime, err := strconv.Atoi(args[1])
    if err != nil || checkarrtime < 1 || checkarrtime > 288{
        return shim.Error("開始充電時間格式錯誤！請重新填寫開始充電時間")
    }
    offer.ArrTime = checkarrtime
    checkdeptime, err := strconv.Atoi(args[2])
    if err != nil || checkdeptime < 1 || checkdeptime > 288 || checkdeptime <= offer.ArrTime {
        return shim.Error("結束充電時間格式錯誤！請重新填寫結束充電時間")
    }
    offer.DepTime = checkdeptime
    checkarrSoC, err := strconv.Atoi(args[3])
    if err != nil || checkarrSoC < 0 || checkarrSoC > 100{
        return shim.Error("目前電量格式錯誤！請重新填寫目前電量")
    }
    offer.ArrSoC = checkarrSoC
    checkdepSoC, err := strconv.Atoi(args[4])
    if err != nil || checkdepSoC < 0 || checkdepSoC > 100 || checkdepSoC <= offer.ArrSoC {
        return shim.Error("目標電量格式錯誤！請重新填寫目標電量")
    }
    offer.DepSoC = checkdepSoC
    checkacdc, err := strconv.Atoi(args[5])
    if err != nil || !(checkacdc == 1 || checkacdc == 2) {
        return shim.Error("充電方式格式錯誤！請重新填寫充電方式")
    }
    offer.Acdc = checkacdc
    if args[6] == "" {
		return shim.Error("尚無使用者登入！請先登入")
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

var match_count int = 0
func (t *SimpleChaincode) match(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

    var match MyMatch
    if !(args[0] == "A" || args[0] == "B" || args[0] == "C") {
        return shim.Error("充電站格式錯誤!")
    }
    match.StationID = args[0]
    checkchargerID, err := strconv.Atoi(args[1])
    if err != nil || checkchargerID < 1 || checkchargerID > 30 {
        return shim.Error("充電樁格式錯誤!")
    }
    match.ChargerID = checkchargerID
    checkmaxSoC, err := strconv.Atoi(args[2])
    if err != nil || checkmaxSoC < 0 || checkmaxSoC > 100 {
        return shim.Error("離場電量格式錯誤!")
    }
    match.MaxSoC = checkmaxSoC
    checkprice, err := strconv.Atoi(args[3])
    if err != nil || checkprice < 0 {
        return shim.Error("電價格式錯誤!")
    }
    match.Price = checkprice
    if args[4] == "" {
		return shim.Error("尚未提出充電申請!")
    }
    match.OfferID = args[4]
    matchID := "match" + strconv.Itoa(match_count)
    match.MatchID = matchID
    match_count++

    MatchAsBytes, _ := json.Marshal(match)
    err = stub.PutState(matchID, MatchAsBytes)
	if err != nil {
		return shim.Error(err.Error())
    }
    return shim.Success([]byte(matchID))
}
var power_count int = 0
func (t *SimpleChaincode) power(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

    var powers []Power
    err = json.Unmarshal([]byte(args[0]), &powers)
    if err != nil {
        return shim.Error(err.Error())
    }

    for _, p := range powers {
        var power Power
        if !(p.StationID == "A" || p.StationID == "B" || p.StationID == "C") {
            return shim.Error("充電站格式錯誤!")
        }
        power.StationID = p.StationID
        if p.ChargerID < 1 || p.ChargerID > 30 {
            return shim.Error("充電樁格式錯誤!")
        }
        power.ChargerID = p.ChargerID
        if p.Power < 0 {
            return shim.Error("充電功率格式錯誤!")
        }
        power.Power = p.Power
        if p.State != -1 && (p.State < 0 || p.State > 100) {
            return shim.Error("配對狀態格式錯誤!")
        }
        power.State = p.State
        if p.TimeStamp < 1 || p.TimeStamp > 288 {
            return shim.Error("時間點格式錯誤!")
        }
        power.TimeStamp = p.TimeStamp
        powerID := "power" + strconv.Itoa(power_count)
        power.PowerID = powerID
        power_count++
    
        PowerAsBytes, _ := json.Marshal(power)
        err = stub.PutState(powerID, PowerAsBytes)
        if err != nil {
            return shim.Error(err.Error())
        }
    }
    return shim.Success(nil)
}

func (t *SimpleChaincode) getUserbyCar(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	return shim.Error("getUserbyCar ERROR!")
}
func (t *SimpleChaincode) showCarbyUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
    return shim.Success([]byte(user.CarID))
}
func (t *SimpleChaincode) showOfferbyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
    var offer MyOffer
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    if args[0] == "" {
		return shim.Error("尚無提出充電申請!")
    }
    offerID := args[0]
    
    offerAsBytes, err := stub.GetState(offerID)
    if err != nil {
            return shim.Error("Failed to find offer - " + offerID)
    }
    json.Unmarshal(offerAsBytes, &offer)
    if offer.OfferID != offerID {
            return shim.Error("Offer does not exist - " + offerID)
    }
    return shim.Success(offerAsBytes)
}
func (t *SimpleChaincode) showMatchbyID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
    var match MyMatch
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    if args[0] == "" {
		return shim.Error("尚無提出充電配對!")
    }
    matchID := args[0]
    
    matchAsBytes, err := stub.GetState(matchID)
    if err != nil {
            return shim.Error("Failed to find match - " + matchID)
    }
    json.Unmarshal(matchAsBytes, &match)
    if match.MatchID != matchID {
            return shim.Error("match does not exist - " + matchID)
    }
    return shim.Success(matchAsBytes)
}
func (t *SimpleChaincode) showMatchbyUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    if args[0] == "" {
		return shim.Error("尚無使用者登入!")
    }
    userID := args[0]

    type myAllmatch struct {
        Matchs []MyMatch `json:"matchs"`
    }
    var allmatch myAllmatch

    matchIterator, err := stub.GetStateByRange("match0", "match999999999")
    if err != nil {
        return shim.Error(err.Error())
    }
    defer matchIterator.Close()

    if !matchIterator.HasNext() {
        return shim.Error("No matches found")
    }
    for matchIterator.HasNext() {
        data, err := matchIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        ValAsBytes := data.Value

        var match MyMatch
        json.Unmarshal(ValAsBytes, &match)

        var offer MyOffer
        offerAsBytes, err := stub.GetState(match.OfferID)
        if err != nil {
                return shim.Error("Failed to find offer - " + match.OfferID)
        }
        json.Unmarshal(offerAsBytes, &offer)
        if offer.OfferID != match.OfferID {
                return shim.Error("Offer does not exist - " + match.OfferID)
        }
        if offer.UserID == userID{
            allmatch.Matchs = append(allmatch.Matchs, match)
        }
    }
    allmatchAsBytes, _ := json.Marshal(allmatch)
    return shim.Success(allmatchAsBytes)
}

func (t *SimpleChaincode) showPowerbyMatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

    var match MyMatch
    if args[0] == "" {
        return shim.Error("尚無提出充電配對!")
    }
    matchID := args[0]
    
    matchAsBytes, err := stub.GetState(matchID)
    if err != nil {
        return shim.Error("Failed to find match - " + matchID)
    }
    json.Unmarshal(matchAsBytes, &match)
    if match.MatchID != matchID {
        return shim.Error("Match does not exist - " + matchID)
    }

    var offer MyOffer
    offerAsBytes, err := stub.GetState(match.OfferID)
    if err != nil {
        return shim.Error("Failed to find offer - " + match.OfferID)
    }
    json.Unmarshal(offerAsBytes, &offer)
    if offer.OfferID != match.OfferID {
        return shim.Error("Offer does not exist - " + match.OfferID)
    }

    type myAllpower struct {
        Powers []Power `json:"powers"`
    }
    var allpower myAllpower
    
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
        if power.StationID == match.StationID && power.ChargerID == match.ChargerID && power.TimeStamp >= offer.ArrTime && power.TimeStamp < offer.DepTime {
            allpower.Powers = append(allpower.Powers, power)
        }
    }
    allpowerAsBytes, _ := json.Marshal(allpower)
	return shim.Success(allpowerAsBytes)
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting EEP Chaincode: %s", err)
	}
}
