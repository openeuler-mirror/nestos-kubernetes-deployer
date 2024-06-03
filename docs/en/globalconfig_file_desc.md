# Global config file description

``` shell
persistdir: /etc/nkd            # File storage path, including global configuration files, cluster configuration files, certificate files, etc. Default path: /etc/nkd       
bootstrapurl:
  bootstrapIgnHost: ""        # Ignition service address (domain name or IP, usually NKD operating environment)
  bootstrapIgnPort: "9080"    # Ignition service port (default 9080, you need to open the firewall port yourself)
```  