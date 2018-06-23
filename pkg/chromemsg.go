package chromemsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"unsafe"

	// external
	"github.com/k0kubun/pp"

	// internal
	"github.com/sniperkit/snk.golang.chromemsg/pkg/version"
)

var nativeEndian binary.ByteOrder = endianness()

var defaultMsgr = messenger{
	port: bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)),
	lock: &sync.RWMutex{},
	rev:  version.Version,
}

type Config struct {
	PortType PortType
	Debug    bool
	Verbose  bool
}

type messenger struct {
	rev  string
	port *bufio.ReadWriter
	lock *sync.RWMutex
	conf *Config
}

func New(port *bufio.ReadWriter) *messenger {
	if port == nil {
		port = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	}
	return &messenger{
		rev:  version.Version,
		port: port,
		lock: &sync.RWMutex{},
	}
}

func NewWithConfig(conf *Config) (*messenger, error) {

	m := &messenger{
		rev:  version.Version,
		lock: &sync.RWMutex{},
	}

	if conf == nil {
		conf = &Config{}
	}

	if conf.Debug {
		pp.Println("conf=", conf)
	}

	if conf.PortType == "buffio" {
		m.port = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	}

	if m.port == nil {
		return nil, errors.New("No port defined...")
	}

	return m, nil
}

func Read(data interface{}) error {
	return defaultMsgr.Read(data)
}

func Write(msg interface{}) error {
	return defaultMsgr.Write(msg)
}

func (m *messenger) Config(conf *Config) error {
	m.lock.RLock()
	m.conf = conf

	if conf.PortType == "buffio" {
		m.port = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	}

	m.lock.RUnlock()
	return nil
}

func (m *messenger) Read(data interface{}) error {
	m.lock.RLock()

	lengthBits := make([]byte, 4)
	_, err := m.port.Read(lengthBits)
	if err != nil {
		return err
	}

	length := nativeToInt(lengthBits)
	content := make([]byte, length)
	_, err = m.port.Read(content)
	if err != nil {
		return err
	}
	json.Unmarshal(content, data)

	m.lock.RUnlock()
	return nil
}

func (m *messenger) Write(msg interface{}) error {
	m.lock.Lock()

	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	length := len(json)
	bits := make([]byte, 4)
	buf := bytes.NewBuffer(bits)
	err = binary.Write(buf, nativeEndian, length)
	if err != nil {
		return err
	}
	_, err = m.port.Write(bits)
	if err != nil {
		return err
	}

	_, err = m.port.Write(json)
	if err != nil {
		return err
	}
	m.lock.Unlock()
	return nil
}

func nativeToInt(bits []byte) int {
	var length uint32
	buf := bytes.NewBuffer(bits)
	binary.Read(buf, nativeEndian, &length)
	return int(length)
}

func endianness() binary.ByteOrder {
	var i int = 1
	bs := (*[unsafe.Sizeof(0)]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
