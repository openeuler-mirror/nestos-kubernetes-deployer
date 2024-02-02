# 全局配置文件说明

``` shell
persistdir: /etc/nkd            # 文件存储路径，包括全局配置文件、集群配置文件以及证书文件等。默认路径：/etc/nkd           
bootstrapurl:
  bootstrap_ign_host: ""        # 用于引导ignition文件的http服务的域名或者IP地址。默认为宿主机ip地址
  bootstrap_ign_port: "9080"    # 用于引导ignition文件的http服务的端口，默认端口号为9080
```  