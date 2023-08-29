package meepo_core

func (mp *Meepo) onTeleportationNew(tp Teleportation) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	mp.teleportations[tp.ID()] = tp
}
