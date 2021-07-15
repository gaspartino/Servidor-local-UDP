package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func ver_erro(erro error) {
	if erro != nil {
		fmt.Println("Erro -> ", erro)
		fmt.Println("Saindo")
		panic("")
	}
}

func ler_servidor(Socket *net.UDPConn) {
	for {
		buffer_rec := make([]byte, 1024)
		tam, _, erro := Socket.ReadFromUDP(buffer_rec)
		ver_erro(erro)
		fmt.Printf("%s", string(buffer_rec[:tam]))

		//Caso receba a mensagem de saida do servidor, ele encerra o cliente
		if string(buffer_rec[0:tam]) == "\nDesconectado do servidor! \nAgradecemos sua preferência :^)" {
			os.Exit(0)
		}
	}
}

func main() {
	fmt.Println("Iniciando cliente")
	End_servidor, erro := net.ResolveUDPAddr("udp", ":10001")
	ver_erro(erro)
	End_local, erro := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	ver_erro(erro)
	fmt.Println("Endereço foi resolvido")
	//Abre o primeiro Socket para ser recebido pelo servidor
	Socket, erro := net.DialUDP("udp", End_local, End_servidor)
	ver_erro(erro)
	fmt.Println("Cliente foi iniciado")

	//Conecta ao socket
	buffer_rec := make([]byte, 1024)
	buffer_conect := []byte("Conectado")
	_, erro = Socket.Write(buffer_conect)

	//Recebe o endereço do Socket que irá tratar este usuário
	tam, End_servidor, erro := Socket.ReadFromUDP(buffer_rec)
	End_servidorInd := string(buffer_rec[:tam])

	//Fecha o primeiro socket
	Socket.Close()

	//Abre um novo socket para ser tratado pelo servidor
	End_servidor, erro = net.ResolveUDPAddr("udp", End_servidorInd)
	ver_erro(erro)
	Socket, erro = net.DialUDP("udp", End_local, End_servidor)
	ver_erro(erro)

	//Conecta ao cliente
	_, erro = Socket.Write(buffer_conect)

	//Inicia uma go routine para receber as mensagens do servidor
	go ler_servidor(Socket)

	//Cria um loop para ficar respondendo o servidor
	buff_scan := bufio.NewReader(os.Stdin)
	for {
		msg_string, erro := buff_scan.ReadString('\n')
		ver_erro(erro)

		msg_string = strings.TrimSuffix(msg_string, "\n")
		buffer_envio := []byte(msg_string)

		_, erro = Socket.Write(buffer_envio)
		ver_erro(erro)
	}
}
