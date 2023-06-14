# 要在CentOS上安装nvm（Node Version Manager）

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
