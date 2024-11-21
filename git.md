# git

## 配置仓库ssh时出现错误尝试解决方法之1

放在 .ssh目录文件上创建一个 config 文件 指定ssh密钥位置

```config
Host *
IdentityFile ～/.ssh/id_rsa
HostKeyAlgorithms ssh-rsa
PubkeyAcceptedKeyTypes ssh-rsa
```

git 合并 单提交

```cmd
git cherry-pick commit-id
```

