package tx

type Protocol string

const (
	UMP  Protocol = "UMP"
	DMP  Protocol = "DMP"
	HRMP Protocol = "HRMP"
	XCMP Protocol = "XCMP"
)

func NewMessage(protocol Protocol) IXCMP {
	switch protocol {
	case UMP:
		return NewUmp()
	case DMP:
		return NewDmp()
	case HRMP:
		return NewHrmp()
	case XCMP:
		return NewHrmp()
	}
	panic("invalid Protocol")
}
