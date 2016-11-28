package main

import (
	"encoding/json"
	"fmt"

	webrtc "github.com/keroserene/go-webrtc"
)

// Based on code from: https://github.com/keroserene/go-webrtc/blob/master/demo/chat/chat.go

var defaultIceServer = "stun:stun.l.google.com:19302"

// SignalHandler is called with a signal that needs to be relayed to the peer out of band.
type SignalHandler func(signal string)
type ChannelHandler func(channel *webrtc.DataChannel)

func Peer(onSignal SignalHandler, onChannel ChannelHandler) (*peer, error) {
	p := &peer{
		OnSignal:  onSignal,
		OnChannel: onChannel,
	}
	err := p.start()
	if err != nil {
		return nil, err
	}
	return p, nil
}

type peer struct {
	pc *webrtc.PeerConnection

	OnSignal  SignalHandler
	OnChannel ChannelHandler
}

func (p *peer) generateOffer() error {
	offer, err := p.pc.CreateOffer() // blocking
	if err != nil {
		return err
	}
	p.pc.SetLocalDescription(offer)
	fmt.Println("XXX: generated offer")
	return nil
}

func (p *peer) generateAnswer() error {
	answer, err := p.pc.CreateAnswer() // blocking
	if err != nil {
		return err
	}
	p.pc.SetLocalDescription(answer)
	fmt.Println("XXX: generated answer")
	return nil
}

func (p *peer) start() error {
	config := webrtc.NewConfiguration(
		webrtc.OptionIceServer(defaultIceServer),
	)

	pc, err := webrtc.NewPeerConnection(config)
	if nil != err {
		return err
	}
	p.pc = pc

	// OnNegotiationNeeded is triggered when something important has occurred in
	// the state of PeerConnection (such as creating a new data channel), in which
	// case a new SDP offer must be prepared and sent to the remote peer.
	pc.OnNegotiationNeeded = func() {
		go p.generateOffer()
	}

	// Once all ICE candidates are prepared, they need to be sent to the remote
	// peer which will attempt reaching the local peer through NATs.
	pc.OnIceComplete = func() {
		// Finished gathering ICE candidates
		sdp := pc.LocalDescription().Serialize()
		go p.OnSignal(sdp)
	}

	// Called when a peer initiates a data channel
	pc.OnDataChannel = p.OnChannel
	return nil
}

func (p *peer) CreateDataChannel(label string) error {
	// Attempting to create the first datachannel triggers ICE.
	dc, err := p.pc.CreateDataChannel(label, webrtc.Init{})
	if nil != err {
		return err
	}

	go p.OnChannel(dc)
	return nil
}

func (p *peer) Connect(signal string) error {
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(signal), &parsed)
	if nil != err {
		return err
	}

	if nil != parsed["sdp"] {
		sdp := webrtc.DeserializeSessionDescription(signal)
		if nil == sdp {
			return fmt.Errorf("invalid sdp: %s", signal)
		}

		err = p.pc.SetRemoteDescription(sdp)
		if nil != err {
			return err
		}
		fmt.Println("SDP " + sdp.Type + " successfully received.")
		if sdp.Type == "offer" {
			go p.generateAnswer()
		}
	}

	// Allow individual ICE candidate messages, but this won't be necessary if
	// the remote peer also doesn't use trickle ICE.
	if nil != parsed["candidate"] {
		ice := webrtc.DeserializeIceCandidate(signal)
		if nil == ice {
			return fmt.Errorf("invalid ICE candidate: %s", signal)
		}
		p.pc.AddIceCandidate(*ice)
		fmt.Println("ICE candidate successfully received.")
	}
	return nil
}
