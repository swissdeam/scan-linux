package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/bendahl/uinput"
	"go.bug.st/serial"
	"golang.design/x/clipboard"
)

func main() {
	clipboard.Init()

	ports := connection()

	// Список обнаруженных портов
	for _, port := range ports {
		fmt.Printf("Найден Com-порт: %v\n", port)
	}

	// Открывает первый порт
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	for {
		if len(ports) == 0 {
			log.Println("Подключите Reader/scaner")
			ports = connection()
			time.Sleep(3 * time.Second)
		} else {
			log.Println("Успешно подключено")
			break
		}
	}

	port, err := serial.Open(ports[0], mode)
	if err != nil {
		log.Println(err)
	}

	buff := make([]byte, 100)
	var code string

	for {
		read(port, buff, code)
	}

}

func connection() []string {

	// Retrieve the port list
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Println(err)
	}
	if len(ports) == 0 {
		log.Printf("Com-портов не найдено!")

	}

	return ports

}

func read(p serial.Port, b []byte, c string) {

	for {
		// Reads up to 100 bytes
		n, err := p.Read(b)
		if err != nil {
			log.Println(err)
		}
		if n == 0 {
			fmt.Println("EOF")

		}

		// fmt.Printf("%s", string(b[:n]))
		// ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		// defer cancel()

		c = c + string((b[:n]))

		// If we receive a newline stop reading
		if strings.Contains(string(b[:n]), "\n") {

			// fmt.Println(i, "после чтения")
			// fmt.Println(temp, "прошлая")
			// fmt.Println(c, "текущая")

			patternNoCode := "No card"
			matched, _ := regexp.MatchString(patternNoCode, c)

			if matched {
				return
			} else {

				fmt.Println(c)

				pattern1, _ := regexp.Compile(`\[[a-zA-Z0-9]+\]`)
				codex1 := pattern1.FindAllString(c, -1)
				// fmt.Println(codex1)

				pattern2, _ := regexp.Compile(`[a-zA-Z0-9]+`)
				codex2 := pattern2.FindAllString(codex1[0], -1)
				// fmt.Println(codex2)

				clipboard.Write(clipboard.FmtText, []byte(codex2[0]))

				// initialize keyboard and check for possible errors
				keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
				if err != nil {
					return
				}

				if runtime.GOOS == "linux" {
					time.Sleep(2 * time.Second)
				}

				defer keyboard.Close()

				keyboard.KeyDown(uinput.KeyLeftctrl)
				keyboard.KeyPress(uinput.KeyV)
				keyboard.KeyUp(uinput.KeyLeftctrl)

				// fmt.Println(c, "Check response")
				// fmt.Println(codex1, "Check response x1")
				// fmt.Println(codex2, "Check response x2")
			}

		}

	}

}
