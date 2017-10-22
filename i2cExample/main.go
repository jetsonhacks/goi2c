package main

import (
	"fmt"
	"time"

	"github.com/jetsonhacks/goi2c/devices/ledBackpack7Segment"
)

func showCountdown(backpack *ledBackpack7Segment.LedBackpack7Segment, counter int) {
	// Countdown from 10 to 0
	ticker := time.NewTicker(time.Second)
	countDownOver := false
	backpack.WriteString("  10")
	counter--
	for {
		if countDownOver == true {
			break
		}
		<-ticker.C
		countString := fmt.Sprintf("%4d", counter)
		backpack.WriteString(countString)
		counter--
		if counter == -1 {
			ticker.Stop()
			countDownOver = true
		}
	}
	// Flash the final 0 for a couple of seconds
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_2hz)
	time.Sleep(2 * time.Second)
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_off)
}

func showClock(backpack *ledBackpack7Segment.LedBackpack7Segment) {
	counter := 590
	for {
		hours := counter / 60
		minutes := counter % 60
		counter++
		clockString := fmt.Sprintf("%2d:%02d", hours, minutes)
		backpack.WriteString(clockString)
		time.Sleep(time.Second)
		if counter > 610 {
			return
		}
	}
}

func demoSeparator(backpack *ledBackpack7Segment.LedBackpack7Segment) {
	backpack.WriteError()
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_2hz)
	time.Sleep(2 * time.Second)
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_off)
}

func main() {
	fmt.Println("Entered Main")
	backpack, err := ledBackpack7Segment.NewLedBackpack7Segment(1, 0x70)
	fmt.Println("backpack: ", backpack)
	if err != nil {
		panic(err)
	}
	defer backpack.Close()

	backpack.Begin()
	defer backpack.End()
	backpack.WriteString("    ")
	backpack.WriteString("1.234")
	time.Sleep(3 * time.Second)
	demoSeparator(backpack)
	// Should ignore last invalid character
	backpack.WriteString("9.8.:7.6.5.")
	time.Sleep(3 * time.Second)
	demoSeparator(backpack)

	backpack.WriteString("12:15")
	time.Sleep(3 * time.Second)
	demoSeparator(backpack)
	// Countdown from 10
	showCountdown(backpack, 10)
	demoSeparator(backpack)
	showClock(backpack)
	demoSeparator(backpack)

	fmt.Println("Leaving Main")
}
