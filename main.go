package main

import (
	"fmt"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
)

func main() {
	tcpServer, err := net.ListenMC("localhost:25565")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		var p pk.Packet

		err = conn.ReadPacket(&p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if p.ID == 0x00 {
			var (
				protocolVersion pk.VarInt
				serverAddress   pk.String
				serverPort      pk.UnsignedShort
				nextState       pk.VarInt
			)

			err = p.Scan(&protocolVersion, &serverAddress, &serverPort, &nextState)
			if err != nil {
				fmt.Println(err)
				continue
			}

		}
	}

}
