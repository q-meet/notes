# Telegram

## Go Telegram 机器人 API

### 前言

Telegram机器人的发送信息功能需要指定聊天窗口ID，接下来介绍的则是如何利用机器人获取Telegram的个人ID及群组ID。

获取个人ID
可以使用@userinfobot进行获取个人ID
也可以给自己的机器人发送信息，然后打开

    https://api.telegram.org/bot(机器人TOKEN)/getUpdates
查找自己发的消息的id字段

获取群组ID
需要将机器人拉入群组（只有管理员或者创始人才有权限）
然后发送/xxx @机器人
最后打开

    https://api.telegram.org/bot(机器人TOKEN)/getUpdates
获取发送信息的id字段（群组ID为-开头）

### 入门

该库被设计为 Telegram Bot API 的简单包装器。我们鼓励您首先阅读Telegram 的文档，以了解机器人的能力。他们还提供了一些解决常见问题的好方法。

安装中

go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5

一个简单的机器人
为了演练基础知识，让我们创建一个简单的回显机器人，它会重复您所说的话来回复您的消息。在继续之前，请确保您从@Botfather获取 API 令牌 。

让我们从构建一个新的BotAPI开始。

```go
package main

import (
    "os"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
    if err != nil {
        panic(err)
    }

    bot.Debug = true
}
```

我们没有将 API 令牌直接输入到文件中，而是使用环境变量。这使得我们可以轻松配置我们的机器人以使用正确的帐户，并防止我们将真实的代币泄露给世界。拥有您的令牌的任何人都可以从您的机器人发送和接收消息！

我们还进行了设置bot.Debug = true，以便获取有关发送到 Telegram 的请求的更多信息。如果运行上面的示例，您将看到有关端点请求的信息getMe。库会自动调用此函数以确保您的令牌按预期工作。它还使用有关机器人的信息填充结构Self中的字段。BotAPI

现在我们已经连接到 Telegram，让我们开始获取更新并执行操作。我们可以在启用调试模式的行之后添加此代码。

```go
    // Create a new UpdateConfig struct with an offset of 0. Offsets are used
    // to make sure Telegram knows we've handled previous values and we don't
    // need them repeated.
    updateConfig := tgbotapi.NewUpdate(0)

    // Tell Telegram we should wait up to 30 seconds on each request for an
    // update. This way we can get information just as quickly as making many
    // frequent requests without having to send nearly as many.
    updateConfig.Timeout = 30

    // Start polling Telegram for updates.
    updates := bot.GetUpdatesChan(updateConfig)

    // Let's go through each update that we're getting from Telegram.
    for update := range updates {
        // Telegram can send many types of updates depending on what your Bot
        // is up to. We only want to look at messages for now, so we can
        // discard any other updates.
        if update.Message == nil {
            continue
        }

        // Now that we know we've gotten a new message, we can construct a
        // reply! We'll take the Chat ID and Text from the incoming message
        // and use it to create a new message.
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
        // We'll also say that this message is a reply to the previous message.
        // For any other specifications than Chat ID or Text, you'll need to
        // set fields on the `MessageConfig`.
        msg.ReplyToMessageID = update.Message.MessageID

        // Okay, we're sending our message off! We don't care about the message
        // we just sent, so we'll discard it.
        if _, err := bot.Send(msg); err != nil {
            // Note that panics are a bad way to handle errors. Telegram can
            // have service outages or network errors, you should retry sending
            // messages or more gracefully handle failures.
            panic(err)
        }
    }
```

恭喜！您已经创建了自己的机器人！

现在您已经掌握了一些基础知识，我们可以开始讨论该库的结构和更高级的功能。

### BotAPI Test

bot_test.go 文件文件中
配置以下信息 获取方法如上前言，即可调用测试方法使用

```text
TestToken
ChatID
Channel
SupergroupChatID
ReplyToMessageID
```
