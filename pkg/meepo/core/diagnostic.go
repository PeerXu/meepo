package meepo_core

import (
	"github.com/pion/webrtc/v3"

	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
)

func (mp *Meepo) Diagnostic() (meepo_interface.DiagnosticReport, error) {
	logger := mp.GetLogger().WithField("#method", "Diagnostic")

	healthReport := mp.routingTable.HealthReport()

	ts, err := mp.ListTransports(mp.context())
	if err != nil {
		logger.WithError(err).Debugf("failed to list transports")
		return nil, err
	}
	tvs := ViewTransports(ts)

	tps, err := mp.ListTeleportations(mp.context())
	if err != nil {
		logger.WithError(err).Debugf("failed to list teleportations")
		return nil, err
	}
	tpvs := ViewTeleportations(tps)

	tcs, err := mp.ListChannels(mp.context())
	if err != nil {
		logger.WithError(err).Debugf("failed to list channels")
		return nil, err
	}
	var cvs []sdk_interface.ChannelView
	for addr, cs := range tcs {
		cvs = append(cvs, ViewChannelsWithAddr(cs, addr)...)
	}

	pc, err := mp.newPeerConnection()
	if err != nil {
		logger.WithError(err).Debugf("failed to new peer connection")
		return nil, err
	}
	defer pc.Close()

	dc, err := pc.CreateDataChannel("_IGNORE_", nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create data channel")
		return nil, err
	}
	defer dc.Close()

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		logger.WithError(err).Debugf("failed to create offer")
		return nil, err
	}
	gatherComplted := webrtc.GatheringCompletePromise(pc)
	if err = pc.SetLocalDescription(offer); err != nil {
		logger.WithError(err).Debugf("failed to gather")
		return nil, err
	}
	<-gatherComplted

	rp := map[string]any{
		"addr":           mp.Addr().String(),
		"transports":     tvs,
		"teleportations": tpvs,
		"channels":       cvs,
		"webrtc": map[string]any{
			"offer": pc.LocalDescription(),
		},
		"routingTable": map[string]any{
			"dhtAlpha": mp.dhtAlpha,
			"healthReport": map[string]any{
				"report":  healthReport.Report,
				"summary": healthReport.Summary,
			},
		},
		"poof": map[string]any{
			"interval": mp.poofInterval,
			"count":    mp.poofCount,
		},
		"mux": map[string]any{
			"enable":    mp.enableMux,
			"ver":       mp.muxVer,
			"buf":       mp.muxBuf,
			"streamBuf": mp.muxStreamBuf,
			"nocomp":    mp.muxNocomp,
		},
		"kcp": map[string]any{
			"enable":      mp.enableKcp,
			"preset":      mp.kcpPreset,
			"crypt":       mp.kcpCrypt,
			"mtu":         mp.kcpMtu,
			"sndwnd":      mp.kcpSndwnd,
			"rcvwnd":      mp.kcpRcvwnd,
			"dataShard":   mp.kcpDataShard,
			"parityShard": mp.kcpParityShard,
		},
	}

	logger.Tracef("diagnostic")

	return rp, nil
}
