# invoices 

`invoices` 用來處理「財政部電子發票平台」寄送的「消費記錄彙整」。

## 主要功能

1. 可設定 `big5`「消費記錄彙整」轉為 `utf-8`。
2. 透過設定，輸入檔格式可為：`.csv`、`.jsn` 或 `.json`、`.xml`，輸出檔格式可為：`.csv`、`.jsn` 或 `.json`、`.xml`。
3. 消費記錄均於 `sqlite3` 資料檔案中備份。

## 使用說明

### 命令行參數(command-line parameters)

```shell
invs [global options]
```
### 功能(`global options`)
1. `--case value`,` -c value`  case-options filename (default: "./inp/case.ini")
2. `--initial`, `-i`           initalizing enviroment of applicaton to inital state
3. `--verbose`, `-b`           verbose output
4. `--help`,` -h`              show help
5. `--version`, `-v`           print the version