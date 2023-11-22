// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sooockets"
	"strconv"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		integral := sooockets.Integral1{
			A:     0,
			B:     0,
			C:     0,
			D:     0,
			Alpha: 0,
			Beta:  0,
		}
		err := c.ReadJSON(&integral)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: ", integral)
		res := computeIntegral(integral)
		err = c.WriteMessage(websocket.TextMessage, []byte(res))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func computeIntegral(integral sooockets.Integral1) string {
	a := integral.A
	b := integral.B
	c := integral.C
	d := integral.D
	alpha := integral.Alpha
	beta := integral.Beta
	res := (a * beta * beta * beta * beta) / 4
	res += (b * beta * beta * beta) / 3
	res += (c * beta * beta) / 2
	res += d * beta
	res -= (a * alpha * alpha * alpha * alpha) / 4
	res -= (b * alpha * alpha * alpha) / 3
	res -= (c * alpha * alpha) / 2
	res -= d * alpha
	resStr := strconv.FormatFloat(res, 'g', -1, 64)
	return resStr

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
