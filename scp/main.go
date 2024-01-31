package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

// 準備
// - 秘密鍵ファイルを用意する
// 	 - mkdir .ssh && ssh-keygen -t ed25519 -f .ssh/go-scp
// - SSHサーバーを用意する
// 	 - docker-compose build --no-cache && docker-compose up -d
// - DockerコンテナにSSHでログインする
// 	 - ssh -i .ssh/go-scp casone@localhost -p 20021
// - Dockerコンテナにファイルを転送する(scpが使えることを確認する)
// 	 - scp -i .ssh/go-scp -P 20021 main.go casone@localhost:/home/casone

func main() {
	// 秘密鍵ファイルの読み込み
	pk, err := os.ReadFile(".ssh/go-scp")
	if err != nil {
		panic(err)
	}

	// 秘密鍵ファイルを解析し、署名情報を取得
	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		panic(err)
	}

	// SSHクライアントの設定を行う
	config := &ssh.ClientConfig{
		User: "casone",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// SSHサーバーに接続
	client, err := ssh.Dial("tcp", "localhost:20021", config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// セッションを開始
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// SSHサーバーにコマンドを送信
	// if err := session.Run("ls -la"); err != nil {
	//   panic(err)
	// }
	// bs, err := session.Output("ls -la")
	// if err != nil {
	//   panic(err)
	// }
	// println(string(bs))

	// リモートのファイル(testfile)をローカルに転送する
	// var buf bytes.Buffer
	// session.Stdout = &buf
	// if err := session.Run("cat testfile"); err != nil {
	//   panic(err)
	// }
	// if err := os.WriteFile("testfile", buf.Bytes(), 0o644); err != nil {
	//   panic(err)
	// }

	stdout, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	if err := session.Start("scp -t /home/casone"); err != nil {
		panic(err)
	}

	f, err := os.Open("Dockerfile")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintf(stdin, "C0644 %d %s\n", fi.Size(), "Dockerfile"); err != nil {
		panic(err)
	}
	recvRemoteStatus(stdout)

	if _, err := io.Copy(stdin, f); err != nil {
		panic(err)
	}
	recvRemoteStatus(stdout)
}

func recvRemoteStatus(r io.Reader) {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("recv: %v\n", buf[0])
}
