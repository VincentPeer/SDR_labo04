package udpserver

import (
	"net"
	"strconv"
)

type UDPAddress struct {
	// The address to listen on
	Address string
	// The port to listen on
	Port string
}

// UDP is a UDP server
type UDP struct {
	// UDPAddress is the connection to the UDP server
	UDPAddress UDPAddress
	// The id of the server
	ID string
	// The connection to the UDP server
	conn *net.UDPConn
	// Whether the server has been started
	Started bool
	done chan struct{}
}

func NewUDPConn(address string, port int) UDPAddress {
	return UDPAddress{
		Address: address,
		Port:    strconv.Itoa(port),
	}
}

// NewUDP creates a new UDP server
func NewUDP(address string, port int, id string) *UDP {
	return &UDP{
		UDPAddress: NewUDPConn(address, port),
		ID:         id,
		Started:    false,
	}
}

// Start initializes the UDP server
func (u *UDP) Start() error {
	addr, err := net.ResolveUDPAddr("udp", u.UDPAddress.Address+":"+u.UDPAddress.Port)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	u.conn = conn
	u.done = make(chan struct{})
	u.Started = true

	return nil
}

// stop closes the UDP server
func (u *UDP) stop() error {
	if !u.Started {
		return nil
	}
	u.Started = false
	u.done <- struct{}{}
	return u.conn.Close()
}

// Listen listens for incoming UDP packets
func (u *UDP) Listen(onMessage func(message string, remoteAddr *UDPAddress), onError func(err error)) {
	if !u.Started {
		err := u.Start()
		if err != nil {
			onError(err)
			return
		}
	}

	defer u.stop()

	buffer := make([]byte, 1024)


	go func() {
		for {
			readLength, remoteAddr, err := u.conn.ReadFromUDP(buffer)
			if err != nil {
				onError(err)
				return
			}
	
			remoteUDP := NewUDPConn(remoteAddr.IP.String(), remoteAddr.Port)
			go onMessage(string(buffer[:readLength]), &remoteUDP)
		}
	}()
	
	for {
		select {
		case <-u.done:
			return
		default:
			if !u.Started {
				close(u.done)
			}
		}
	}	
}

// Send sends a UDP packet to the specified address
func (u *UDP) Send(remoteUDP *UDPAddress, message string) error {
	if !u.Started {
		err := u.Start()
		if err != nil {
			return err
		}
	}

	remoteAddr, err := net.ResolveUDPAddr("udp", remoteUDP.Address+":"+remoteUDP.Port)
	if err != nil {
		return err
	}

	_, err = u.conn.WriteToUDP([]byte(message), remoteAddr)
	if err != nil {
		return err
	}

	return nil
}
