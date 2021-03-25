package meepo

import "github.com/PeerXu/meepo/pkg/teleportation"

func (mp *Meepo) GetTeleportation(name string) (teleportation.Teleportation, error) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	return mp.getTeleportationNL(name)
}

func (mp *Meepo) getTeleportationNL(name string) (teleportation.Teleportation, error) {
	var tp teleportation.Teleportation
	var ok bool

	tp, ok = mp.getTeleportationSourceNL(name)
	if !ok {
		tp, ok = mp.getTeleportationSinkNL(name)
		if !ok {
			return nil, TeleportationNotExistError
		}
	}

	return tp, nil
}

func (mp *Meepo) addTeleportationSource(name string, ts *teleportation.TeleportationSource) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	mp.addTeleportationSourceNL(name, ts)
}

func (mp *Meepo) addTeleportationSourceNL(name string, ts *teleportation.TeleportationSource) {
	mp.teleportationSources[name] = ts
}

func (mp *Meepo) removeTeleportationSource(name string) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	mp.removeTeleportationSourceNL(name)
}

func (mp *Meepo) removeTeleportationSourceNL(name string) {
	delete(mp.teleportationSources, name)
}

func (mp *Meepo) addTeleportationSink(name string, ts *teleportation.TeleportationSink) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	mp.addTeleportationSinkNL(name, ts)
}

func (mp *Meepo) addTeleportationSinkNL(name string, ts *teleportation.TeleportationSink) {
	mp.teleportationSinks[name] = ts
}

func (mp *Meepo) removeTeleportationSink(name string) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	mp.removeTeleportationSinkNL(name)
}

func (mp *Meepo) removeTeleportationSinkNL(name string) {
	delete(mp.teleportationSinks, name)
}

func (mp *Meepo) listTeleportations() ([]teleportation.Teleportation, error) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	return mp.listTeleportationsNL()
}

func (mp *Meepo) listTeleportationsNL() ([]teleportation.Teleportation, error) {
	var teleportations []teleportation.Teleportation

	for _, ts := range mp.teleportationSources {
		teleportations = append(teleportations, ts)
	}

	for _, ts := range mp.teleportationSinks {
		teleportations = append(teleportations, ts)
	}

	return teleportations, nil
}

func (mp *Meepo) listTeleportationsByPeerID(id string) ([]teleportation.Teleportation, error) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()

	return mp.listTeleportationsByPeerIDNL(id)
}

func (mp *Meepo) listTeleportationsByPeerIDNL(id string) ([]teleportation.Teleportation, error) {
	xs, err := mp.listTeleportationsNL()
	if err != nil {
		return nil, err
	}

	var ys []teleportation.Teleportation
	for _, x := range xs {
		if x.Transport().PeerID() == id {
			ys = append(ys, x)
		}
	}

	return ys, nil
}

func (mp *Meepo) getTeleportationSource(name string) (*teleportation.TeleportationSource, bool) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	return mp.getTeleportationSourceNL(name)
}

func (mp *Meepo) getTeleportationSourceNL(name string) (*teleportation.TeleportationSource, bool) {
	ts, ok := mp.teleportationSources[name]
	return ts, ok
}

func (mp *Meepo) getTeleportationSink(name string) (*teleportation.TeleportationSink, bool) {
	mp.teleportationsMtx.Lock()
	defer mp.teleportationsMtx.Unlock()
	return mp.getTeleportationSinkNL(name)
}

func (mp *Meepo) getTeleportationSinkNL(name string) (*teleportation.TeleportationSink, bool) {
	ts, ok := mp.teleportationSinks[name]
	return ts, ok
}
