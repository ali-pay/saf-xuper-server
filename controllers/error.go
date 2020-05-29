package controllers

const (
	Unknown_Error = 0
)

var errMap = map[int]string{
	Unknown_Error: "未知错误",
}

func GetError(code int) string {
	return errMap[code]
}

/**
合约未部署 contract for account not confirmed
合约已存在 contract counter11 already exists
账户无效 The number of words in the Mnemonic sentence is not valid. It must be within [12, 15, 18, 21, 24]
合约账户无效 get account `XC1234567812345672@xuper` error: Key not found
合约账户已存在 account already exists
无法连接 connection refused
该账户没有足够的xuper NOT_ENOUGH_UTXO_ERROR
权限不够 RWACL_INVALID_ERROR
主链无法设置 xuper is forbidden
 */
