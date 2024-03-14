### 命令行操作
```shell
# 编译命令行工具
go mod tidy
cd ./cmd
go build cmd.go
mv cmd wallet

# --help
./wallet trx --help

# 产生地址, 生成"trx_key.json"和"trx_insert_addr.sql"
./wallet trx generate 2

# 查看地址余额, 生成"trx_insert_balance.sql"
# addresses.json 参照 "doc/trx_addresses.json.example"
 ./wallet trx balance addresses.json

# 查看当前gasPrice, 控制台打印出当前 gasPrice, 此值跟faucet有关
 ./wallet trx gasprice

# 生成未签名文件"trx_unsigned_non_addr.json"和"trx_unsigned.json"文件
# "trx_unsigned_non_addr.json"为生成过程中出错的地址
# help info: wallet trx collect [token] [addresses.json] [gasPrice] [toAddr]
./wallet trx collect trx trx_addrs.json 2000000 TPfHmimxS2drgvxzDiAbgFahrcsBoZG9EQ
./wallet trx collect usdt trx_addrs.json 50000000 TPfHmimxS2drgvxzDiAbgFahrcsBoZG9EQ

# 生成未签名文件"eth_unsigned_non_addr.json"和"eth_unsigned.json"文件
# "eth_unsigned_non_addr.json"为生成过程中出错的地址
# help info: wallet faucetETH [token] [addresses.json] [gasPrice] [fromAddr]
./wallet trx faucet trx trx_addrs.json 50000000 TPfHmimxS2drgvxzDiAbgFahrcsBoZG9EQ

# 签名, 生成"eth_signed.json"和"eth_signed_non_addr.json"
# "eth_signed_non_addr.json"为生成过程中出错的地址, "eth_signed.json"中包含广播数组
# help info: wallet eth sign [eth_unsigned.json] [eth_key.json]
./wallet eth sign eth_unsigned.json eth_key.json 

# 启动服务器
go run main.go

# 广播
curl --location 'http://127.0.0.1:10000/eth/sendRawTransaction' \
--header 'Content-Type: application/json' \
--data '{"raw_trans": "f86f80850c3d73f2ac825208941d8b917c99ff73c8fdc8568e53654914eeeba28e871f9ae1606652a0808401546d72a06246d445816b8c0c48911fe73152638848c64c6d87542584a6a9bd06728b946fa070adba7320fa227008c82002854c6a9fc40407c2d7ff0c3667d5cad1e4a4dcfa"}'
```
