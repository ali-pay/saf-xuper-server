# 如果编译发生错误

**1. jason-cn-dev：将该字段改成：ali-pay**

**2. xupercc：将该字段改成：saf-xuper-server**

# 超级链的后端

## [查看API文档](超级链API文档_v0629.md)

## 使用说明

环境配置
```bash
sudo apt-get install docker.io
sudo apt-get install gcc
```

修改权限
```bash
sudo chmod 777 xdev/*.sh
sudo chmod 777 xdev/xdev
```

编译
```go
go build -mod vendor
//此时会在当前目录下生成xupercc可执行文件
```

运行
```bash
./xupercc

或者后台执行

nohop ./xupercc &
```

配置说明

conf/app.ini：这是后端服务的配置
（如无必要，不要改动）

```ini
[app]
# 后端的运行模式，可选：debug or release
app_mode = release
# 运行时临时文件的存放位置
root_path = runtime/

[server]
# http协议的类型，可选：http or https
protocol = http
# 后台监听的端口
http_port = 8080

[code]
# 合约代码的最大大小，单位MB
file_max_size = 2
# 合约代码支持的类型
file_exts = cc,go
# 合约代码的存放位置
code_path = contract_code/
# 合约代码编译后的存放位置
wasm_path = contract_wasm/

[log]
# 日志文件的存放位置
file_path = logs/
# 日志文件名
file_name = gin.log
# 路由日志文件名
router_file = router.log
# 运行时日志文件名
runtime_file = runtime.log

[req]
# 助记词的类型，可选：1-2 中/英文
language = 1
# 助记词的长度，可选：1-3 8/16/24位
strength = 1

[cache]
# 区块和交易的缓存大小，建议不要超过15
size = 10
```

conf/sdk.yaml：这是sdk的配置，用来跟链交互用的
（如无必要，不要改动）

```yaml
#背书节点的地址
endorseServiceHost: "127.0.0.1:37101"

complianceCheck:
  # 是否需要进行合规性背书
  isNeedComplianceCheck: false
  # 是否需要支付手续费
  isNeedComplianceCheckFee: false
  # 手续费
  complianceCheckEndorseServiceFee: 0
  # 收手续费的地址
  complianceCheckEndorseServiceFeeAddr: WwLgfAatHyKx2mCJruRaML4oVf7Chzp42
  # 给手续费的地址
  complianceCheckEndorseServiceAddr: ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt

# 加密算法配置, 国密:gm
crypto: "xchain"

#创建平行链所需要的最低费用，需要跟节点启动时的配置一样
minNewChainAmount: "100"
```
