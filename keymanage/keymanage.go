package keymanage

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strings"
)

// expected: root:x:0:0:root:/root:/bin/bash
func convertJson(str string) string {
	lines := strings.Split(str, "\n")
	for _, v := range lines {
		split := strings.Split(v, ":")
		for i, v := range split {
			if i == 0 {
				fmt.Println(v)
			}
		}
	}
	return ""
}

// パスワードの初期化　rootのみの実行？
func initPassword() {

}

// ユーザの一覧表示
func UsersList() {
	var cmd string = ""
	var output []byte = nil
	var err error = nil

	cmd = "cat /etc/passwd | grep /bin/bash"
	output, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
	fmt.Println(string(output))
}

// ユーザの追加
// すでに対象サーバにユーザが存在する場合は何もせずにreturnする
func UserAdd(username string) {
	var cmd string = ""
	var output []byte = nil
	var err error = nil

	// 追加しようとしているユーザがすでに存在するかチェックする
	cmd = "cat /etc/passwd | grep -e \"^" + username + ":\""
	output, _ = doCmd(cmd)
	// すでにユーザが存在する場合はreturn
	if output != nil && len(output) != 0 {
		log.Println(string(output))
		return
	}
	cmd = "sudo useradd " + username
	_, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
}

// ユーザの削除
func UserDel(username string) {
	// rootの扱い
	var cmd string = ""
	var output []byte = nil
	var err error = nil

	// 追加しようとしているユーザがすでに存在するかチェックする
	cmd = "cat /etc/passwd | grep -e \"^" + username + ":\""
	output, _ = doCmd(cmd)
	//　ユーザが存在しない場合はreturn
	if output != nil && len(output) == 0 {
		log.Println("target user is already deleted: " + username)
		return
	}
	cmd = "sudo userdel " + username
	_, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
}

// ユーザに権限追加
// sudoが付いているユーザにつけようとしても問題ない
func AuthAdd(username string) {
	var cmd string = ""
	// var output []byte = nil
	var err error = nil

	cmd = "sudo usermod -aG wheel " + username
	_, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
}

// ユーザから権限削除
//-Gオプション、-aオプションを利用しても、副グループの一部を削除することはできない。
//もしも複数登録している副グループのうちの１つを削除する場合は、次のような作業手順が考えられる。
//idコマンドで副グループ一覧を取得する
//副グループ一覧をクリップボードなどにコピーする。
//削除したい副グループを除し、「,」で繋ぐ。
//-Gオプションで削除した副グループを省いた一覧を副グループとして設定する。
//このように副グループの削除は少々面倒な作業になる。
func delAuth() {
	//var cmd string = ""
	//// var output []byte = nil
	//var err error = nil
	//
	//cmd = "sudo usermod -aG wheel"
	//_, err = doCmd(cmd)
	//if err != nil {
	//	log.Fatalf("[Execute Command] %v, %v", cmd , err)
	//}
}

func doCmd(cmd string) ([]byte, error) {
	ip := "127.0.0.1"
	port := "2222"
	user := "vagrant"
	key_path := "/Users/shin/git/keymanage/id_rsa"
	buf, err := ioutil.ReadFile(key_path)
	if err != nil {
		log.Fatalf("could not read keypath: %s", key_path)
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	conn, err := ssh.Dial("tcp", ip+":"+port, config)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	return session.Output(cmd)
}
