package main 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/peer"
)

type SimpleChaincode struct {
}

type booking struct {
	ObjectType	string `json:"docType"` //docType is used to distinguish the various types of objects in state database  not sure why its needed
	AgentID		string `json:"agentId"`   //unique agent id
	BookingID	int	   `json:"bookingId"` //booking ID
	TotalBill	float32	`json:"totalBill"`  //no of days client will be staying		
	IsSettled	bool   `json:"isSettled"`  //is booking is pending or completed
	Payables	float32`json:"payables"`   //final payables to agent
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Printf("Error starting SimpleChaincode")
		}
}

//Init initializes chaincode
//==============================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Hello World")
	return shim.Success(nil)
}

//Invoke - entry points for invokations
//================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args:=stub.GetFunctionAndParameters()
	fmt.Println("invoke is running" +function)
	
	caller:= args[0]
	
	//divert different functions
	if function == "initBooking" && caller == "TA" {  //create a new booking
		return t.initBooking(stub,args)
	} else if function == "updateSettled" && caller == "HOTEL" {
		return t.updateSettled(stub,args)
	} else if function == "updatePayables" && caller == "CTAC" {
		return t.updatePayables(stub,args)
	} else if function == "readStatus" {
		return t.readStatus(stub, args)
	}
		
		
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
	
}

func(t *SimpleChaincode) updatePayables (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// arg0 will be caller
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting Booking ID")
	}
	
	bookingID := args[1]
	bookingAsBytes, err := stub.PutState(bookingID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} else if bookingAsBytes == nil {
		jsonResp := "{\"Error\":\"Booking ID does not exist: " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} 
	
	bookingTemp := booking{}
	
	err:= json.Unmarshal(bookingAsBytes, &bookingTemp)
	
	if err != nil {
		return shim.Error(err.Error())
	}
	
	
	if bookingTemp.IsSettled  != true {
		return shim.Error("payables cant be updated as settlement is pending for booking id" +bookingID)
	}
	
	bookingTemp.Payables = bookingTemp.TotalBill *0.02
	
	bookingAsJsonBytes, err := json.Marshal(bookingTemp)
	err = stub.PutState(bookingID, bookingAsJsonBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("payables are updated are successfully updated")
	
	
	
	
}


func (t *SimpleChaincode) updateSettled (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// arg0 will be caller arg1 booking id  arg2 totalbill
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting Booking ID")
	}
	
	bookingID := strconv.Atoi(args[1])
	bill := strconv.ParseFloat(args[2], 32)
	
	bookingAsBytes, err := stub.GetState(bookingID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} else if bookingAsBytes == nil {
		jsonResp := "{\"Error\":\"Booking ID does not exist: " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} 
	
	
	bookingTemp := booking{}
	err:= json.Unmarshal(bookingAsBytes, &bookingTemp)
	
	if err != nil {
		return shim.Error(err.Error())
	}
	
	bookingTemp.IsSettled = true
	bookingTemp.Totalbill = bill
	
	bookingAsJsonBytes, err := json.Marshal(bookingTemp)
	err = stub.PutState(bookingID, bookingAsJsonBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("settled status and total bill are successfully updated")
	
}


func (t *SimpleChaincode) initBooking (stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	// when calling init booking we need to pass all 5 arguments 
	//arg0 - caller type
	var err error
	
	agentID := strings.ToLower(args[1])
	bookingID := strconv.Atoi(args[2])
	totalBill := args[3])
	isSettled := args[4]
	payables := args[5]
	
	
	objectType := "booking"
	booking := &booking(objectType, agentID, bookingID, totalBill, isSettled, payables)
	bookingAsJsonBytes,err := json.Marshal(booking)
	if err != nil {
		return shim.Error(err.Error())
	}
	//put the marshal into chaincode
	err = stub.PutState(bookingID, bookingAsJsonBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// make index later
	
}

func (t *SimpleChaincode) readStatus (stub shim.ChaincodeStubInterface, args []string)  pb.response {
	
	var bookingid int
	var jsonResp string
	var err error
	var c booking
	
	bookingid = args[1]
	
	bookingAsBytes, err := stub.GetState(bookingid)
	
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + bookingid + "\"}"
		return shim.Error(jsonResp)
	} else if bookingAsBytes == nil {
		jsonResp = "{\"Error\":\"booking id does not exist: " + booking id  + "\"}"
		return shim.Error(jsonResp)
	}
	
	return shim.Success(bookingAsBytes)
	
}

