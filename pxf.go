package main

import (
	"fmt"
	"image"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

var maxY atomic.Int32

func main() {
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			n := maxY.Add(3)
			if n > 700 {
				maxY.Store(400)
			}
		}
	}()

	f, err := os.Open("img.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	bounds := image.Bounds()
	size := bounds.Size()

	fmt.Printf("size: %v\n", size)

	const sizeX, sizeY = 1920, 1080
	/*
					var pixels [sizeX][sizeY]string
					for x := 0; x < sizeX; x++ {
						for y := 0; y < sizeY; y++ {
							ix := x % (size.X)
							iy := y % (size.Y)

							if ix < 3 || iy < 3 {
								//pixels[x][y] = "PX 0 0 ff00ff\n"
								pixels[x][y] = fmt.Sprintf("PX %d %d ff00ff\n", x, y)
								continue
							}

							ix -= 3
							iy -= 3

							col := image.At(ix, iy)
							R, G, B, A := col.RGBA()
							if A == 0 {
								continue
							}

				ip link add link ens0s1 ens0s1:0 type macvlan

			ip link add mv1 link enp0s1 type macvlan mode bridge
			ip link add mv2 link enp0s1 type macvlan mode bridge
			ip link add mv3 link enp0s1 type macvlan mode bridge
			ip link add mv4 link enp0s1 type macvlan mode bridge
			ip link add mv5 link enp0s1 type macvlan mode bridge
			ip link add mv6 link enp0s1 type macvlan mode bridge
			ip link add mv7 link enp0s1 type macvlan mode bridge
			ip link add mv8 link enp0s1 type macvlan mode bridge
			ip link add mv9 link enp0s1 type macvlan mode bridge
			ip link add mv10 link enp0s1 type macvlan mode bridge
			ip link add mv11 link enp0s1 type macvlan mode bridge
			ip link add mv12 link enp0s1 type macvlan mode bridge
			ip link add mv13 link enp0s1 type macvlan mode bridge
			ip link add mv14 link enp0s1 type macvlan mode bridge
			ip link add mv15 link enp0s1 type macvlan mode bridge
			ip link add mv16 link enp0s1 type macvlan mode bridge
			ip link add mv17 link enp0s1 type macvlan mode bridge
			ip link add mv18 link enp0s1 type macvlan mode bridge
			ip link add mv19 link enp0s1 type macvlan mode bridge
			ip link add mv20 link enp0s1 type macvlan mode bridge
			ip link add mv21 link enp0s1 type macvlan mode bridge
			ip link add mv22 link enp0s1 type macvlan mode bridge
			ip link add mv23 link enp0s1 type macvlan mode bridge
			ip link add mv24 link enp0s1 type macvlan mode bridge
			ip link add mv25 link enp0s1 type macvlan mode bridge
			ip link add mv26 link enp0s1 type macvlan mode bridge
			ip link add mv27 link enp0s1 type macvlan mode bridge
			ip link add mv28 link enp0s1 type macvlan mode bridge
			ip link add mv29 link enp0s1 type macvlan mode bridge
			ip link add mv30 link enp0s1 type macvlan mode bridge
			ip link add mv31 link enp0s1 type macvlan mode bridge
			ip link add mv32 link enp0s1 type macvlan mode bridge
			ip link add mv33 link enp0s1 type macvlan mode bridge
			ip link add mv34 link enp0s1 type macvlan mode bridge
			ip link add mv35 link enp0s1 type macvlan mode bridge
			ip link add mv36 link enp0s1 type macvlan mode bridge
			ip link add mv37 link enp0s1 type macvlan mode bridge
			ip link add mv38 link enp0s1 type macvlan mode bridge
			ip link add mv39 link enp0s1 type macvlan mode bridge
			ip link add mv40 link enp0s1 type macvlan mode bridge

		dhcpcd mv1 &
		dhcpcd mv2 &
		dhcpcd mv3 &
		dhcpcd mv4 &
		dhcpcd mv5 &
		dhcpcd mv6 &
		dhcpcd mv7 &
		dhcpcd mv8 &
		dhcpcd mv9 &
		dhcpcd mv10 &
		dhcpcd mv11 &
		dhcpcd mv12 &
		dhcpcd mv13 &
		dhcpcd mv14 &
		dhcpcd mv15 &
		dhcpcd mv16 &
		dhcpcd mv17 &
		dhcpcd mv18 &
		dhcpcd mv19 &
		dhcpcd mv20 &
		dhcpcd mv21 &
		dhcpcd mv22 &
		dhcpcd mv23 &
		dhcpcd mv24 &
		dhcpcd mv25 &
		dhcpcd mv26 &
		dhcpcd mv27 &
		dhcpcd mv28 &
		dhcpcd mv29 &
		dhcpcd mv30 &
		dhcpcd mv31 &
		dhcpcd mv32 &
		dhcpcd mv33 &
		dhcpcd mv34 &
		dhcpcd mv35 &
		dhcpcd mv36 &
		dhcpcd mv37 &
		dhcpcd mv38 &
		dhcpcd mv39 &
		dhcpcd mv40 &





								pixels[x][y] = fmt.Sprintf("PX %d %d %02x%02x%02x\n", x, y, R>>8, G>>8, B>>8)
						}
					}*/

	nwork := 40
	nconn := 1

	sources := `151.217.1.242
151.217.1.243
151.217.1.244
151.217.1.245
151.217.1.221
151.217.1.219
151.217.1.218
151.217.1.207
151.217.1.214
151.217.1.246
151.217.1.249
151.217.1.212
151.217.1.226
151.217.1.248
151.217.1.204
151.217.1.232
151.217.1.247
151.217.1.222
151.217.1.228
151.217.1.235
151.217.1.208
151.217.1.234
151.217.1.229
151.217.1.206`

	sourceIPs := strings.Split(sources, "\n")

	nconn *= len(sourceIPs)

	var connWg sync.WaitGroup

	connWg.Add(nwork)

	var conns []net.Conn
	var connLk sync.Mutex
	connCond := sync.NewCond(&connLk)

	for i := 0; i < nconn; i++ {
		/*conn, err := net.Dial("tcp", "151.217.15.79:1337")
		if err != nil {
			panic(err)
		}*/

		src := sourceIPs[i%len(sourceIPs)]
		// trim
		src = strings.TrimSpace(src)

		dst := "151.217.15.79"

		conn, err := net.DialTCP("tcp", &net.TCPAddr{IP: net.ParseIP(src)}, &net.TCPAddr{IP: net.ParseIP(dst), Port: 1337})
		if err != nil {
			panic(err)
		}

		if err := conn.SetNoDelay(true); err != nil {
			panic(err)
		}

		file, err := conn.File()
		if err != nil {
			fmt.Println("Error retrieving file descriptor:", err)
			os.Exit(1)
		}
		//defer file.Close()

		fd := int(file.Fd())

		// Set the TOS field
		// For example, set TOS to 0x28, which is a common value for AF41 (Assured Forwarding)
		err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TOS, 0x28)
		if err != nil {
			fmt.Println("Error setting TOS:", err)
			os.Exit(1)
		}

		conns = append(conns, conn)
	}

	for wid := 0; wid < nwork; wid++ {
		go func(wid int) {
			//offy := wid
			fmt.Printf("[%d]", wid)

			connWg.Done()
			connWg.Wait()

			sendBuf := make([]byte, 0, 1900*300*8)

			/*for i := 0; i < 1900*300; i++ { // 1900x900
				x := i % 1900
				y := offy*300 + i/1900

				x = 1920 - x - 1
				y = 1080 - y - 1
				/*_, err := conn.Write([]byte(pixels[x][y]))
				if err != nil {
					panic(err)
				}* /

				// trim to image size
				if x >= 2*size.X || y >= 3*size.Y {
					continue
				}
				if x < size.X || y < size.Y*2 {
					continue
				}

				sendBuf = append(sendBuf, pixels[x][y]...)
			}*/

			// rand startX, startY

			for {
				sendBuf = sendBuf[:0]
				iters := 5

				max := int(maxY.Load()) + 1

				for i := 0; i < iters; i++ {
					startX := rand.Intn(1800)

					startY := rand.Intn(80) + max

					for i := 0; i < size.X*size.Y; i++ { // 1900x900
						x := i % size.X
						y := i / size.X

						if x*y*127%13 > 9 {
							continue
						}

						px := image.At(x, y)
						R, G, B, A := px.RGBA()

						if A == 0 {
							continue
						}

						sendBuf = append(sendBuf, fmt.Sprintf("PX %d %d %02x%02x%02x\n", startX+x, startY+y, R>>8, G>>8, B>>8)...)
					}
				}

				// get a conn
				connLk.Lock()
				for len(conns) == 0 {
					connCond.Wait()
				}
				conn := conns[0]
				conns = conns[1:]
				connLk.Unlock()

				_, err := conn.Write(sendBuf)
				if err != nil {
					panic(err)
				}

				connLk.Lock()
				conns = append(conns, conn)
				connLk.Unlock()
				connCond.Signal()
			}
		}(wid)
	}

	select {}
}
