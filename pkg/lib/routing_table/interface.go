package routing_table

type ID interface {
	Bytes() []byte
	Equal(ID) bool
}

type RoutingTable interface {
	LocalID() ID
	TableSize() int
	BucketSize(cpl int) int
	CommonPrefixLen(x ID) int
	GenRandIDWithCpl(cpl int) ID

	AddID(x ID) error
	RemoveID(x ID) error

	CloserIDs(x ID, count int, excludes []ID) (ys []ID, found bool)
	ClosestIDs(x ID, count int) (ys []ID, found bool)
}
