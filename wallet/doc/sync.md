## 同步充值记录

### 配置需要同步充值记录的地址

- 少量地址可直接在./config/config.go中配置
- 大量地址可直接在启动程序的本地文件配置
  - eth地址文件名 "./eth_addresses.json"
  - btc地址文件名 "./btc_addresses.json"
  - trx地址文件名 "./trx_addresses.json"
   
### 创建数据库表

```sql
//insertUserSql := "INSERT INTO btc_deposit(addr, tx_hash, value, vout) VALUES(?, ?, ?, ?)"
CREATE TABLE `btc_deposit` (
    `addr` varchar(128) NOT NULL,
    `tx_hash` varchar(128) NOT NULL
    `value` int NOT NULL,
    `vout` int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `eth_deposit` (
    `addr` varchar(128) NOT NULL,
    `tx_hash` varchar(128) NOT NULL,
    `value` varchar(128) NOT NULL,
    `contract_addr` varchar(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `trx_deposit` (
    `addr` varchar(128) NOT NULL,
    `tx_hash` varchar(128) NOT NULL,
    `value` varchar(128) NOT NULL,
    `contract_addr` varchar(128) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```



### 配置数据库

./config/config.go
