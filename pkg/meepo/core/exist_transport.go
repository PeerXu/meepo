package meepo_core

func (mp *Meepo) existTransport(target Addr) bool {
	mp.transportsMtx.Lock()
	defer mp.transportsMtx.Unlock()
	_, found := mp.transports[target]
	return found
}
