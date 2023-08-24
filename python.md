# python

## macOS 安装和管理多个Python版本

安装步骤

安装homebrew：🚀

```shell
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
$ brew -v
Homebrew 1.6.9
Homebrew/homebrew-core (git revision 5707e; last commit 2018-07-09)
```

安装pyenv：🛰

``` bash
brew update
brew install pyenv
pyenv -v

pyenv 1.2.5
```

安装管理多个Python：

```shell
$ pyenv install 2.7.15
$ pyenv install 3.7.0
$ pyenv versions
  system
  2.7.15
* 3.7.0 (set by /Users/john/.pyenv/version)
```

注：星号指定当前的版本

切换版本：

```shell
$ pyenv global 2.7.15
$ pyenv versions
  system
* 2.7.15 (set by /Users/john/.pyenv/version)
  3.7.0
$ python --version
Python 2.7.15
```

pyenv常用的命令说明：

```text
使用方式: pyenv <命令> [<参数>]

命令:
  commands    查看所有命令
  local       设置或显示本地的Python版本
  global      设置或显示全局Python版本
  shell       设置或显示shell指定的Python版本
  install     安装指定Python版本
  uninstall   卸载指定Python版本)
  version     显示当前的Python版本及其本地路径
  versions    查看所有已经安装的版本
  which       显示安装路径

```

注：使用local、global、shell，设置Python版本时需要跟上参数（版本号），查看则不需要。
