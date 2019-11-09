package notify

// Javascript Binding
//
// IC.Notifr()
//	 .newSMS()
//	 .newEmail()
//	 .newSlack()
//	 .send(msg)
//
type NotifyBind struct {
	notify *Notify
}

func (b *NotifyBind) NewSMS() *SMS {
	return NewSMS()
}

func (b *NotifyBind) NewEmail() *Email {
	return NewEmail()
}

func (b *NotifyBind) NewSlack() *SlackMessage {
	return NewSlackMessage()
}

func (b *NotifyBind) Send(msg interface{}) {
	b.notify.Send(msg)
}
