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
	var lineNum int
	for {
		lineNum++
		fmt.Println("line:", lineNum)

		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if !strings.HasPrefix(line, "INSERT INTO") {
			for _, w := range writers {
				if _, err := w.WriteString(line); err != nil {
					panic(err)
				}
			}
			continue
		}

		var tablename string
		switch {
		case strings.HasPrefix(line, "INSERT INTO `users` VALUES "):
			tablename = "users"
		case strings.HasPrefix(line, "INSERT INTO `user_bans` VALUES "):
			tablename = "user_bans"
		case strings.HasPrefix(line, "INSERT INTO `user_decks` VALUES "):
			tablename = "user_decks"
		case strings.HasPrefix(line, "INSERT INTO `user_devices` VALUES "):
			tablename = "user_devices"
		case strings.HasPrefix(line, "INSERT INTO `user_login_bonuses` VALUES "):
			tablename = "user_login_bonuses"
		case strings.HasPrefix(line, "INSERT INTO `user_cards` VALUES "):
			tablename = "user_cards"
		case strings.HasPrefix(line, "INSERT INTO `user_items` VALUES "):
			tablename = "user_items"
		case strings.HasPrefix(line, "INSERT INTO `user_present_all_received_history` VALUES "):
			tablename = "user_present_all_received_history"
		default:
			for _, w := range writers {
				if _, err := w.WriteString(line); err != nil {
					panic(err)
				}
			}
			continue
		}

		inserts := make([]string, cnt)
		prefix := "INSERT INTO `" + tablename + "` VALUES ("
		for i := range inserts {
			inserts[i] = prefix
		}
		values := strings.TrimPrefix(line, prefix)
		values = strings.TrimSuffix(values, ");\n")
		records := strings.Split(values, "),(")
		for _, record := range records {
			cols := strings.Split(record, ",")
			userIDStr := cols[1]
			if tablename == "users" {
				userIDStr = cols[0]
			}
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				fmt.Println("cols:", cols)
				panic(err)
			}
			shardNum := int(userID) % cnt
			if inserts[shardNum] != prefix {
				inserts[shardNum] += "),("
			}
			inserts[shardNum] += record
		}
		for i := range inserts {
			inserts[i] += ");"
			if _, err := writers[i].WriteString(inserts[i] + "\n"); err != nil {
				panic(err)
			}
		}
	}
}
