package common

type DebugFlags struct {
	// 测试rpc timeout机制
	RPCNotSend bool
    MailBoxStopKeepRunService bool
}

var debugFlags *DebugFlags

func init() {
	debugFlags = &DebugFlags{
		RPCNotSend: false,
        MailBoxStopKeepRunService: false,
	}
}

func GetDebugFlags() *DebugFlags {
	return debugFlags
}
