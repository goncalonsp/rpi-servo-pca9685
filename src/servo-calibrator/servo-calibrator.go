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

	"github.com/golang/glog"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	"github.com/kidoman/embd/motion/servo"

	_ "github.com/kidoman/embd/host/rpi"
)

func listenForEntry(anglesChan chan int) {
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
	defer glog.Flush()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	d := pca9685.New(bus, 0x40)
	d.Freq = 50
	defer d.Close()

	pwm0 := d.ServoChannel(0)
	servo0 := servo.New(pwm0)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	pulseChan := make(chan int)
	go listenForEntry(pulseChan)

	defer func() {
		servo0.SetAngle(1000)
	}()

	fmt.Print("Enter values to experiment, try slow increments:\n")

	for {
		select {
		case pulse := <-pulseChan:
			fmt.Printf("> Setting to %d microseconds\n", pulse)
			servo0.PWM.SetMicroseconds(int(pulse))
		case <-c:
			return
		}
	}
}
