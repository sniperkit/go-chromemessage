package chromemsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"

	// external
	"github.com/k0kubun/pp"

	// internal
	"github.com/sniperkit/snk.golang.chromemsg/pkg/version"
)

var nativeEndian binary.ByteOrder = endianness()

var defaultMsgr = Messenger{
	port: bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)),
	lock: &sync.RWMutex{},
	rev:  version.Version,
}

type Config struct {
	PortType PortType
	Debug    bool
	Verbose  bool
}

type Messenger struct {
	rev  string
	port *bufio.ReadWriter
	lock *sync.RWMutex
	conf *Config
}

func New(port *bufio.ReadWriter) *Messenger {
	if port == nil {
		port = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	}
	return &Messenger{
		rev:  version.Version,
		port: port,
		lock: &sync.RWMutex{},
	}
}

func NewWithConfig(conf *Config) (*Messenger, error) {
	m := &Messenger{
		rev:  version.Version,
		lock: &sync.RWMutex{},
	}

	if conf == nil {
		conf = &Config{}
	}

	if conf.Debug {
		pp.Println("conf=", conf)
	}

	switch conf.PortType {
	case Bufio:
		fmt.Println("creating new port wrapper wuth bufio...")
		m.port = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))

	default:
		return nil, errors.New(
			fmt.Sprintf("Unkown port wrapper... conf.PortType=%s", conf.PortType),
		)

	}

	return m, nil
}

func Read(data interface{}) error {
	return defaultMsgr.Read(data)
}

func Write(msg interface{}) error {
	return defaultMsgr.Write(msg)
}

func (m *Messenger) Config(conf *Config) error {
	m.lock.RLock()
	m.conf = conf

	if conf.PortType == "buffio" {
		m.port = bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	}

	m.lock.RUnlock()
	return nil
}

func (m *Messenger) Read(data interface{}) error {
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

func (m *Messenger) Write(msg interface{}) error {
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
