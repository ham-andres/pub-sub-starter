package pubsub

type SimpleQueueType int 

const (
	Durable SimpleQueueType = iota
	Transient 
)


func setQueueType(qt SimpleQueueType)(bool,bool,bool) {
	var isDurable, autoDelete, exclusive bool 
	if qt == Durable {
		isDurable, autoDelete, exclusive = true,false ,false
	} else {
		isDurable, autoDelete, exclusive = false, true, true
	}
	return isDurable, autoDelete, exclusive

}
