package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

var (
	index           int
	clientesSockets map[*net.UDPAddr]*net.UDPConn
	filmesCad       map[int][]string
	clientesCad     map[string]string
	clientesON      map[string]*net.UDPAddr
	clientesChat    map[string]*net.UDPAddr
)

func contemCliente(end *net.UDPAddr) (string, bool) {
	for id, endAux := range clientesON {
		if endAux == end {
			return id, true
		}
	}
	return "", false
}

func contemFilme(filme string, diretor string) bool {
	for ind := 1; ind <= index-1; ind++ {
		if filme == filmesCad[ind][0] && diretor == filmesCad[ind][1] {
			return true
		}
	}
	return false
}

func ver_erro(erro error) {
	if erro != nil {
		fmt.Println("Erro -> ", erro)
		fmt.Println("Saindo")
		panic("")
	}
}

func add_filme(filme string, descricao string, diretor string) {
	filmesCad[index] = append(filmesCad[index], filme)
	filmesCad[index] = append(filmesCad[index], diretor)
	filmesCad[index] = append(filmesCad[index], descricao)
	index++
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Offline(Socket *net.UDPConn, end_cliente *net.UDPAddr) {
	buffer_recMsg := make([]byte, 1024)
	buffer_recEsc := make([]byte, 1024)

	buffer_msg1 := []byte("\nBem vindo ao servidor\n----------------------\n|(1)->Entrar\n|(2)->Cadastre-se\n|(3)-> Sair\n|Digite sua escolha: ")

	_, erro := Socket.WriteToUDP(buffer_msg1, end_cliente)
	ver_erro(erro)
	tam, _, erro := Socket.ReadFromUDP(buffer_recEsc)
	ver_erro(erro)

	if string(buffer_recEsc[0:tam-1]) == "1" { //Opção de logar
		buffer_loginID := []byte("\nDigite seu ID: ")

		//Pede o ID do login
		_, erro = Socket.WriteToUDP(buffer_loginID, end_cliente)
		ver_erro(erro)
		tam, _, erro := Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Verificando se o ID está nos clientes cadastrados e recebe a senha relacionada
		senha, verifUsCad := clientesCad[string(buffer_recMsg[0:tam-1])]
		//Verificando se o ID já está inserido nos clientes online
		_, verifUsDif := clientesON[string(buffer_recMsg[0:tam-1])]

		if verifUsDif { //Se o ID estiver entre os online
			//Informa que a conta ja está logada
			buffer_loginIDError1 := []byte("Este ID já está logado!\n")

			_, erro = Socket.WriteToUDP(buffer_loginIDError1, end_cliente)
			ver_erro(erro)

		} else if verifUsCad { //Se o ID estiver cadastrado
			id := string(buffer_recMsg[0 : tam-1])

			buffer_loginSenha := []byte("Digite sua senha: ")
			_, erro = Socket.WriteToUDP(buffer_loginSenha, end_cliente)
			ver_erro(erro)

			tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
			ver_erro(erro)

			//Compara se a senha recebida é a mesma que a relacionada ao ID
			if string(buffer_recMsg[0:tam-1]) == senha {
				buffer_loginSucess := []byte("\nLogado com sucesso!\n")

				_, erro = Socket.WriteToUDP(buffer_loginSucess, end_cliente)
				ver_erro(erro)

				//Adiciona aos clientes que estão Online
				clientesON[id] = end_cliente
			} else { //Se não a senha digitada estava incorreta
				buffer_loginSenhaError := []byte("\nSenha incorreta!\n")

				_, erro = Socket.WriteToUDP(buffer_loginSenhaError, end_cliente)
				ver_erro(erro)
			}
		} else { //Caso não encontre o ID na lista, será informado que o usuário não está cadastrado
			buffer_loginIDError2 := []byte("Usuário não cadastrado!\n")

			_, erro = Socket.WriteToUDP(buffer_loginIDError2, end_cliente)
			ver_erro(erro)
		}
	} else if string(buffer_recEsc[0:tam-1]) == "2" { //Opção de cadastro
		buffer_cadID := []byte("\nDigite seu novo ID: ")

		//Pede o ID a ser cadastrado
		_, erro = Socket.WriteToUDP(buffer_cadID, end_cliente)
		ver_erro(erro)
		tam, _, erro := Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Verifica se o ID está cadastrado
		id := string(buffer_recMsg[0 : tam-1])
		_, verifUs := clientesCad[string(buffer_recMsg[0:tam-1])]
		if verifUs == false { //Caso não esteja cadastrado
			//Pede a senha
			buffer_cadSenha := []byte("Digite sua nova senha: ")

			_, erro = Socket.WriteToUDP(buffer_cadSenha, end_cliente)
			ver_erro(erro)
			tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
			ver_erro(erro)

			//Registra o cadastro
			senha := string(buffer_recMsg[0 : tam-1])
			clientesCad[id] = senha

			//Informa o sucesso no cadastro
			buffer_cadSucess := []byte("\nCadastrado com sucesso!\n")
			_, erro = Socket.WriteToUDP(buffer_cadSucess, end_cliente)
			ver_erro(erro)

		} else { //Se encontrar o ID cadastrado, informa o erro
			buffer_cadIDError := []byte("\nUsuário já existente!\n")

			_, erro = Socket.WriteToUDP(buffer_cadIDError, end_cliente)
			ver_erro(erro)
		}

	} else if string(buffer_recEsc[0:tam-1]) == "3" { //Opção de sair
		buffer_end := []byte("\nDesconectado do servidor! \nAgradecemos sua preferência :^)")

		_, erro = Socket.WriteToUDP(buffer_end, end_cliente)
		ver_erro(erro)

		//Deleta o Socket da lista de Socket's, o qual tratava este cliente
		delete(clientesSockets, end_cliente)
	} else {
		buffer_msgError := []byte("\nOpção inválida!\n")

		_, erro = Socket.WriteToUDP(buffer_msgError, end_cliente)
		ver_erro(erro)
	}
}

func Online(Socket *net.UDPConn, end_cliente *net.UDPAddr) {
	buffer_recMsg := make([]byte, 1024)
	buffer_recEsc := make([]byte, 1024)

	buffer_msg2 := []byte("\n----------------------------------\n|(1)-> Inverter String \n|(2)-> IMC\n|(3)-> Número Aleatório \n|(4)-> Mensagem direta \n|(5)-> Chat \n|(6)-> Cadastrar filme \n|(7)-> Filme Aleatório \n|(8)-> Sair\n|Digite sua escolha: ")
	buffer_input := []byte("\nDigite uma mensagem: ")

	_, erro := Socket.WriteToUDP(buffer_msg2, end_cliente)
	ver_erro(erro)
	tam, _, erro := Socket.ReadFromUDP(buffer_recEsc)
	ver_erro(erro)

	if string(buffer_recEsc[0:tam-1]) == "1" { //Opção de inverter as mensagens
		buffer_msgOp1 := []byte("\nA partir de agora, todas suas mensagens serão retornadas ao contrário (Digite 'SAIR' para voltar ao menu de escolhas)\n")

		_, erro = Socket.WriteToUDP(buffer_msgOp1, end_cliente)
		ver_erro(erro)

		tam = 1
		for string(buffer_recMsg[0:tam-1]) != "SAIR" {
			//Pede a mensagem
			_, erro = Socket.WriteToUDP(buffer_input, end_cliente)
			ver_erro(erro)
			tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
			ver_erro(erro)

			if string(buffer_recMsg[0:tam-1]) != "SAIR" {
				//Inverte, salva e envia
				msg_rev := Reverse(string(buffer_recMsg[0 : tam-1]))
				msg_env := []byte(msg_rev)
				_, erro = Socket.WriteToUDP(msg_env, end_cliente)
			}
		}
	} else if string(buffer_recEsc[0:tam-1]) == "2" { //Opção de calcular IMC
		buffer_msgOp2 := []byte("\n-------Você entrou na calculadora de IMC-------\n")
		buffer_msgOp2Alt := []byte("Informe sua altura: ")
		buffer_msgOp2Peso := []byte("Informe seu peso: ")
		buffer_msgOp2IMC := []byte("Seu IMC é: ")

		_, erro = Socket.WriteToUDP(buffer_msgOp2, end_cliente)
		ver_erro(erro)

		_, erro = Socket.WriteToUDP(buffer_msgOp2Alt, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Após pedir a altura, a salva como string e depois a converte em float
		AltStr := string(buffer_recMsg[0 : tam-1])
		AltFloat, _ := strconv.ParseFloat(AltStr, 64)

		_, erro = Socket.WriteToUDP(buffer_msgOp2Peso, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Após pedir o peso, o salva como string e depois a converte em float
		PesoStr := string(buffer_recMsg[0 : tam-1])
		PesoFloat, _ := strconv.ParseFloat(PesoStr, 64)

		//Realiza o cálculo
		IMCFloat := PesoFloat / (AltFloat * AltFloat)

		//Converte o float em string e depois o transforma em byte para enviar ao usuário
		IMCStr := strconv.FormatFloat(IMCFloat, 'f', 6, 64)
		buffer_envio := []byte("IMC = " + IMCStr)
		_, erro = Socket.WriteToUDP(buffer_msgOp2IMC, end_cliente)
		ver_erro(erro)
		_, erro = Socket.WriteToUDP(buffer_envio, end_cliente)
		ver_erro(erro)

		//Informa a condição do usuário com base no resultado
		if IMCFloat < 18.5 {
			buffer_msgOp2Caso1 := []byte("\nVocê está abaixo do peso!\n\n")
			_, erro = Socket.WriteToUDP(buffer_msgOp2Caso1, end_cliente)
			ver_erro(erro)

		} else if IMCFloat < 24.9 {
			buffer_msgOp2Caso2 := []byte("\nVocê está com o peso normal!\n\n")
			_, erro = Socket.WriteToUDP(buffer_msgOp2Caso2, end_cliente)
			ver_erro(erro)

		} else if IMCFloat < 29.9 {
			buffer_msgOp2Caso3 := []byte("\nVocê está com sobrepeso!\n\n")
			_, erro = Socket.WriteToUDP(buffer_msgOp2Caso3, end_cliente)
			ver_erro(erro)

		} else if IMCFloat < 34.9 {
			buffer_msgOp2Caso4 := []byte("\nVocê está com obesidade de 1° grau!\n")
			_, erro = Socket.WriteToUDP(buffer_msgOp2Caso4, end_cliente)
			ver_erro(erro)

		} else if IMCFloat < 39.9 {
			buffer_msgOp2Caso5 := []byte("\nVocê está com obesidade de 2° grau!\n")
			_, erro = Socket.WriteToUDP(buffer_msgOp2Caso5, end_cliente)
			ver_erro(erro)

		} else {
			buffer_msgOp2Caso6 := []byte("\nVocê está com obesidade morbida!\n")
			_, erro = Socket.WriteToUDP(buffer_msgOp2Caso6, end_cliente)
			ver_erro(erro)
		}

	} else if string(buffer_recEsc[0:tam-1]) == "3" { //Opção de número aleatório
		buffer_msgOp3 := []byte("\nDigite dois numeros separados por uma virgula (Ex: 14,8)\n")

		_, erro = Socket.WriteToUDP(buffer_msgOp3, end_cliente)
		ver_erro(erro)
		_, erro = Socket.WriteToUDP(buffer_input, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Salva a resposta em uma string
		Resp := string(buffer_recMsg[:tam-1])

		//Salva os números em inteiros
		num1, _ := strconv.Atoi(Resp[:strings.Index(Resp, ",")])
		num2, _ := strconv.Atoi(Resp[strings.Index(Resp, ",")+1:])

		//Analisa qual o maior e menor número, para então informar o número aleatório
		if num1 == num2 {
			buffer_envio := []byte("O número escolhido foi: " + strconv.Itoa(num1))
			_, erro = Socket.WriteToUDP(buffer_envio, end_cliente)
			ver_erro(erro)
		} else if num1 <= num2 {
			numAlt := rand.Intn(num2-num1) + num1
			buffer_envio := []byte("O número escolhido foi: " + strconv.Itoa(numAlt))
			_, erro = Socket.WriteToUDP(buffer_envio, end_cliente)
			ver_erro(erro)
		} else {
			numAlt := rand.Intn(num1-num2) + num2
			buffer_envio := []byte("O número escolhido foi: " + strconv.Itoa(numAlt))
			_, erro = Socket.WriteToUDP(buffer_envio, end_cliente)
			ver_erro(erro)
		}

	} else if string(buffer_recEsc[0:tam-1]) == "4" { //Opção de mensagem direta
		_, erro = Socket.WriteToUDP(buffer_input, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		env_msg := string(buffer_recMsg[:tam-1])

		//Pergunta o usuário destino
		buffer_msgOp4 := []byte("\nDigite o ID destino da mensagem: ")
		_, erro = Socket.WriteToUDP(buffer_msgOp4, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Busca o ID destino e origem
		idDest := string(buffer_recMsg[0 : tam-1])
		idOrig, _ := contemCliente(end_cliente)

		//Verifica se o usuário destino está cadastrado
		_, verifUs := clientesCad[idDest]
		if verifUs {
			//Verifica se o usuário destino está online para receber e pega seu endereço
			end_dest, verifUs := clientesON[idDest]
			if idDest == idOrig { //Se o ID de origem for o mesmo do destino, informa ao usuário o equivoco
				buffer_msgOp4Error2 := []byte("\nEste é seu proprio ID!\n")
				_, erro = Socket.WriteToUDP(buffer_msgOp4Error2, end_cliente)
				ver_erro(erro)

			} else if verifUs { //Se o usuário destino estiver online
				buffer_envio := []byte("\nMensagem de " + idOrig + "!\nA mensagem é: " + env_msg)

				//Busca o Socket que está tratando o usuário destino
				SocketDest, _ := clientesSockets[end_dest]

				//Manda a mensagem
				_, erro = SocketDest.WriteToUDP(buffer_envio, end_dest)
				ver_erro(erro)
			} else { //Caso o usuário destino não esteja online, informa ao usuário origem
				buffer_msgOp4Error3 := []byte("\nUsuário " + idDest + " não está online, tente mais tarde!\n")
				_, erro = Socket.WriteToUDP(buffer_msgOp4Error3, end_cliente)
				ver_erro(erro)
			}
		} else { //Informa que o usuário não existe
			buffer_msgOp4Error1 := []byte("\nUsuário não existente\n")
			_, erro := Socket.WriteToUDP(buffer_msgOp4Error1, end_cliente)
			ver_erro(erro)

		}

	} else if string(buffer_recEsc[0:tam-1]) == "5" { //Opção do chat
		buffer_msgOp5 := []byte("\nVocê está entrando no chat! (Digite 'SAIR' quando quiser voltar ao menu de opções)\n")
		buffer_msgOp5Conect := []byte("Conectado ao chat\n\n")
		_, erro = Socket.WriteToUDP(buffer_msgOp5, end_cliente)
		ver_erro(erro)
		_, erro = Socket.WriteToUDP(buffer_msgOp5Conect, end_cliente)
		ver_erro(erro)

		//Registra o usuário na mapa que contem os participantes do chat
		ID, _ := contemCliente(end_cliente)
		clientesChat[ID] = end_cliente

		Conectado := true
		for Conectado {
			//Ficará recebendo as mensagens
			tam, _, erro := Socket.ReadFromUDP(buffer_recMsg)
			ver_erro(erro)

			Conectado = string(buffer_recMsg[0:tam-1]) != "SAIR"

			//Percorre todos o clientes que estão no chat
			for _, endDest := range clientesChat {
				if endDest != end_cliente && Conectado { //Verifica se o endereço é diferente da origem
					buffer_envMsg := []byte(ID + ": " + string(buffer_recMsg[:tam]) + "\n")

					//Recebe o Socket que está tratando o cliente
					SocketChat := clientesSockets[endDest]
					_, erro = SocketChat.WriteToUDP(buffer_envMsg, endDest)
					ver_erro(erro)
				}
			}
		}
		delete(clientesChat, ID)
	} else if string(buffer_recEsc[:tam-1]) == "6" {
		buffer_msgOp6 := []byte("\n----------Bem vindo ao Cadastro de Filmes----------")
		buffer_msgOp6Nome := []byte("\nDigite o nome do filme a ser cadastrado: ")
		buffer_msgOp6Descr := []byte("Digite a descrição do filme: ")
		buffer_msgOp6Diret := []byte("Digite o nome do diretor do filme: ")
		buffer_msgOp6Error := []byte("Falha na tentativa! Filme já cadastrado!\n")
		buffer_msgOp6Sucess := []byte("\nFilme cadastrado! Obrigado :^)\n")

		//Receber nome do filme
		_, erro = Socket.WriteToUDP(buffer_msgOp6, end_cliente)
		ver_erro(erro)
		_, erro = Socket.WriteToUDP(buffer_msgOp6Nome, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)
		nome := string(buffer_recMsg[:tam-1])

		//Receber diretor do filme
		_, erro = Socket.WriteToUDP(buffer_msgOp6Diret, end_cliente)
		ver_erro(erro)
		tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)
		diret := string(buffer_recMsg[:tam-1])

		//Verifica se o filme não está cadastrado
		if contemFilme(nome, diret) != true {
			//Receber descricao do filme
			_, erro = Socket.WriteToUDP(buffer_msgOp6Descr, end_cliente)
			ver_erro(erro)
			tam, _, erro = Socket.ReadFromUDP(buffer_recMsg)
			ver_erro(erro)
			descr := string(buffer_recMsg[:tam-1])

			//Cadastra o filme
			add_filme(nome, descr, diret)
			_, erro = Socket.WriteToUDP(buffer_msgOp6Sucess, end_cliente)
			ver_erro(erro)
		} else {
			//Caso esteja informa ao usuário
			_, erro = Socket.WriteToUDP(buffer_msgOp6Error, end_cliente)
			ver_erro(erro)
		}

	} else if string(buffer_recEsc[:tam-1]) == "7" {
		buffer_end := []byte("\nFilme selecionado\n")
		_, erro = Socket.WriteToUDP(buffer_end, end_cliente)
		ver_erro(erro)

		//Gera um número aleatório
		numAlt := rand.Intn(index-1) + 1

		//Usa o número para buscar o filme relacionado
		buffer_nome := []byte("O nome do filme: " + filmesCad[numAlt][0])
		_, erro = Socket.WriteToUDP(buffer_nome, end_cliente)
		ver_erro(erro)
		buffer_diret := []byte("\nDirigido por: " + filmesCad[numAlt][1])
		_, erro = Socket.WriteToUDP(buffer_diret, end_cliente)
		ver_erro(erro)
		buffer_descr := []byte("\nDescrição: " + filmesCad[numAlt][2])
		_, erro = Socket.WriteToUDP(buffer_descr, end_cliente)
		ver_erro(erro)

	} else if string(buffer_recEsc[:tam-1]) == "8" {
		buffer_deconect := []byte("\nDesconectado do servidor! \nAgradecemos sua preferência :^)")
		_, erro = Socket.WriteToUDP(buffer_deconect, end_cliente)
		ver_erro(erro)

		//Retira o cliente dos onlines e apaga seu Socket
		id, _ := contemCliente(end_cliente)
		delete(clientesON, id)
		delete(clientesSockets, end_cliente)
		ver_erro(erro)
	} else {
		buffer_end := []byte("\nOpção inválida. Tente novamente!")
		_, erro = Socket.WriteToUDP(buffer_end, end_cliente)
		ver_erro(erro)
	}
}

func server(Socket *net.UDPConn, end_cliente *net.UDPAddr) {
	for {
		//Analisando se o endereço do cliente está registrado nos clientes online
		_, verOnline := contemCliente(end_cliente)
		if verOnline == false {
			Offline(Socket, end_cliente)
		}

		if verOnline == true {
			Online(Socket, end_cliente)
		}
	}
}

func main() {
	//Iniciando as variaveis globais
	index = 3
	clientesChat = make(map[string]*net.UDPAddr)
	clientesSockets = make(map[*net.UDPAddr]*net.UDPConn)
	filmesCad = make(map[int][]string)
	clientesON = make(map[string]*net.UDPAddr)
	clientesCad = make(map[string]string)

	//Colocando alguns filmes nas listas e alguns cadastros de usuarios, para facilitar os testes
	clientesCad["carlos"] = "3711"
	clientesCad["joao"] = "3422"
	filmesCad[1] = append(filmesCad[1], "Magico de OZ")
	filmesCad[1] = append(filmesCad[1], "Victor Fleming")
	filmesCad[1] = append(filmesCad[1], "Dorothy e seu cachorro Totó são levados para a terra mágica de Oz quando um ciclone passa pela fazenda de seus avós no Kansas. Eles viajam em direção à Cidade Esmeralda para encontrar o Mago Oz e no caminho encontram um Espantalho, que precisa de um cérebro, um Homem de Lata sem um coração e um Leão Covarde que quer coragem. O Mago pede ao grupo que tragam a vassoura da Bruxa Malvada do Oeste a fim de ganharem sua juda.")

	filmesCad[2] = append(filmesCad[2], "Transformers")
	filmesCad[2] = append(filmesCad[2], "Michael Bay")
	filmesCad[2] = append(filmesCad[2], "O destino da humanidade está em jogo quando duas raças de robôs, os Autobots e os vilões Decepticons, chegam à Terra. Os robôs possuem a capacidade de se transformarem em diferentes objetos mecânicos enquanto buscam a chave do poder supremo com a ajuda do jovem Sam.")

	//Criando Socket para receber clientes
	fmt.Println("Iniciando o servidor...")
	End, erro := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	ver_erro(erro)

	fmt.Println("Enderenço resolvido")
	SocketGlob, erro := net.ListenUDP("udp", End)
	ver_erro(erro)

	fmt.Println("Servidor iniciado")

	for i := 3; ; i++ {
		buffer_recMsg := make([]byte, 1024)

		//Recebe o cliente
		tam, end_cliente, erro := SocketGlob.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Criando o endereço para o Socket que irá tratar o cliente
		EndSocketCliente := "127.0.0.1:1000" + strconv.Itoa(i)
		buffer_SocketID := []byte(EndSocketCliente)

		//Criando o novo Socket
		_, erro = SocketGlob.WriteToUDP(buffer_SocketID, end_cliente)
		End, erro := net.ResolveUDPAddr("udp", EndSocketCliente)
		SocketCliente, erro := net.ListenUDP("udp", End)

		//Esperando uma mensagem do cliente para se conectar
		tam, end_cliente, erro = SocketCliente.ReadFromUDP(buffer_recMsg)
		ver_erro(erro)

		//Registra o Socket na lista dos Socket's e chama a função principal
		clientesSockets[end_cliente] = SocketCliente
		fmt.Println("CLIENT ", end_cliente, "-->", string(buffer_recMsg[0:tam]))
		go server(SocketCliente, end_cliente)
	}
}
