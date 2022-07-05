package api

import (
	pb "common/module/proto/message_service"
	notificationPb "common/module/proto/notification_service"
	postEvents "common/module/saga/post_events"
	events "common/module/saga/user_events"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"message/module/domain/model"
	"time"
)

func MapNewUser(command *events.UserCommand) *model.User {
	user := &model.User{
		Id:       primitive.NewObjectID(),
		UserId:   command.User.UserId,
		Username: command.User.Username,
		Email:    command.User.Email,
		Settings: model.NotificationSettings{
			Posts:       true,
			Messages:    true,
			Connections: true,
		},
	}
	return user
}

func MapNewPostNotification(command *postEvents.PostNotificationCommand) *model.Notification {

	notification := &model.Notification{
		Id:               primitive.NewObjectID(),
		Timestamp:        time.Now(),
		Content:          command.Notification.Content,
		RedirectPath:     command.Notification.RedirectPath,
		Read:             false,
		Type:             model.POST,
		NotificationFrom: command.Notification.NotificationFrom,
		NotificationTo:   command.Notification.NotificationTo,
	}

	return notification
}

func MapPostNotificationReply(replyType postEvents.PostNotificationReplyType) (reply *postEvents.PostNotificationReply) {
	reply = &postEvents.PostNotificationReply{
		Type: replyType,
	}
	return reply
}

func MapUserReply(user *model.User, replyType events.UserReplyType) (reply *events.UserReply) {
	reply = &events.UserReply{
		Type: replyType,
		PostUser: events.PostUser{
			Id:       user.Id,
			UserId:   user.UserId,
			Username: user.Username,
			Email:    user.Email,
		},
	}
	return reply
}

func MapMessageReply(message *model.Message, receiver string, sender string) (reply *pb.Message) {
	reply = &pb.Message{
		SenderUsername:   sender,
		ReceiverUsername: receiver,
		MessageText:      message.MessageText,
		TimeSent:         message.TimeSent.String(),
	}
	return reply
}
func MapNewMessage(messageText string, receiverId string, senderId string) *model.Message {
	message := &model.Message{
		Id:          primitive.ObjectID{},
		SenderId:    senderId,
		ReceiverId:  receiverId,
		MessageText: messageText,
		TimeSent:    time.Now(),
	}
	return message
}

func MapNotificationResponse(notification *model.Notification) *notificationPb.Notification {
	id := notification.Id.Hex()

	notificationPb := &notificationPb.Notification{
		Id:               id,
		Content:          notification.Content,
		From:             notification.NotificationFrom,
		To:               notification.NotificationTo,
		RedirectPath:     notification.RedirectPath,
		NotificationType: mapNotificationTypeToString(notification.Type),
	}

	return notificationPb
}

func MapSettingsResponse(settings *model.NotificationSettings) *notificationPb.NotificationSettings {

	settingsPb := &notificationPb.NotificationSettings{
		Posts:       settings.Posts,
		Messages:    settings.Messages,
		Connections: settings.Connections,
	}

	return settingsPb
}

func MapChangeSettingsRequest(request *notificationPb.ChangeSettingsRequest) *model.NotificationSettings {

	settingsModel := &model.NotificationSettings{
		Posts:       request.Settings.Posts,
		Messages:    request.Settings.Messages,
		Connections: request.Settings.Connections,
	}

	return settingsModel
}

func mapNotificationTypeToString(notificationType model.NotificationType) string {
	if notificationType == model.POST {
		return "POST"
	}
	if notificationType == model.MESSAGE {
		return "MESSAGE"
	} else {
		return "PROFILE"
	}

}
