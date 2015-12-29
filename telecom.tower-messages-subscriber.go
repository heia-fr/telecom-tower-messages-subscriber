// Copyright 2015 Jacques Supcik, HEIA-FR
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// 2015-12-29 | JS | First version

//
// Telecom Tower Messages Subscriber
//

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/nats-io/nats"
	"github.com/vharitonsky/iniflags"
	"log"
	"net/http"
)

type Line struct {
	Text  string `json:"text"`
	Font  int    `json:"font"`
	Color string `json:"color"`
}

type RollingMessage struct {
	Body         []Line `json:"body"`
	Introduction []Line `json:"introduction"`
	Conclusion   []Line `json:"conclusion"`
	Separator    []Line `json:"separator"`
}

func main() {
	var natsUrl = flag.String("nats-url", nats.DefaultURL, "NATS URL")
	var natsSubject = flag.String("nats-subject", "heiafr.telecomtower.bot", "NATS Subject")
	var towerUrl = flag.String("tower-url", "http://localhost:8080/message", "Tower URL")

	iniflags.Parse()

	nc, err := nats.Connect(*natsUrl)
	if err != nil {
		log.Fatal(err)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	recvCh := make(chan *RollingMessage)
	ec.BindRecvChan(*natsSubject, recvCh)
	for message := range recvCh {
		if b, err := json.Marshal(message); err != nil {
			log.Println(err)
		} else {
			if _, err = http.Post(
				*towerUrl,
				"application/json; charset=utf-8", bytes.NewBuffer(b)); err != nil {
				log.Println(err)
			}

		}
	}

}
