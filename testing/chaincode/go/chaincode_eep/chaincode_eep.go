package main

import (
    "fmt"
    "encoding/json"
    "strconv"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}
type Offer struct {
    OfferID string           `json:"offerID"`
    CarNum int               `json:"carNum"`
    ArrTime int              `json:"arrTime"`
    DepTime int              `json:"depTime"`
    ArrSoC int               `json:"arrSoC"`
    DepSoC int               `json:"depSoC"`
    Acdc int                 `json:"acdc"`
    Capacity int             `json:"capacity"`
	Location_x float64       `json:"location_x"`
	Location_y float64       `json:"location_y"`
}
type Match struct {
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
var offer_count int = 0
func (t *SimpleChaincode) offer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

    var offer Offer
    checkcarnum, err := strconv.Atoi(args[0])
    if err != nil || checkcarnum < 1 {
        return shim.Error("電動汽車編號格式錯誤!")
    }
    offer.CarNum = checkcarnum
    checkarrtime, err := strconv.Atoi(args[1])
    if err != nil || checkarrtime < 1 || checkarrtime > 288{
        return shim.Error("開始充電時間格式錯誤!")
    }
    offer.ArrTime = checkarrtime
    checkdeptime, err := strconv.Atoi(args[2])
    if err != nil || checkdeptime < 1 || checkdeptime > 288 || checkdeptime <= offer.ArrTime {
        return shim.Error("結束充電時間格式錯誤!")
    }
    offer.DepTime = checkdeptime
    checkarrSoC, err := strconv.Atoi(args[3])
    if err != nil || checkarrSoC < 0 || checkarrSoC > 100{
        return shim.Error("開始電池狀態格式錯誤!")
    }
    offer.ArrSoC = checkarrSoC
    checkdepSoC, err := strconv.Atoi(args[4])
    if err != nil || checkdepSoC < 0 || checkdepSoC > 100 || checkdepSoC <= offer.ArrSoC {
        return shim.Error("結束電池狀態格式錯誤!")
    }
    offer.DepSoC = checkdepSoC
    checkacdc, err := strconv.Atoi(args[5])
    if err != nil || !(checkacdc == 1 || checkacdc == 2) {
        return shim.Error("快/慢充格式錯誤!")
    }
    offer.Acdc = checkacdc
    checkcapacity, err := strconv.Atoi(args[6])
    if err != nil || checkcapacity < 0 {
        return shim.Error("電池容量格式錯誤!")
    }
    offer.Capacity = checkcapacity
    checklocation_x, err := strconv.ParseFloat(args[7], 64)
    if err != nil {
        return shim.Error("起始位置x格式錯誤!")
    }
    offer.Location_x = checklocation_x
    checklocation_y, err := strconv.ParseFloat(args[8], 64)
    if err != nil {
        return shim.Error("起始位置y格式錯誤!")
    }
    offer.Location_y = checklocation_y

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
    type Alloffer struct {
        Offers []Offer `json:"offers"`
    }
    var alloffer Alloffer
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

        var offer Offer
        json.Unmarshal(ValAsBytes, &offer)
		alloffer.Offers = append(alloffer.Offers, offer)
    }

	allofferAsBytes, _ := json.Marshal(alloffer)
	return shim.Success(allofferAsBytes)
}

var match_count int = 0
func (t *SimpleChaincode) match(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

    var match Match
    if !(args[0] == "A" || args[0] == "B" || args[0] == "C") {
        return shim.Error("充電站格式錯誤!")
    }
    match.StationID = args[0]
    checkchargerID, err := strconv.Atoi(args[1])
    if err != nil || checkchargerID < 1 {
        return shim.Error("充電樁格式錯誤!")
    }
    match.ChargerID = checkchargerID
    checkmaxSoC, err := strconv.Atoi(args[2])
    if err != nil || checkmaxSoC < 0 || checkmaxSoC > 100 {
        return shim.Error("最高可充電池狀態格式錯誤!")
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
    return shim.Success(nil)
}
func (t *SimpleChaincode) showAllMatch(stub shim.ChaincodeStubInterface) pb.Response {
    type Allmatch struct {
        Matchs []Match `json:"matchs"`
    }
    var allmatch Allmatch
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

        var match Match
        json.Unmarshal(ValAsBytes, &match)
		allmatch.Matchs = append(allmatch.Matchs, match)
    }
	allmatchAsBytes, _ := json.Marshal(allmatch)
	return shim.Success(allmatchAsBytes)
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
        if !(p.State == 0 || p.State == 1) {
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

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting EEP Chaincode: %s", err)
	}
}
