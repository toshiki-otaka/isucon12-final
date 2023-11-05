package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		panic("specify shard count")
	}
	cnt, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	f, err := os.Open("./webapp/sql/4_alldata_exclude_user_presents.sql")
	if err != nil {
		fmt.Println("Run:")
		fmt.Println("$ wget https://github.com/isucon/isucon12-final/releases/download/initial_data_20220912/initial_data.tar.gz")
		fmt.Println("$ tar -xvf initial_data.tar.gz")
		panic(err)
	}
	defer f.Close()

	writers := make([]*os.File, cnt)
	for i := 0; i < cnt; i++ {
		f, err := os.Create("./webapp/sql/4_alldata_exclude_user_presents_" + strconv.Itoa(i) + ".sql")
		if err != nil {
			panic(err)
		}
		writers[i] = f
	}

	r := bufio.NewReader(f)
	for {
		b, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line := string(b)

		if !strings.HasPrefix(line, "INSERT INTO") {
			continue
		}

		end := 100
		if len(line) < end {
			end = len(line)
		}
		fmt.Println(len(line), string(line)[:end])
	}
}
