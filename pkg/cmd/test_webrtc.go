package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/pkg/lib/stun"
)

var (
	testWebrtcCmd = &cobra.Command{
		Use:  "webrtc",
		RunE: meepoTestWebrtc,
		Args: cobra.NoArgs,
	}

	testWebrtcOptions struct {
		Offerer  bool
		Answerer bool
	}
)

func meepoTestWebrtc(cmd *cobra.Command, args []string) error {
	if !(testWebrtcOptions.Offerer || testWebrtcOptions.Answerer) {
		return fmt.Errorf("either offerer or answerer should be set to true")
	}

	var se webrtc.SettingEngine
	api := webrtc.NewAPI(webrtc.WithSettingEngine(se))
	pc, err := api.NewPeerConnection(webrtc.Configuration{ICEServers: []webrtc.ICEServer{{URLs: stun.STUNS}}})
	if err != nil {
		return err
	}

	if testWebrtcOptions.Offerer {
		pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
			fmt.Printf("offerer: peer connection state: %s\n", s.String())
		})

		dc, err := pc.CreateDataChannel("data", nil)
		if err != nil {
			return err
		}
		dc.OnOpen(func() {
			dc.OnMessage(func(m webrtc.DataChannelMessage) {
				fmt.Printf("offerer: recvMsg: %s\n", string(m.Data))
			})
		})

		offer, err := pc.CreateOffer(nil)
		if err != nil {
			return err
		}
		if err = pc.SetLocalDescription(offer); err != nil {
			return err
		}

		<-webrtc.GatheringCompletePromise(pc)

		offerStr, err := encodeSessionDescription(*pc.LocalDescription())
		if err != nil {
			return err
		}

		fmt.Println("copy offer to answerer:")
		fmt.Println(offerStr)
		fmt.Println()

		fmt.Println("paste answer from answerer:")
		answerStr, err := readFromStdin()
		if err != nil {
			return err
		}

		var answer webrtc.SessionDescription
		if err = decodeSessionDescription(answerStr, &answer); err != nil {
			return err
		}

		if err = pc.SetRemoteDescription(answer); err != nil {
			return err
		}
	} else {
		pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
			fmt.Printf("answerer: peer connection state: %s\n", s.String())
		})

		pc.OnDataChannel(func(dc *webrtc.DataChannel) {
			dc.OnOpen(func() {
				for t := range time.NewTicker(5 * time.Second).C {
					msg := t.String()
					if err := dc.SendText(msg); err != nil {
						fmt.Printf("answerer: error: %s\n", err.Error())
						return
					}
					fmt.Printf("answerer: sendMsg: %s\n", msg)
				}
			})
		})

		var offer webrtc.SessionDescription
		fmt.Println("paste offer from offerer:")
		offerStr, err := readFromStdin()
		if err != nil {
			return err
		}
		if err = decodeSessionDescription(offerStr, &offer); err != nil {
			return err
		}
		if err = pc.SetRemoteDescription(offer); err != nil {
			return err
		}

		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			return err
		}

		if err = pc.SetLocalDescription(answer); err != nil {
			return err
		}

		<-webrtc.GatheringCompletePromise(pc)

		answerStr, err := encodeSessionDescription(*pc.LocalDescription())
		if err != nil {
			return err
		}

		fmt.Println("copy answer to offerer:")
		fmt.Println(answerStr)
		fmt.Println()
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)
	<-c
	return nil
}

func init() {
	fs := testWebrtcCmd.Flags()

	fs.BoolVar(&testWebrtcOptions.Offerer, "offerer", false, "as offerer")
	fs.BoolVar(&testWebrtcOptions.Answerer, "answerer", false, "as answerer")

	testCmd.AddCommand(testWebrtcCmd)
}
