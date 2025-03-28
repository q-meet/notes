# git

## 配置仓库ssh时出现错误尝试解决方法之1

放在 .ssh目录文件上创建一个 config 文件 指定ssh密钥位置

```config
Host *
IdentityFile ～/.ssh/id_rsa
HostKeyAlgorithms ssh-rsa
PubkeyAcceptedKeyTypes ssh-rsa
```

## 语句操作

### 单独合并提交

git cherry-pick commit_id

### 合并某个分支上的一系列commits

在一些特性情况下，合并单个commit并不够，你需要合并一系列相连的commits。这种情况下就不要选择cherry-pick了，rebase 更适合。

```git checkout -b <newBranchName> <to-commit-id>```  
创建一个新的分支，指明新分支的最后一个commit

``git rebase --onto <branchName> <from-commit-id>``
变基这个新的分支到最终要合并到的分支，指明从哪个特定的commit开始
如： 还以上面的为例，假设你需要合并feature分支的commit 76cada ~62ecb3 到master分支。 首先需要基于feature创建一个新的分支，并指明新分支的最后一个commit：

``git checkout -b <newbranch> 62ecb3``
然后，rebase这个新分支的commit到master（–onto master）。76cada^ 指明你想从哪个特定的commit开始。

git rebase --onto master 76cada^
得到的结果就是feature分支的commit 76cada ~62ecb3 都被合并到了master分支。

再合并的过程中可能出现冲突，出现冲突，必须手动解决后，然后 运行git rebase --continue。

系统对冲突的提示：
fix conflicts and then run "git rebase --continue"
use "git rebase --skip" to skip this patch
use "git rebase --abort" to check out the original branch
