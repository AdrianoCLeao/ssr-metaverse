package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

var config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type ThreadSafeWriter struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	return t.Conn.WriteJSON(v)
}

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	Websocket      *ThreadSafeWriter
}

type Peers struct {
	ListLock    sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type Room struct {
	Peers *Peers
}

func NewRoom() *Room {
	return &Room{
		Peers: &Peers{
			Connections: make([]PeerConnectionState, 0),
			TrackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
		},
	}
}

func WebRTCHandler(room *Room) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Erro ao fazer upgrade para WebSocket:", err)
			return
		}
		defer ws.Close()

		peerConnection, err := webrtc.NewPeerConnection(config)
		if err != nil {
			log.Println("Erro ao criar PeerConnection:", err)
			return
		}
		defer peerConnection.Close()

		for _, kind := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeAudio, webrtc.RTPCodecTypeVideo} {
			_, err := peerConnection.AddTransceiverFromKind(kind, webrtc.RTPTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionRecvonly,
			})
			if err != nil {
				log.Println("Erro ao adicionar transceiver:", err)
				return
			}
		}

		writer := &ThreadSafeWriter{Conn: ws}
		newPeer := PeerConnectionState{
			PeerConnection: peerConnection,
			Websocket:      writer,
		}

		room.Peers.AddPeer(newPeer)

		peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
			if candidate == nil {
				return
			}
			candidateJSON, err := json.Marshal(candidate.ToJSON())
			if err != nil {
				log.Println("Erro ao serializar ICE candidate:", err)
				return
			}
			msg := websocketMessage{
				Event: "candidate",
				Data:  string(candidateJSON),
			}
			if err := writer.WriteJSON(msg); err != nil {
				log.Println("Erro ao enviar ICE candidate:", err)
			}
		})

		peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
			log.Printf("Peer recebeu track: ID=%s, Kind=%s", remoteTrack.ID(), remoteTrack.Kind().String())
			localTrack, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, remoteTrack.ID(), remoteTrack.StreamID())
			if err != nil {
				log.Println("Erro ao criar localTrack:", err)
				return
			}

			room.Peers.ListLock.Lock()
			room.Peers.TrackLocals[remoteTrack.ID()] = localTrack
			room.Peers.ListLock.Unlock()

			room.Peers.SignalPeerConnections()

			rtcpBuf := make([]byte, 1500)
			for {
				packet, _, err := remoteTrack.ReadRTP()
				if err != nil {
					log.Println("Erro ao ler RTP:", err)
					break
				}
				if err := localTrack.WriteRTP(packet); err != nil {
					log.Println("Erro ao escrever RTP:", err)
				}
				if _, _, err := receiver.Read(rtcpBuf); err != nil {
					break
				}
			}
		})

		var candidateBuffer []webrtc.ICECandidateInit

		for {
			_, msgBytes, err := ws.ReadMessage()
			if err != nil {
				log.Println("Erro ao ler mensagem WebSocket:", err)
				break
			}
			var msg websocketMessage
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Println("Erro ao decodificar mensagem:", err)
				continue
			}
			switch msg.Event {
			case "offer":
				var offer webrtc.SessionDescription
				if err := json.Unmarshal([]byte(msg.Data), &offer); err != nil {
					log.Println("Erro ao decodificar offer:", err)
					continue
				}
				if err := peerConnection.SetRemoteDescription(offer); err != nil {
					log.Println("Erro ao definir remote description:", err)
					continue
				}

				for _, c := range candidateBuffer {
					if err := peerConnection.AddICECandidate(c); err != nil {
						log.Println("Erro ao adicionar candidato buffered:", err)
					}
				}
				candidateBuffer = nil

				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					log.Println("Erro ao criar answer:", err)
					continue
				}
				if err := peerConnection.SetLocalDescription(answer); err != nil {
					log.Println("Erro ao definir local description:", err)
					continue
				}
				<-webrtc.GatheringCompletePromise(peerConnection)
				answerJSON, err := json.Marshal(*peerConnection.LocalDescription())
				if err != nil {
					log.Println("Erro ao serializar answer:", err)
					continue
				}
				response := websocketMessage{
					Event: "answer",
					Data:  string(answerJSON),
				}
				if err := writer.WriteJSON(response); err != nil {
					log.Println("Erro ao enviar answer:", err)
				}
			case "candidate":
				var candidate webrtc.ICECandidateInit
				if err := json.Unmarshal([]byte(msg.Data), &candidate); err != nil {
					log.Println("Erro ao decodificar candidate:", err)
					continue
				}
				if peerConnection.RemoteDescription() == nil {
					candidateBuffer = append(candidateBuffer, candidate)
					log.Println("Buffering ICE candidate, remote description ainda não definida.")
				} else {
					if err := peerConnection.AddICECandidate(candidate); err != nil {
						log.Println("Erro ao adicionar candidate:", err)
					}
				}
			}
		}

		room.Peers.RemovePeer(newPeer)
	}
}

func (p *Peers) AddPeer(peer PeerConnectionState) {
	p.ListLock.Lock()
	p.Connections = append(p.Connections, peer)
	p.ListLock.Unlock()
	p.SignalPeerConnections()
}

func (p *Peers) RemovePeer(peer PeerConnectionState) {
	p.ListLock.Lock()
	defer p.ListLock.Unlock()
	for i, conn := range p.Connections {
		if conn == peer {
			p.Connections = append(p.Connections[:i], p.Connections[i+1:]...)
			break
		}
	}
	p.SignalPeerConnections()
}

func (p *Peers) SignalPeerConnections() {
	p.ListLock.Lock()
	defer p.ListLock.Unlock()

	for i := range p.Connections {
		pc := p.Connections[i].PeerConnection

		existingTracks := map[string]bool{}
		for _, sender := range pc.GetSenders() {
			if sender.Track() == nil {
				continue
			}
			existingTracks[sender.Track().ID()] = true
			if _, ok := p.TrackLocals[sender.Track().ID()]; !ok {
				if err := pc.RemoveTrack(sender); err != nil {
					log.Println("Erro ao remover track:", err)
				}
			}
		}

		for trackID, track := range p.TrackLocals {
			if !existingTracks[trackID] {
				if _, err := pc.AddTrack(track); err != nil {
					log.Println("Erro ao adicionar track:", err)
				}
			}
		}

		offer, err := pc.CreateOffer(nil)
		if err != nil {
			log.Println("Erro ao criar offer:", err)
			continue
		}
		if err := pc.SetLocalDescription(offer); err != nil {
			log.Println("Erro ao definir local description:", err)
			continue
		}
		offerJSON, err := json.Marshal(offer)
		if err != nil {
			log.Println("Erro ao serializar offer:", err)
			continue
		}
		msg := websocketMessage{
			Event: "offer",
			Data:  string(offerJSON),
		}
		if err := p.Connections[i].Websocket.WriteJSON(msg); err != nil {
			log.Println("Erro ao enviar offer:", err)
		}
	}
}

func DispatchKeyFrames(p *Peers) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		p.ListLock.Lock()
		for _, conn := range p.Connections {
			for _, receiver := range conn.PeerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}
				err := conn.PeerConnection.WriteRTCP([]rtcp.Packet{
					&rtcp.PictureLossIndication{MediaSSRC: uint32(receiver.Track().SSRC())},
				})
				if err != nil {
					log.Println("Erro ao enviar RTCP PLI:", err)
				}
			}
		}
		p.ListLock.Unlock()
	}
}
