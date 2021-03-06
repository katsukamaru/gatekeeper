package keymanage

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strings"
)

type Server struct {
	Name  string
	Users []User
}

type User struct {
	Name string
	Sudo bool
}

// expected: root:x:0:0:root:/root:/bin/bash
func convertJson(str string) []string {
	var list []string
	lines := strings.Split(str, "\n")
	for _, v := range lines {
		split := strings.Split(v, ":")
		for i, v := range split {
			if i == 0 {
				fmt.Println(v)
				list = append(list, v)
			}
		}
	}
	return list
}

// expected: wheel:x:10:katsukamaru,minami
func wheelUsers(str string) []string {
	trimed := strings.Split(str, "\n")
	split := strings.Split(trimed[0], ":")
	return strings.Split(split[3], ",")
}

// パスワードの初期化　rootのみの実行？
func initPassword() {

}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// ユーザの一覧表示
func UsersList() []User {
	var cmd = ""
	var output []byte = nil
	var err error = nil

	cmd = "cat /etc/passwd | grep /bin/bash"
	output, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
	list := convertJson(string(output))

	cmd = "cat /etc/group | grep wheel"
	output, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
	wheelUsers := wheelUsers(string(output))

	var Users []User
	for _, value := range list {
		// root は除く
		if value == "root" || value == "" {
			continue
		}
		if contains(wheelUsers, value) {
			Users = append(Users, User{value, true})
			continue
		}
		Users = append(Users, User{value, false})
	}
	return Users
}

func setSshConf() {

	// mkdir ~/.ssh
	// cat pub >> ~/.ssh/authorized_keys
	// chown -R username:Group ~/.ssh
	// chmod 700 ~/.ssh
	// chmod 600 ~/.ssh/authorized_keys

}

// TODO カンマ区切りで入力される想定に対応する
// ユーザの追加
// すでに対象サーバにユーザが存在する場合は何もせずにreturnする
func UserAdd(username string) {
	var cmd = ""
	var output []byte = nil
	var err error = nil

	// var pub = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDGc9KWcrVFgbFuVxMkU4rbfr3wvdI+eZxzb4Ys2BVVt+YUt+RdMcynWq22wofLy5wXyAQDJcaBItytX3TPHCP0HBbd/j403LIUHgFq7IpmfrHNkSs4PkpKsLPt3GTZmBwkqfFH9+6myg237RRUVLk80Rgz1V+JunVGdyc6L8KLqeCZ9xCaL3MkZIz8nm8GmhNNMvHXyxsMQMKjMA5uSQft6Xr12NVuL6we4/qIYS/9GHavjZ4lXj90vUfiHBOd4CoTh+wPDqF4gqHt5Ds8k/ObP0OIKoC/d/angjcAlykjNwaCjamhHC2Tb83hfyLVzbuQq68KnNK+7QVD6ypIFLPH"

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
	var cmd = ""
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
	var cmd = ""
	// var output []byte = nil
	var err error = nil

	cmd = "sudo usermod -aG wheel " + username
	_, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
}

// ユーザから権限削除
// id -nG minami -> minami wheel
// sudo usermod -G minami username
func delAuth(username string) {
	var cmd = ""
	var output []byte = nil
	var err error = nil
	cmd = "id -nG " + username
	output, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}

	groups := string(output)
	split := strings.Split(groups, " ")
	var result []string
	for _, num := range split {
		if num != "wheel" {
			result = append(result, num)
		}
	}
	cmd = "sudo usermod -G minami " + username
	_, err = doCmd(cmd)
	if err != nil {
		log.Fatalf("[Execute Command] %v, %v", cmd, err)
	}
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

// --- Must
// view
// set .ssh/pub key
// AuthAdd and AuthDel

// --- Want
// conf
// initial password
// test
// Login
// package
