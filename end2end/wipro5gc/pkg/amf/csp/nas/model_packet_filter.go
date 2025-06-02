package nas

// type PacketFilterList interface{}

// type PacketFilterModDel struct {
// 	Identifier uint8
// }

type PacketFilter struct {
	Identifier uint8
	Direction  string
	Components []string
}
