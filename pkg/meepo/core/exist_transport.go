package meepo_core

func (mp *Meepo) existTransport(target Addr) bool {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()
	return mp.existTransportNL(target)
}

func (mp *Meepo) existTransportNL(target Addr) bool {
	_, found := mp.transports[target]
	return found
}
