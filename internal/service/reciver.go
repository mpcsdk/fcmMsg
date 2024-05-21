// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

type (
	IReceiver interface{}
)

var (
	localReceiver IReceiver
)

func Receiver() IReceiver {
	if localReceiver == nil {
		panic("implement not found for interface IReceiver, forgot register?")
	}
	return localReceiver
}

func RegisterReceiver(i IReceiver) {
	localReceiver = i
}
