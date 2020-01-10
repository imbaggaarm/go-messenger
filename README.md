# go-messenger [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GitHub version](https://badge.fury.io/gh/imbaggaarm%2Fgo-messenger.svg)](https://badge.fury.io/gh/imbaggaarm%2Fgo-messenger)
A [Facebook Messenger API](https://developers.facebook.com/docs/messenger-platform) wrapper written in **Golang**,
alternative to pymessenger.

## Features
- [x] [Send raw message](https://developers.facebook.com/docs/messenger-platform/reference/send-api/) - SendRawMessage(payload)
- [x] [Send action](https://developers.facebook.com/docs/messenger-platform/send-api-reference/sender-actions) - SendAction(recipientId, action, notificationType)
- [x] [Send message](https://developers.facebook.com/docs/messenger-platform/send-messages) - SendMessage(recipientId, message)
- [x] [Send text message](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_text) - SendTextMessage(recipientId, text)
- [x] [Send quick replies](https://developers.facebook.com/docs/messenger-platform/send-messages/quick-replies) - SendQuickReplies(recipientId, text, quickReplies)
- [x] [Send attachment message](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments) - SendAttachmentMessage(recipientId, attachment)
- [x] [Send attachment with url](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments) - SendAttachmentUrl(recipientId, attachmentType)
- [x] [Send generic message](https://developers.facebook.com/docs/messenger-platform/reference/template/generic) - SendGenericMessage(recipientId, elements)
- [x] [Send button message](https://developers.facebook.com/docs/messenger-platform/send-messages/buttons) - SendButtonMessage(recipientId, text, buttons)
- [x] [Send image with url](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments) - SendImageUrl(recipientId, imageUrl)
- [x] [Send audio with url](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments) - SendAudioUrl(recipientId, audioUrl)
- [x] [Send video with url](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments) - SendVideoUrl(recipientId, videoUrl)
- [x] [Send file with url](https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments) - SendFileUrl(recipientId, fileUrl)
- [x] [Set get started button](https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/get-started-button) - SetGetStarted(gsPayload)
- [x] [Remove get started button](https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/#delete) - RemoveGetStarted()
- [x] [Set persistent menu](https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/persistent-menu) - SetPersistentMenu(pmPayload)
- [x] [Remove persistent menu](https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/#delete) - RemovePersistentMenu()
## Getting Started
### Installation
```
go get -u github.com/imbaggaarm/go-messenger
```
### Example
```Go
import messenger "github.com/imbaggaarm/go-messenger"

apiVersion :=  messenger.DefaultApiVersion // or the version that you want
bot = messenger.NewBot(accessToken, apiVersion)

textMessage := "Hello! Can you hear me?"
bot.sendTextMessage(recipientId, textMessage)
```
## Usage
- [VNUChatbot](https://www.facebook.com/vnuchat/) is a chat-with-stranger chatbot, with back-end written in Golang and use this package. 
## Other
### Future of go-messenger
There are a lot of missing functions in this package, 
I'm planning to make this better and better in the future. Here are some things will be implemented soon:
- An example (chat with strangers)
- Attachment with file messages
### Contact
Follow and contact me on [Twitter](http://twitter.com/baggaarm). If you find an issue, just [open a ticket](https://github.com/imbaggaarm/go-messenger/issues/new). 
Pull requests are warmly welcome as well.