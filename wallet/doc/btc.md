### 命令行操作
```shell
# 编译命令行工具
go mod tidy
cd ./cmd
go build cmd.go
mv cmd wallet

# --help
./wallet btc --help

# 产生地址, 生成"btc_key.json"和"btc_insert_addr.sql"
./wallet btc generate 2

# 查看地址余额, 生成"btc_insert_balance.sql"
# addresses.json 参考 "btc_addresses.json.example
 ./wallet btc balance addresses.json

# 查看当前gasPrice, 控制台打印出当前 gasPrice
 ./wallet btc gasprice

# 生成未签名文件"btc_unsigned_non_addr.json"和"btc_unsigned.json"文件
# "btc_unsigned_non_addr.json"为生成过程中出错的地址
# help info: wallet btc collect [token] [addresses.json] [gasPrice] [toAddr]
./wallet btc collect btc addresses.json 100 bc1pt0asquzczqgem4g3cwtnhtfashyk6h3e76fwe40nvhqwdnjusxzqn0s8dq

# 签名, 生成"btc_signed.json"和"btc_signed_non_addr.json"
# "btc_signed_non_addr.json"为生成过程中出错的地址, "btc_signed.json"中包含广播数据
# help info: wallet btc sign [btc_unsigned.json] [btc_key.json]
./wallet btc sign btc_unsigned.json btc_key.json 

# 启动服务器
go run main.go

# 广播
curl --location 'http://127.0.0.1:10000/btc/sendRawTransaction' \
--header 'Content-Type: application/json' \
--data '{"raw_trans": "020000000001012a0726b649613a24e8188889855836f63e331499eb121838614784300a97ae291500000000fdffffff014a01000000000000225120c461dc300140be0ea291c8b0d6302fcd66fb80e7d42781072db46bbd2e05aaab0340cf0878e35fd6c93ed48d7595b845ba73b4fe0249abccda921d13d057b07fd8eb7040d20c8224676a9c135b68ee9a74074e2e260e47a8b24ac9336209a4d159577e2029d9004f20338c77f5ee8d2a256f6d0e47304c587722c601f048925b69063180ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800387b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a226c69676f222c22616d74223a2231303030303030227d6821c029d9004f20338c77f5ee8d2a256f6d0e47304c587722c601f048925b6906318000000000"}'
```
