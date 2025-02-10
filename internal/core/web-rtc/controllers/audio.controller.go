package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
)

type SDP struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

func AudioOfferHandler(c *gin.Context) {
	var offer SDP
	if err := c.BindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SDP offer"})
		return
	}

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating PeerConnection: " + err.Error()})
		return
	}

	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Receiving track: ID=%s, tipo=%s", remoteTrack.ID(), remoteTrack.Kind().String())

		if remoteTrack.Kind() == webrtc.RTPCodecTypeAudio {
			localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, "audio", "pion")
			if err != nil {
				log.Println("Error creating local track:", err)
				return
			}

			sender, err := peerConnection.AddTrack(localTrack)
			if err != nil {
				log.Println("Error adding local track:", err)
				return
			}

			go func() {
				rtcpBuf := make([]byte, 1500)
				for {
					packet, _, err := remoteTrack.ReadRTP()
					if err != nil {
						log.Println("Error reading RTP:", err)
						return
					}

					if err = localTrack.WriteRTP(packet); err != nil {
						log.Println("Error writing RTP:", err)
						return
					}

					if _, _, err := sender.Read(rtcpBuf); err != nil { // todo: this is ugly
						return
					}
				}
			}()
		}
	})

	sdpOffer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer.SDP,
	}
	if err := peerConnection.SetRemoteDescription(sdpOffer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error defining remote description: " + err.Error()})
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating answer: " + err.Error()})
		return
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error describing local description: " + err.Error()})
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gatherComplete

	answerJSON, err := json.Marshal(peerConnection.LocalDescription())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error serializing answer: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", answerJSON)
}
