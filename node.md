# node

## 要在CentOS上安装nvm（Node Version Manager）

首先，确保您的系统已安装curl和tar。如果没有安装，请运行以下命令进行安装：

```shell
sudo yum install curl
sudo yum install tar
```

接下来，使用curl下载nvm安装脚本。您可以在 <https://github.com/nvm-sh/nvm> 上找到最新的nvm版本。

运行以下命令下载nvm：

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/version/install.sh | bash
```

将 version 替换为您想要安装的nvm版本号。例如，要安装v0.35.3版本，您可以运行以下命令：

```bash
curl -o- <https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh> | bash
```

安装完成后，重新加载bash shell配置文件，以使nvm命令生效：

```bash
source ~/.bashrc
```

现在，您可以使用nvm命令安装和管理Node.js版本。例如，要安装Node.js v14.16.1版本，可以运行以下命令：

```bash
nvm install v14.16.1
```

安装完成后，您可以使用以下命令查看已安装的Node.js版本：

```bash
nvm ls
```

您可以使用以下命令切换使用的Node.js版本：

```bash
nvm use v14.16.1
```

如果您想将某个版本设置为默认版本，请使用以下命令：

```bash
nvm alias default v14.16.1
```

这样，每次打开新的终端窗口时，nvm都会自动使用默认的Node.js版本。

## 在macOS上安装NVM

第一步:删除现有Node版本
如果你的系统已经安装了node，请先卸载它。我的系统已经通过Homebrew安装了node。所以先把它卸载了。如果还没有安装就跳过。

```shell
brew uninstall --ignore-dependencies node
brew uninstall --force node
```

node npm 卸载

```shell
sudo npm uninstall npm -g

sudo rm -rf /usr/local/lib/node /usr/local/lib/node_modules /var/db/receipts/org.nodejs.*

sudo rm -rf /usr/local/include/node /Users/$USER/.npm

sudo rm /usr/local/bin/node

sudo rm /usr/local/share/man/man1/node.1

sudo rm /usr/local/lib/dtrace/node.d
```

### 在Mac上安装NVM

brew install

```shell
brew update 
brew install nvm
```

接下来，在home目录中为NVM创建一个文件夹。

```shell
mkdir ~/.nvm 
```

现在，配置所需的环境变量。在你的home中编辑以下配置文件

```shell
vim ~/.bash_profile 
```

然后，在 ~/.bash_profile（或~/.zshrc，用于macOS Catalina或更高版本）中添加以下几行

```shell
export NVM_DIR=~/.nvm
source $(brew --prefix nvm)/nvm.sh
```

按ESC + :wq 保存并关闭你的文件。
接下来，将该变量加载到当前的shell环境中。在下一次登录，它将自动加载。

```shell
source ~/.bash_profile
```

NVM已经安装在你的macOS系统上。
下一步，在nvm的帮助下安装你需要的Node.js版本即可。

## 用NVM安装Node.js

首先，看看有哪些Node版本可以安装。要查看可用的版本，请输入。

```bash
nvm ls-remote 
```

现在，你可以安装上述输出中列出的任何版本。你也可以使用别名，如node代表最新版本，lts代表最新的LTS版本，等等。

```bash
nvm install node     ## 安装最后一个长期支持版本
nvm install 10
```

安装后，你可以用以下方法来验证所安装的node.js是否安装成功。

```bash
nvm ls 
```

如果你在系统上安装了多个版本，你可以在任何时候将任何版本设置为默认版本。要设置节点10为默认版本，只需使用。

```bash
nvm use 10
```

同样地，你可以安装其他版本，如Node其他版本，并在它们之间进行切换。
