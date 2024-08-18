# TVBox Api Admin

一个用 Go 实现的简单 Web 服务器，用于远程管理tvbox生效的api，提供了一个简易管理界面，可以添加、编辑和删除链接，并更改链接的状态。


- **自定义域名绑定**: 支持通过命令行参数绑定自定义域名或 IP 地址。

## 运行

```
wget https://mirror.ghproxy.com/https://github.com/sligter/tvboxadm/releases/latest/download/tvbox_adm-linux-amd64
mv tvbox_adm-linux-amd64 /usr/bin/tvbox_adm
sudo chmod +x /usr/bin/tvbox_adm
tvbox_adm -port=2345 -domain=xxx.com(optional)
```
管理后台/admin
默认账户`admin`密码`admin`

首次进入需要至后台添加链接并激活即可

### 作为系统服务后台运行

创建系统服务文件:
```
sudo vim /etc/systemd/system/tvbox_adm.service
```
在文件中添加以下内容:
```
[Unit]
Description=TVBox Api Admin
After=network.target

[Service]
ExecStart=/usr/bin/tvbox_adm -port=2345 -domain=xxx.com(optional)
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
```

保存并关闭文件。
重新加载 systemd 配置:
```
sudo systemctl daemon-reload
```
启动服务:
```
sudo systemctl start tvbox_adm
```
设置开机自启:
```
sudo systemctl enable tvbox_adm
```
查看服务状态:
```
sudo systemctl status tvbox_adm
```

# docker部署
```
docker run -d --name tvbox_adm \
  -p 2345:2345 \
  --restart unless-stopped \
  bradleylzh/tvbox_adm:v1.0
```

### 绑定域名
```
docker run -d --name tvbox_adm \
  -p 2345:2345 \
  --restart unless-stopped \
  bradleylzh/tvbox_adm:v1.0 \
  --domain=tvbox.212138.xyz
```

### docker compose
```
version: '3'
services:
  tvbox_adm:
    image: bradleylzh/tvbox_adm:v1.0
    container_name: tvbox_adm
    ports:
      - "2345:2345"
    environment:
      - DOMAIN=xxx.com
    restart: unless-stopped
```