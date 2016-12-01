package main

import (
	"flag"
	"os"
	"os/signal"
	//"time"
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	"github.com/kidoman/embd/motion/servo"

	_ "github.com/kidoman/embd/host/rpi"
)

func listenForAngles(anglesChan chan int) {
	inputReader := bufio.NewReader(os.Stdin)
    for {
        s, err := inputReader.ReadString('\n')
        if err != nil {
            fmt.Printf("Read error: %v", err)
            return
        }
        i, _ := strconv.Atoi(strings.Trim(s, "\n\t "))
        anglesChan <- i
    }
    close(anglesChan)
}

func main() {
	flag.Parse()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	d := pca9685.New(bus, 0x40)
	d.Freq = 50
	defer d.Close()

	pwm0 := d.ServoChannel(0)
	pwm1 := d.ServoChannel(1)
	servo0 := servo.New(pwm0)
	servo1 := servo.New(pwm1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	anglesChan := make(chan int)
	go listenForAngles(anglesChan)

	defer func() {
		servo0.SetAngle(90)
		servo1.SetAngle(90)
	}()

	for {
		select {
		case angle := <- anglesChan:
			fmt.Printf("> Setting angle to %d degrees\n", angle)
			servo0.SetAngle(angle)
			servo1.SetAngle(angle)
		case <-c:
			return
		}
	}
}