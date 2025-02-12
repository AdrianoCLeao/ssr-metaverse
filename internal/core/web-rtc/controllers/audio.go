package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
	"ssr-metaverse/internal/core/error"
)

// SDP representa a estrutura do SDP recebido via JSON.
type SDP struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

// Client representa cada cliente conectado ao SFU.
type Client struct {
	id             string
	pc             *webrtc.PeerConnection
	// subscriptions mapeia (por exemplo) o id do publicador para a TrackLocal que receberá
	// os pacotes RTP daquele publicador.
	subscriptions  map[string]*webrtc.TrackLocalStaticRTP
	// published indica se o cliente já publicou (enviou) um track de áudio.
	published      bool
	// publisherCodec guarda o codec do track de áudio publicado (necessário para criação
	// das tracks de assinatura).
	publisherCodec *webrtc.RTPCodecCapability
}

// Variáveis globais para controlar os clientes conectados.
var (
	sfuClients    = make(map[string]*Client)
	sfuClientsMu  sync.Mutex
	clientCounter int
)

// AudioOfferHandler lida com o SDP offer de um cliente e configura o PeerConnection.
// Além disso, implementa a lógica básica de SFU para áudio.
func AudioOfferHandler(c *gin.Context) {
	var offer SDP

	if err := c.BindJSON(&offer); err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid SDP offer: " + err.Error(),
		})
		return
	}

	// Configuração do STUN (para ICE)
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error creating PeerConnection: " + err.Error(),
		})
		return
	}

	// Cria um identificador único para o cliente.
	sfuClientsMu.Lock()
	clientCounter++
	clientID := fmt.Sprintf("client-%d", clientCounter)
	client := &Client{
		id:            clientID,
		pc:            peerConnection,
		subscriptions: make(map[string]*webrtc.TrackLocalStaticRTP),
	}
	// Adiciona o novo cliente à lista global.
	sfuClients[clientID] = client

	// Se já existirem clientes que publicaram áudio, adiciona tracks de assinatura
	// para que o novo cliente receba o áudio deles.
	for _, other := range sfuClients {
		if other.id == clientID {
			continue
		}
		if other.published && other.publisherCodec != nil {
			newTrack, err := webrtc.NewTrackLocalStaticRTP(*other.publisherCodec, "audio", other.id)
			if err != nil {
				log.Printf("Error creating subscription track for publisher %s: %v", other.id, err)
				continue
			}
			if _, err := client.pc.AddTrack(newTrack); err != nil {
				log.Printf("Error adding subscription track for publisher %s: %v", other.id, err)
				continue
			}
			client.subscriptions[other.id] = newTrack
		}
	}
	sfuClientsMu.Unlock()

	// Quando um track for recebido (isto é, quando o cliente enviar áudio),
	// este callback é chamado.
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Client %s received track: ID=%s, Kind=%s", clientID, remoteTrack.ID(), remoteTrack.Kind().String())
		if remoteTrack.Kind() == webrtc.RTPCodecTypeAudio {
			// O cliente está publicando áudio.
			sfuClientsMu.Lock()
			client.published = true
			codecCap := remoteTrack.Codec().RTPCodecCapability
			client.publisherCodec = &codecCap
			// Para cada outro cliente já conectado, cria uma track de assinatura para receber
			// o áudio deste publicador.
			for _, other := range sfuClients {
				if other.id == clientID {
					continue
				}
				if _, exists := other.subscriptions[clientID]; !exists {
					newTrack, err := webrtc.NewTrackLocalStaticRTP(codecCap, "audio", clientID)
					if err != nil {
						log.Printf("Error creating subscription track for client %s: %v", other.id, err)
						continue
					}
					if _, err = other.pc.AddTrack(newTrack); err != nil {
						log.Printf("Error adding subscription track for client %s: %v", other.id, err)
						continue
					}
					other.subscriptions[clientID] = newTrack
				}
			}
			sfuClientsMu.Unlock()

			// Inicia a leitura dos pacotes RTP e encaminha para todos os assinantes (clientes
			// diferentes do publicador).
			rtcpBuf := make([]byte, 1500)
			for {
				packet, _, err := remoteTrack.ReadRTP()
				if err != nil {
					log.Printf("Error reading RTP from client %s: %v", clientID, err)
					break
				}

				// Para cada cliente conectado (exceto o próprio publicador) encaminha o pacote RTP.
				sfuClientsMu.Lock()
				for _, other := range sfuClients {
					if other.id == clientID {
						continue
					}
					if track, ok := other.subscriptions[clientID]; ok {
						if err := track.WriteRTP(packet); err != nil {
							log.Printf("Error writing RTP to client %s: %v", other.id, err)
						}
					}
				}
				sfuClientsMu.Unlock()

				// Opcionalmente, lê pacotes RTCP do receptor (para feedback, etc).
				if _, _, err := receiver.Read(rtcpBuf); err != nil {
					break
				}
			}
		}
	})

	// Define o SDP remoto (oferta recebida).
	sdpOffer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offer.SDP,
	}
	if err := peerConnection.SetRemoteDescription(sdpOffer); err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error setting remote description: " + err.Error(),
		})
		return
	}

	// Cria e define a descrição local (a resposta).
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error creating answer: " + err.Error(),
		})
		return
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error setting local description: " + err.Error(),
		})
		return
	}

	// Aguarda a conclusão da gathering dos ICE candidates.
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	<-gatherComplete

	answerJSON, err := json.Marshal(peerConnection.LocalDescription())
	if err != nil {
		error.RespondWithError(c, error.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error serializing answer: " + err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "application/json", answerJSON)
}
