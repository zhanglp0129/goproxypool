package constant

type ProxyAddressStatus byte

const (
	// PendingDetection 待检测
	PendingDetection ProxyAddressStatus = iota + 1
	// AcceptDetection 通过检测
	AcceptDetection
	// FailedDetection 未通过检测
	FailedDetection
)
