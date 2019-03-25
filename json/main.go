package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type MyStructSub struct {
	Sub1 int32
	Sub2 uint16
	Rsb1 uint16
	Sub3 uint8
	Rsb2 [3]uint8
}

type MyStruct struct {
	Data1 int32
	Data2 uint16
	Data3 uint8
	Sub   MyStructSub
}

func JsontoBinary(src, dst string) error {

	// ファイルオープン
	file, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("os.OpenFile", err)
		return err
	}
	defer file.Close()

	// 構造体をバイナリにする
	var mystruct MyStruct
	err = binary.Read(file, binary.LittleEndian, &mystruct)
	if err != nil {
		log.Println("binary.Read", err)
		return err
	}

	// 構造体をJSONに変換
	jsonBytes, err := json.Marshal(mystruct)
	if err != nil {
		log.Println("json.Marshal", err)
		return err
	}

	// プリフィックスなし、スペース4つでインデント
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	buf, err := ioutil.ReadAll(out)
	if err != nil {
		log.Println("ioutil.ReadAll", err)
		return err
	}

	// ファイルへ書き込み
	err = ioutil.WriteFile(dst, buf, 0666)
	if err != nil {
		log.Println("ioutil.WriteFile", err)
		return err
	}

	return nil
}

func BinarytoJson(src, dst string) error {

	// ファイル読み込み
	bytes, err := ioutil.ReadFile(src)

	// JSON を 構造体に変換
	var mystruct MyStruct
	err = json.Unmarshal(bytes, &mystruct)
	if err != nil {
		fmt.Println("json.Unmarshal:", err)
		return err
	}

	// ファイルオープン
	file, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 構造体をバイナリに変換してファイルへ書き込む
	err = binary.Write(file, binary.LittleEndian, &mystruct)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	src1 := "src.bin"
	src2 := "src.json"
	dst1 := "dst.json"
	dst2 := "dst.bin"

	// ファイルがなければ新規作成
	f, err := os.Stat(src1)
	if os.IsNotExist(err) || f.IsDir() {
		var mystruct MyStruct
		file, err := os.OpenFile(src1, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer file.Close()
		err = binary.Write(file, binary.LittleEndian, &mystruct)
		if err != nil {
			return
		}
	}

	// json -> bin
	err = JsontoBinary(src1, dst1)
	if err != nil {
		return
	}

	// ファイルがなければ新規作成
	f, err = os.Stat(src2)
	if os.IsNotExist(err) || f.IsDir() {

		src, err := os.Open(dst1)
		if err != nil {
			return
		}
		defer src.Close()

		dst, err := os.Create(src2)
		if err != nil {
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return
		}
	}

	// json -> bin
	err = BinarytoJson(src2, dst2)
	if err != nil {
		return
	}

}
