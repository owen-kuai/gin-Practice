# gin practice 介绍
   
#构建
本地构建可以执行 make build 命令

## 本地启动：
配置ide 添加启动参数
```bash
run -c /var/lib/configs # 指定自己的配置文件地址
```

# 配置说明
```bash
host: "0.0.0.0"  # 绑定的ip 默认为 0.0.0.30
port: 8899       # 默认监听端口
enable_swagger: true  # 是否开启swagger 页面，
logger:
  path: /var//logs  # 日志文件夹
  maxSize: 512            # 单个日志文件的大小 单位 MB
  maxAge: 15              #  日志文件的保存时间 单位 天
  maxBackups: 10          # 日志文件的备份数量
  compress: false         # 是否压缩
```