package tx

type protocol string

const (
	UMP  protocol = "UMP"
	DMP  protocol = "DMP"
	HRMP protocol = "HRMP"
	XCMP protocol = "XCMP"
)

func NewMessage(protocol protocol) IXCMP {
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
	panic("invalid protocol")
}
