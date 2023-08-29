package meepo_core

func (mp *Meepo) onTeleportationClose(tp Teleportation) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	delete(mp.teleportations, tp.ID())
}
