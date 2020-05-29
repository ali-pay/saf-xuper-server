package xkernel

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"github.com/xuperchain/xuperchain/core/pb"
)

const (
	NEW     = "NewAccount"
	ACCOUNT = "SetAccountAcl"
	METHOD  = "SetMethodAcl"
)

type Acl struct {
	xchain.Xchain `json:"-"`

	Pm struct {
		Rule        int     `json:"rule"`
		AcceptValue float32 `json:"acceptValue"`
	} `json:"pm"`
	AksWeight map[string]float32 `json:"aksWeight"`
}

func InitAcl(account *account.Account, node, bcName, contractAccount string) *Acl {
	commConfig := config.GetInstance()

	return &Acl{
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
		},
	}
}

func (c *Acl) AclIR(action, contractName, contractMethod string, aks map[string]float32) *pb.InvokeRequest {
	args := make(map[string][]byte)

	switch action {
	case NEW:
		args["account_name"] = []byte(c.ContractAccount)

	case ACCOUNT:
		args["account_name"] = []byte(c.ContractAccount)

	case METHOD:
		args["contract_name"] = []byte(contractName)
		args["method_name"] = []byte(contractMethod)
	}

	var aksWeight string
	format := `"%s": %.1f,`
	//遍历所有地址
	for k, v := range aks {
		aksWeight += fmt.Sprintf(format, k, v)
	}
	//去掉最后的逗号
	aksWeight = aksWeight[:strings.LastIndex(aksWeight, ",")]

	acl := `
        {
            "pm": {
                "rule": 1,
                "acceptValue": 1.0
            },
            "aksWeight": {
                %s
            }
        }
        `
	acl = fmt.Sprintf(acl, aksWeight)
	args["acl"] = []byte(acl)

	//for k, v := range args {
	//	fmt.Println(k,string(v))
	//}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: action,
		Args:       args,
	}
}

func (c *Acl) AclDoit(action, contractName, contractMethod string, aks map[string]float32) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.Pre(action, contractName, contractMethod, aks)
	if err != nil {
		log.Printf("Acl preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.Post(preSelectUTXOResponse)
}

func (c *Acl) Pre(action, contractName, contractMethod string, aks map[string]float32) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := c.AclIR(action, contractName, contractMethod, aks)
	invokeRequests = append(invokeRequests, invokeRequest)

	var authRequires []string
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

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

func (c *Acl) Post(preExeWithSelRes *pb.PreExecWithSelectUTXOResponse) (string, error) {
	// populates fields
	var authRequires []string
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	c.Initiator = c.Account.Address
	c.AuthRequire = authRequires
	c.InvokeRPCReq = nil
	c.PreSelUTXOReq = nil
	c.Fee = strconv.Itoa(int(preExeWithSelRes.Response.GasUsed))
	c.Amount = "0"

	return c.GenCompleteTxAndPost(preExeWithSelRes)
}

func (c *Acl) CreateContractAccount() (gas int64, acl *Acl, txid string, err error) {
	aks := map[string]float32{
		c.Account.Address: 1,
	}

	// preExe
	resp, err := c.Pre(NEW, "", "", aks)
	if err != nil {
		log.Printf("Acl Pre failed, err: %v", err)
		return
	}

	// post
	txid, err = c.Post(resp)
	if err != nil {
		log.Printf("Acl Post failed, err: %v", err)
		return
	}

	//合约应答数据
	//var datas []string
	//for _, v := range resp.Response.Responses {
	//	//记录合约错误
	//	if v.Status != 200 {
	//		datas = append(datas, v.String())
	//		continue
	//	}
	//
	//	//记录数据
	//	body := string(v.Body)
	//	body = strings.ReplaceAll(body, " ", "")
	//	body = strings.ReplaceAll(body, "\n", "")
	//
	//	datas = append(datas, body)
	//}
	//str := fmt.Sprint(datas)

	acl = new(Acl)
	err = json.Unmarshal(resp.Response.Responses[0].Body, acl)
	if err != nil {
		log.Printf("unmarshal acl fail, err: %v", err)
		return
	}
	return resp.Response.GasUsed, acl, txid, err
}
