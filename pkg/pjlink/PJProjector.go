package pjlink

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

const pjLinkPort = "4352"
const ON = "on"
const OFF = "off"
const ERRA = "ERRA"
const OK = "OK"

type Projector struct {
	Address  string
	Port     string
	Password string
}

func NewProjector(ip, port, password string) *Projector {
	if port == "0" || port == "" {
		port = pjLinkPort
	}
	return &Projector{
		Address:  ip,
		Port:     port,
		Password: password,
	}
}

//--------------------------------------------------------------------------------------------------------------------//
//--------------- Functional Calls -----------------------------------------------------------------------------------//
//--------------------------------------------------------------------------------------------------------------------//

// --------------- Power ----------------------------------------------------------------------------------------------//
func (pr *Projector) GetPowerStatus() (*Response, error) {
	req := Request{
		Class:     1,
		Command:   "POWR",
		Parameter: "?",
	}
	return pr.SendRequest(req)
}

func (pr *Projector) PowerOn() error {
	req := Request{
		Class:     1,
		Command:   "POWR",
		Parameter: "1",
	}
	resp, err := pr.SendRequest(req)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return errors.New("could not turn on Projector")
}

func (pr *Projector) PowerOff() error {
	req := Request{
		Class:     1,
		Command:   "POWR",
		Parameter: "0",
	}
	resp, err := pr.SendRequest(req)
	if err != nil {
		return err
	}
	if resp.Success() {
		return nil
	}
	return errors.New("could not turn off Projector")
}

func (p *Projector) GetProperty(property string) (string, error) {
	var request Request

	request.Class = 1
	request.Command = property
	request.Parameter = "?"

	resp, err := p.SendRequest(request)

	if err != nil {
		return "", err
	}

	/*	log.Printf("response size for %s: %d\n", property, len(resp.Response))
		for i := 0; i < len(resp.Response); i++ {
			log.Printf("response %d: %s\n", i, resp.Response[i])
		} */

	return resp.Response[0], nil
}

func (p *Projector) GetPropertyArray(property string) ([]string, error) {
	var request Request

	request.Class = 1
	request.Command = property
	request.Parameter = "?"

	resp, err := p.SendRequest(request)

	if err != nil {
		return make([]string, 0), err
	}

	return resp.Response, nil
}

func (p *Projector) SetProperty(property string, val string) error {
	var request Request

	request.Class = 1
	request.Command = property
	request.Parameter = val

	_, err := p.SendRequest(request)

	return err
}

// --------------------------------------------------------------------------------------------------------------------//
// Low-Level Calls
// --------------------------------------------------------------------------------------------------------------------//
func (pr *Projector) SendRequest(request Request) (*Response, error) {
	if err := request.Validate(); err != nil { //malformed command, don't send
		return nil, err
	} else { //send request and parse response into struct
		response, requestError := pr.sendRawRequest(request)
		if requestError != nil {
			return nil, requestError
		}

		return response, nil
	}
}

func (pr *Projector) sendRawRequest(request Request) (*Response, error) {
	//establish TCP connection with PJLink device
	connection, connectionError := pr.connectToPJLink()
	defer func() {
		if connection == nil {
			return
		}

		connection.Close()
	}()

	if connectionError != nil {
		return nil, connectionError
	}

	// Define a split function that separates on carriage return (i.e '\r').
	onCarriageReturn := func(data []byte, atEOF bool) (advance int, token []byte,
		err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == '\r' {
				return i + 1, data[:i], nil
			}
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens
		// after this but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}

	//setup scanner
	scanner := bufio.NewScanner(connection)
	scanner.Split(onCarriageReturn)
	scanner.Scan() //grab a line
	challenge := strings.Split(scanner.Text(), " ")

	seed := pr.checkAuthentication(challenge)
	stringCommand := request.toRaw(seed, pr.Password)

	//send command
	connection.Write([]byte(stringCommand))
	scanner.Scan() //grab response line

	resp := NewPJResponse()
	err := resp.Parse(scanner.Text())
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// attempts to establish a TCP socket with the specified IP:port
// success: returns populated pjlinkConn struct and nil error
// failure: returns empty pjlinkConn and error
func (pr *Projector) connectToPJLink() (net.Conn, error) {
	protocol := "tcp" //PJLink always uses TCP
	timeout := 10     //represents seconds

	connection, connectionError := net.DialTimeout(protocol, net.JoinHostPort(pr.Address, pr.Port), time.Duration(timeout)*time.Second)
	if connectionError != nil {
		return nil, fmt.Errorf("failed to establish a connection with device @ %s. error msg: %s", pr.Address, connectionError.Error())
	}

	return connection, connectionError
}

// check if this Projector uses authentication. If so return true and the given seed. Otherwise false and an empty string.
func (pr *Projector) checkAuthentication(response []string) (seed string) {
	if response[0] != "PJLINK" {
		return ""
	}
	if response[1] == "0" {
		return ""
	} else if response[1] == "1" {
		return response[2]
	}

	return ""
}
