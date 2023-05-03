package main

import (
	"io"
	"log"
	"net"
	"sync"
	"temp/conv"
	"strings"
	"strconv"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8080")
	if err != nil {
		log.Fatal(err)
	}

	a := conv.CelsiusToFahrenheit(3)
	log.Println(a)

	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					switch msg := string(buf[:n]); msg {

  				        case "ping":
						_, err = c.Write([]byte("pong"))

					case "Kjevik;SN39040;18.03.2022 01:50;6":

						a1 := strings.Split((msg), ";")

						a2, err := strconv.ParseFloat(a1[3], 64)

						a3 := conv.CelsiusToFahrenheit(a2)

						a4 := strconv.FormatFloat(a3, 'f', -1, 64)
										
						_, err = c.Write([]byte(a1[0] + ";" + a1[1] + ";" + a1[2] + ";" + a4))

						if err != nil {
							panic(err)
						}

					default:
						_, err = c.Write(buf[:n])
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
