package gregor

import (
	"net"
	"time"

	context "golang.org/x/net/context"
)

type InBandMsgType int

const (
	InBandMsgTypeNone   InBandMsgType = 0
	InBandMsgTypeUpdate InBandMsgType = 1
	InBandMsgTypeSync   InBandMsgType = 2
)

type UID interface {
	Bytes() []byte
}

type MsgID interface {
	Bytes() []byte
}

type DeviceID interface {
	Bytes() []byte
}

type System interface {
	String() string
}

type Category interface {
	String() string
}

type Body interface {
	Bytes() []byte
}

type Metadata interface {
	UID() UID
	MsgID() MsgID
	CTime() time.Time
	SetCTime(time.Time)
	DeviceID() DeviceID
	InBandMsgType() InBandMsgType
}

type MessageWithMetadata interface {
	Metadata() Metadata
}

type InBandMessage interface {
	MessageWithMetadata
	ToStateUpdateMessage() StateUpdateMessage
	ToStateSyncMessage() StateSyncMessage
	Merge(m1 InBandMessage) error
}

type StateUpdateMessage interface {
	MessageWithMetadata
	Creation() Item
	Dismissal() Dismissal
}

type StateSyncMessage interface {
	MessageWithMetadata
}

type OutOfBandMessage interface {
	System() System
	UID() UID
	Body() Body
}

type TimeOrOffset interface {
	Time() *time.Time
	Offset() *time.Duration
}

type Item interface {
	MessageWithMetadata
	DTime() TimeOrOffset
	NotifyTimes() []TimeOrOffset
	Body() Body
	Category() Category
}

type MsgRange interface {
	EndTime() TimeOrOffset
	Category() Category
}

type Dismissal interface {
	MsgIDsToDismiss() []MsgID
	RangesToDismiss() []MsgRange
}

type State interface {
	Items() ([]Item, error)
	ItemsInCategory(c Category) ([]Item, error)
}

type Message interface {
	ToInBandMessage() InBandMessage
	ToOutOfBandMessage() OutOfBandMessage
}

// MessageConsumer consumes state update messages. It's half of
// the state machine protocol
type MessageConsumer interface {
	// ConsumeMessage is called on a new incoming message to mutate the state
	// of the state machine. Of course messages can be "inband" which actually
	// perform state mutations, or might be "out-of-band" that just use the
	// Gregor broadcast mechanism to make sure that all clients get the
	// notification.
	ConsumeMessage(m Message) error
}

// StateMachine is the central interface of the Gregor system. Various parts of the
// server and client infrastructure will implement various parts of this interface,
// to ensure that the state machine can be replicated, and that it can be queried.
type StateMachine interface {
	MessageConsumer

	// State returns the state for the user u on device d at time t.
	// d can be nil, in which case the global state (across all devices)
	// is returned. If t is nil, then use Now, otherwise, return the state
	// at the given time.
	State(u UID, d DeviceID, t TimeOrOffset) (State, error)

	// InBandMessagesSince returns all messages since the given time
	// for the user u on device d.  If d is nil, then we'll return
	// all messages across all devices.  If d is a device, then we'll
	// return global messages and per-device messages for that device.
	InBandMessagesSince(u UID, d DeviceID, t TimeOrOffset) ([]InBandMessage, error)
}

type ObjFactory interface {
	MakeUID(b []byte) (UID, error)
	MakeMsgID(b []byte) (MsgID, error)
	MakeDeviceID(b []byte) (DeviceID, error)
	MakeBody(b []byte) (Body, error)
	MakeCategory(s string) (Category, error)
	MakeItem(u UID, msgid MsgID, deviceid DeviceID, ctime time.Time, c Category, dtime *time.Time, body Body) (Item, error)
	MakeDismissalByRange(uid UID, msgid MsgID, devid DeviceID, ctime time.Time, c Category, d time.Time) (InBandMessage, error)
	MakeDismissalByID(uid UID, msgid MsgID, devid DeviceID, ctime time.Time, d MsgID) (InBandMessage, error)
	MakeStateSyncMessage(uid UID, msgid MsgID, devid DeviceID, ctime time.Time) (InBandMessage, error)
	MakeState(i []Item) (State, error)
	MakeMetadata(uid UID, msgid MsgID, devid DeviceID, ctime time.Time, i InBandMsgType) (Metadata, error)
	MakeInBandMessageFromItem(i Item) (InBandMessage, error)
}

type NetworkInterfaceIncoming interface {
	ConsumeMessage(c context.Context, m Message) error
}

type NetworkInterfaceOutgoing interface {
	BroadcastMessage(c context.Context, m Message) error
}

type NetworkInterface interface {
	NetworkInterfaceOutgoing
	Serve(i NetworkInterfaceIncoming) error
}

type MainLoopServer interface {
	Serve(n net.Listener) error
}

func UIDFromMessage(m Message) UID {
	ibm := m.ToInBandMessage()
	if ibm != nil && ibm.Metadata() != nil {
		return ibm.Metadata().UID()
	}
	if oobm := m.ToOutOfBandMessage(); oobm != nil {
		return oobm.UID()
	}
	return nil
}
