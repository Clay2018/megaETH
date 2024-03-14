实时同步区块数据, 将属于地址的充值交易插入数据库,并更新交易状态


查询所有钱包地址的余额,归集


BTC

ETH, usdt-erc20, usdc-erc20
TRX, USDT_TRC20
https://github.com/okx/js-wallet-sdk/blob/main/packages/coin-tron/tests/trx.test.ts


关于交互建议:
- 签名侧不需要数据库，离线和在线通过json文件联系
- 命令行程序实现命令
    -  生成key.json
    -  根据当前地址文件(address.json)和支持币种文件(token.json)，生成未签名文件(unsigned.json)
    -  根据当前文件(unsigned.json, key.json, gas_price.json), 生成已签名文件unsigned.json
    -  启动程序(若连接在线数据库，将插入充值数据)，广播
