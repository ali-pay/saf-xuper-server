package controllers

const (
	ErrorUnknown                        = "unknown error"
	ErrorContractForAccountNotConfirmed = "contract for account not confirmed"
	ErrorContractAlreadyExists          = "contract * already exists"
	ErrorNotEnoughUtxo                  = "NOT_ENOUGH_UTXO_ERROR"
	ErrorRwaclInvalid                   = "RWACL_INVALID_ERROR"
	ErrorConnectionRefused              = "connection refused"
	ErrorAccountAlreadyExists           = "account already exists"
)

var errMap = map[string]string{
	ErrorUnknown:                        "未知错误",
	ErrorContractForAccountNotConfirmed: "合约未部署",
	ErrorContractAlreadyExists:          "合约已存在",
	ErrorNotEnoughUtxo:                  "账户余额不够",
	ErrorRwaclInvalid:                   "账户权限不够",
	ErrorConnectionRefused:              "无法链接或请求有误导致链接被拒接",
	ErrorAccountAlreadyExists:           "合约账户已存在",
}

func GetError(code string) string {
	return errMap[code]
}
