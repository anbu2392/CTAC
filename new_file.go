package main 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type booking struct {
	ObjectType	string `json:"docType"` //docType is used to distinguish the various types of objects in state database  not sure why its needed
	AgentID		string `json:"agentId"`   //unique agend id
	BookingID	int	   `json:"bookingId"` //booking ID
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
	return shim.Success(nil)
}

//Invoke - entry points for invokations
//================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
	function,args:=stub.GetFunctionAndParameters()
	fmt.Println("invoke is running" +function)
	
	//divert different functions
	if function == "initBooking"{  //create a new booking
		return t.initBooking(stub,args)
	}
	
}

func (t *SimpleChaincode) initBooking(stub shim.ChaincodeStubInterface) pb.Response{
	var err error
	
}