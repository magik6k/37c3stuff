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
			n := maxY.Add(5)
			if n > 1000 {
				maxY.Store(10)
			}
		}
	}()

	var images []image.Image
	var sizes []image.Point

	loadimg := func(p string) {
		f, err := os.Open(p)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		im, _, err := image.Decode(f)
		if err != nil {
			panic(err)
		}

		bounds := im.Bounds()
		size := bounds.Size()

		images = append(images, im)
		sizes = append(sizes, size)
	}

	loadimg("img.png")
	loadimg("img1.png")
	loadimg("img2.png")
	loadimg("img3.png")
	loadimg("img4.png")

	nwork := 40
	nconn := 3

	sources := `151.217.1.242
151.217.1.243
151.217.1.244
151.217.1.247
151.217.1.221
151.217.1.208
151.217.1.218
151.217.1.217
151.217.1.214
151.217.1.246
151.217.1.249
151.217.1.212
151.217.1.226
151.217.1.248
151.217.1.204
151.217.1.232
151.217.1.231
151.217.1.227
151.217.1.228
151.217.1.235
151.217.1.234
151.217.1.207
151.217.1.206
151.217.1.229
151.217.1.224
151.217.1.245
151.217.1.237`

	sourceIPs := strings.Split(sources, "\n")

	nconn *= len(sourceIPs)

	var connWg sync.WaitGroup

	connWg.Add(nwork)

	var conns []net.Conn
	var connLk sync.Mutex
	connCond := sync.NewCond(&connLk)

	go func() {

		for i := 0; i < nconn; i++ {

			src := sourceIPs[i%len(sourceIPs)]
			// trim
			src = strings.TrimSpace(src)

			dst := "151.217.15.90"

			dialer := net.Dialer{
				Timeout:   170 * time.Millisecond, // Set the connection timeout to 70 milliseconds
				LocalAddr: &net.TCPAddr{IP: net.ParseIP(src)},
			}

		retry:
			conn, err := dialer.Dial("tcp", net.JoinHostPort(dst, "1337"))
			if err != nil {
				time.Sleep(1 * time.Millisecond)
				fmt.Printf("retrying %s: %s\n", src, err)
				goto retry
			}

			fmt.Println("connected to", conn.RemoteAddr())

			if err := conn.(*net.TCPConn).SetNoDelay(true); err != nil {
				panic(err)
			}

			file, err := conn.(*net.TCPConn).File()
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

			connLk.Lock()
			conns = append(conns, conn)
			connLk.Unlock()
			connCond.Signal()
		}
	}()

	for wid := 0; wid < nwork; wid++ {
		go func(wid int) {
			//offy := wid
			fmt.Printf("[%d]", wid)

			connWg.Done()
			connWg.Wait()

			sendBuf := make([]byte, 0, 1900*300*8)

			// rand startX, startY

			for {
				sendBuf = sendBuf[:0]
				iters := 5

				max := int(maxY.Load()) + 1

				for i := 0; i < iters; i++ { // prepare iter buffers
					startX := rand.Intn(1800)

					startY := rand.Intn(80) + max

					imidx := rand.Intn(len(images))
					image := images[imidx]
					size := sizes[imidx]

					for i := 0; i < size.X*size.Y; i++ { // 1900x900
						x := i % size.X
						y := i / size.X

						if x*y*127%13 > 7 {
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
					// bad conn
					go func() {
						conn.Close()
						la := conn.(*net.TCPConn).LocalAddr()
						ra := conn.(*net.TCPConn).RemoteAddr()

						dialer := net.Dialer{
							Timeout:   170 * time.Millisecond, // Set the connection timeout to 70 milliseconds
							LocalAddr: &net.TCPAddr{IP: net.ParseIP(la.String())},
						}

					retry:
						nconn, err := dialer.Dial("tcp", ra.String())
						if err != nil {
							time.Sleep(1 * time.Millisecond)
							goto retry
						}

						fmt.Println("reconnected to", nconn.RemoteAddr())

						if err := nconn.(*net.TCPConn).SetNoDelay(true); err != nil {
							panic(err)
						}

						file, err := nconn.(*net.TCPConn).File()
						if err != nil {
							fmt.Println("Error retrieving file descriptor:", err)
							os.Exit(1)
						}

						fd := int(file.Fd())

						// Set the TOS field
						// For example, set TOS to 0x28, which is a common value for AF41 (Assured Forwarding)
						err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TOS, 0x28)
						if err != nil {
							fmt.Println("Error setting TOS:", err)
							os.Exit(1)
						}

						connLk.Lock()
						conns = append(conns, nconn)
						connLk.Unlock()
						connCond.Signal()
					}()
					continue
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
