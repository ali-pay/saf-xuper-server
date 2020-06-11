package xkernel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/jason-cn-dev/xuper-sdk-go/account"
	"github.com/jason-cn-dev/xuper-sdk-go/config"
	"github.com/jason-cn-dev/xuper-sdk-go/pb"
	"github.com/jason-cn-dev/xuper-sdk-go/xchain"
)

const (
	UPGEADE = "Upgrade"
	DEPLOY  = "Deploy"
)

type Contract struct {
	ContractName string
	xchain.Xchain
}

func InitContract(account *account.Account, node, bcName, contractAccount, contractName string) *Contract {
	commConfig := config.GetInstance()

	return &Contract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
		},
	}
}

func convertToXuperContractArgs(args map[string]string) map[string][]byte {
	argmap := make(map[string][]byte)
	for k, v := range args {
		argmap[k] = []byte(v)
	}
	return argmap
}

func (c *Contract) ContractIR(action, codepath, runtime string, arg map[string]string) *pb.InvokeRequest {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	contractCode, err := ioutil.ReadFile(codepath)
	if err != nil {
		log.Printf("get wasm contract code error: %v", err)
		return nil
	}
	desc := &pb.WasmCodeDesc{
		Runtime: runtime,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(c.ContractAccount),
		"contract_name": []byte(c.ContractName),
		"contract_code": contractCode,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
	}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: action,
		Args:       args,
	}
}

func (c *Contract) ContractDoit(action, codepath, runtime string, arg map[string]string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.Pre(action, codepath, runtime, arg)
	if err != nil {
		log.Printf("Contract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.Post(preSelectUTXOResponse)
}

func (c *Contract) Pre(action, codepath, runtime string, arg map[string]string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := c.ContractIR(action, codepath, runtime, arg)
	invokeRequests = append(invokeRequests, invokeRequest)

	var authRequires []string
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck{
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}
	fmt.Println(authRequires)
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}
	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee),
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}

func (c *Contract) Post(preExeWithSelRes *pb.PreExecWithSelectUTXOResponse) (string, error) {
	// populates fields
	var authRequires []string
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck{
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}
	c.Initiator = c.Account.Address
	c.AuthRequire = authRequires
	c.InvokeRPCReq = nil
	c.PreSelUTXOReq = nil
	c.Fee = strconv.Itoa(int(preExeWithSelRes.Response.GasUsed))
	c.TotalToAmount = "0"

	return c.GenCompleteTxAndPost(preExeWithSelRes, "")
}
