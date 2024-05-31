# 全局配置文件说明

``` shell
persistdir: /etc/nkd            # 文件存储路径，包括全局配置文件、集群配置文件以及证书文件等。默认路径：/etc/nkd           
bootstrapurl:
  bootstrapIgnHost: ""        # 点火服务地址（域名或ip，一般为NKD运行环境）
  bootstrapIgnPort: "9080"    # 点火服务端口（默认9080，需自行开放防火墙端口）
```  