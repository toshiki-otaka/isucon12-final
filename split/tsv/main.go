package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		panic("specify shard count")
	}
	cnt, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	f, err := os.Open("../webapp/sql/5_user_presents_not_receive_data.tsv")
	if err != nil {
		fmt.Println("Run:")
		fmt.Println("$ wget https://github.com/isucon/isucon12-final/releases/download/initial_data_20220912/initial_data.tar.gz")
		fmt.Println("$ tar -xvf initial_data.tar.gz")
		fmt.Println("$ cp ./webapp/sql/* ../webapp/sql/")
		fmt.Println("$ rm -rf ./benchmarker/ ./webapp/")
		fmt.Println("$ rm initial_data.tar.gz")
		fmt.Println("")
		panic(err)
	}
	defer f.Close()

	writers := make([]*csv.Writer, cnt)
	for i := 0; i < cnt; i++ {
		f, err := os.Create("../webapp/sql/5_user_presents_not_receive_data_" + strconv.Itoa(i) + ".tsv")
		if err != nil {
			panic(err)
		}
		w := csv.NewWriter(f)
		w.Comma = '\t'
		defer w.Flush()
		writers[i] = w
	}

	r := csv.NewReader(f)
	r.Comma = '\t'

	var line int
	for {
		line++
		fmt.Println("line:", line)

		first := r.InputOffset() == 0

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if first {
			for _, w := range writers {
				if err := w.Write(record); err != nil {
					panic(err)
				}
			}
			continue
		}

		userID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			panic(err)
		}
		if err := writers[int(userID)%cnt].Write(record); err != nil {
			panic(err)
		}
	}
}
