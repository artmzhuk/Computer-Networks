// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sooockets"
	"strconv"
	"strings"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		str := ""
		fmt.Scanln(&str)
		args := strings.Split(str, ";")
		a, _ := strconv.ParseFloat(args[0], 64)
		b, _ := strconv.ParseFloat(args[1], 64)
		cI, _ := strconv.ParseFloat(args[2], 64)
		d, _ := strconv.ParseFloat(args[3], 64)
		alpha, _ := strconv.ParseFloat(args[4], 64)
		beta, _ := strconv.ParseFloat(args[5], 64)
		integral := sooockets.Integral1{
			A:     a,
			B:     b,
			C:     cI,
			D:     d,
			Alpha: alpha,
			Beta:  beta,
		}
		fmt.Println(integral)
		c.WriteJSON(&integral)
	}

}
