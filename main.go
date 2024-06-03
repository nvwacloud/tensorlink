/*
 * Tensorlink
 * Copyright (C) 2024 Andy <550896603@qq.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const (
	LightProtocol  = "light"
	NativeProtocol = "native"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	role := flag.String("role", "", "Role of the program: 'server' for server, 'client' for client")
	sendPort := flag.String("send_port", "", "Send port")
	recvPort := flag.String("recv_port", "", "Receive port")
	protocol := flag.String("net", "", "Protocol type: 'light' or 'native'")
	serverIP := flag.String("ip", "", "Server IP address (client mode only)")
	flag.Parse()

	switch *role {
	case "server":
		if *sendPort == "" || *recvPort == "" || *protocol == "" {
			printServerUsage()
			return
		}
		if *protocol != LightProtocol && *protocol != NativeProtocol {
			printServerUsage()
			return
		}
		startServer(*sendPort, *recvPort, *protocol)
	case "client":
		if *serverIP == "" || *sendPort == "" || *recvPort == "" || *protocol == "" {
			printClientUsage()
			return
		}
		if *protocol != LightProtocol && *protocol != NativeProtocol {
			printClientUsage()
			return
		}
		startClient(*serverIP, *sendPort, *recvPort, *protocol)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("  -role server -net [protocol] -recv_port [receive port] -send_port [send port]  Start as server")
	fmt.Println("  -role client -ip [server ip] -net [protocol] -send_port [send port] -recv_port [receive port]  Start as client")
}

func printServerUsage() {
	fmt.Println("Server Usage: ")
	fmt.Println("  -role server -net [protocol] -recv_port [receive port] -send_port [send port]")
	fmt.Println("Example: ")
	fmt.Println("  -role server -net native -recv_port 9998 -send_port 9999")
}

func printClientUsage() {
	fmt.Println("Client Usage: ")
	fmt.Println("  -role client -ip [server ip] -net [protocol] -send_port [send port] -recv_port [receive port]")
	fmt.Println("Example: ")
	fmt.Println("  -role client -ip 192.168.2.2 -net native -send_port 9998 -recv_port 9999")
}

func startServer(sendPort, recvPort, protocol string) {
	args := []string{"-s", sendPort, "-r", recvPort, "-n", protocol}
	runVCUDACommand("./vcuda", args...)
}

func startClient(serverIP, sendPort, recvPort, protocol string) {
	args := []string{serverIP, protocol, sendPort, recvPort}
	runVCUDACommand("./vcuda.exe", args...)
}

func runVCUDACommand(path string, args ...string) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to start vcuda: %v\n", err)
	}
}
