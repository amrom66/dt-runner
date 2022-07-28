# 介绍

dt-runner是一个集中式的CI触发平台。在传统的gitlab ci中，研发人员使用.gitlab-ci.yaml文件控制项目的ci流程，但是，这种流程存在权限的管控漏洞，研发人员可以在CI脚本中写入任意内容，从而导致 CI流程的不可控。dt-runner是一个集中式的CI触发工具，其禁用了项目的CI文件，改有集中控制，加强项目的可控性。

* 监听gitlab所有事件
* 创建CRD
* 运行k8s pod
* 采集运行结果
* 发送邮件

## 使用手册

* 准备环境

1. golang 1.18

2. k8s 1.21+

3. kubectl

* 编译

```shell
make buildamd
```

* 运行

```shell
sudo install dt-runner-amd64 /usr/local/bin/dt-runner
mv dt-runner.service /etc/systemd/system/dt-runner.service

adduser amrom

sudo chown -R amrom:amrom /etc/systemd/system/dt-runner.service
sudo chown -R amrom:amrom /usr/local/bin/dt-runner

systemctl daemon-reload
systemctl start dt-runner

journalctl -u dt-runner
```

* 注册到gitlab系统服务

在`gilab`管理员页面，找到`system hook`，添加一个新的`hook`，地址为上一步的`dt-runner`的地址

* 验证可用

vrmanager.yaml

```yaml
apiVersion: apps.dtwave.com/v1
kind: Ci
metadata:
  name: vrmanager
spec:
  model: "model-sample"
  repo: "http://192.168.90.154:32578/root/vrmanager.git"
  branch: "main"
  term:
    schedule: '28 23 1/7 * *'
    events:
    - push
    - commit
  variables:
    PROJECT_NAME: vrmanager
    SHUXI_VESION: 0.0.1
```

```shell
kubectl apply -f vrmanager.yaml
```

注意：`repo`应该是`gitlab`的仓库地址

正常`commit`信息到`repo`，观察`dt-runner`的日志，如果有新的`commit`，则会创建新的`pod`
