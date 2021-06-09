# G2WW (Grafana 2 Wechat Work)
Grafna webhook--企业微信机器人

## Build g2ww

```
go build
```

## Run g2ww

```
# 修改配置文件
# config.yml
# nohup ./g2ww & 或 pm2 start pm2.yml
```

## 填入企业微信机器人地址
参考
```
# webhook https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=f5134a5e-7413-41e3-90a5-3543f95f887b
# post http://localhost:2408/send/f5134a5e-7413-41e3-90a5-3543f95f887b
```
