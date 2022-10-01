package main

import (
	"encoding/json"
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/google/uuid"
)

func main() {
	tcpServer, err := net.ListenMC(":25565")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		clientConn, err := tcpServer.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		var handshakePacket pk.Packet

		err = clientConn.ReadPacket(&handshakePacket)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if handshakePacket.ID == 0x00 { //Handshake
			var (
				protocolVersion pk.VarInt
				serverAddress   pk.String
				serverPort      pk.UnsignedShort
				nextState       pk.VarInt
			)

			err = handshakePacket.Scan(&protocolVersion, &serverAddress, &serverPort, &nextState)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if nextState == 1 { //nextState 1 = Status
				var statusRequest pk.Packet
				clientConn.ReadPacket(&statusRequest)
				clientConn.WritePacket(pk.Marshal(0x00, pk.String(serverListPayload{
					Version: version{
						Name:     "1.16.4",
						Protocol: 754,
					},
					Players: players{
						Max:    10,
						Online: 1,
						Sample: []sample{
							{
								Name: "Working...",
								Id:   uuid.New().String(),
							},
						},
					},
					Description: chat.Text("§cKStars Proxy Server"),
				}.marchal(),
				)))

				var pingRequest pk.Packet
				clientConn.ReadPacket(&pingRequest)
				var pingPayload pk.Long
				pingRequest.Scan(&pingPayload)

				clientConn.WritePacket(pk.Marshal(0x01, pingPayload))
			} else if nextState == 2 {
				connMc, err := net.DialMC("raspberrypi:25565")
				if err != nil {
					fmt.Println(err)
					continue
				}

				connMc.WritePacket(handshakePacket)

				errorCloseProxy := false

				go func() {
					for errorCloseProxy == false {
						var readPacket pk.Packet
						err := clientConn.ReadPacket(&readPacket)
						if err != nil {
							fmt.Println(err)
							errorCloseProxy = true
							break
						}
						connMc.WritePacket(readPacket)
					}
				}()

				for errorCloseProxy == false {
					var writePacket pk.Packet
					err := connMc.ReadPacket(&writePacket)
					if err != nil {
						fmt.Println(err)
						errorCloseProxy = true
						break
					}
					clientConn.WritePacket(writePacket)
				}
			}
		}
	}

}

type serverListPayload struct {
	Version     version      `json:"version"`
	Players     players      `json:"players"`
	Description chat.Message `json:"description"`
	Favicon     string       `json:"favicon"`
}

type version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []sample `json:"sample"`
}

type sample struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (s serverListPayload) marchal() string {
	marshal, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(marshal)
}
