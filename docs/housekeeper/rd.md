# housekeeper

## operator 功能的开发
```
1.配置生命周期管理
2.配置多版本管理
3.配置下发
```

## proxy 功能的开发
```
1.接收operator下发的数据
2.传递数据到agent
3.业务pod的驱离
```

## agent 功能的开发
```
1.接收proxy下发的数据
2.不可变以原子更新方式进行os升级，升级失败回滚
3.node采集信息
```