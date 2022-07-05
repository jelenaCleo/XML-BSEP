package orchestrators

import (
	saga "common/module/saga/messaging"
	events "common/module/saga/post_events"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewPostOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*PostOrchestrator, error) {
	o := &PostOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *PostOrchestrator) LikePost(postId primitive.ObjectID, liker string, postOwner string) error {
	fmt.Println("evo me u orchestratoru u create noty ")

	events := &events.PostNotificationCommand{
		Type: events.LikePost,
		Notification: events.Notification{
			Content:          liker + " liked your post.",
			RedirectPath:     "/post/" + postId.Hex(),
			NotificationFrom: liker,
			NotificationTo:   postOwner,
		},
	}

	return o.commandPublisher.Publish(events)
}

func (o *PostOrchestrator) handle(reply *events.PostNotificationReply) events.PostNotificationReplyType {
	if reply.Type == events.NotificationSent {
		fmt.Println("Senttttttttt")
	}
	if reply.Type == events.UnknownReply {
		fmt.Println("Unknown ")
	}
	fmt.Println("upsssss")
	return events.UnknownReply

}
