##### 服务器需要先配置好ssh免密登录
##### 参数：
- -remoteDir string  服务器同步目录
- -remoteHost string 服务器地址
- -remoteUser string 登录用户名
- -watchedDir string 受监控目录
##### eg:
```bash
-watchedDir="C:\Test" -remoteHost="0.0.0.0" -remoteUser="root" -remoteDir="~/testSync"
```

