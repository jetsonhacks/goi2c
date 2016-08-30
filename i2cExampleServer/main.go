// I2C example for Adafruit 7 Segment Display with I2C Backpack with Server Side Events 
// URL: localhost:8000/test
package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"github.com/jetsonhacks/goi2c/devices/ledBackpack7Segment"
	"log"
	"net/http"
)


// Write a string to the 7 segment display and also send it out as a Server Side Event 
func writeString ( backpack *ledBackpack7Segment.LedBackpack7Segment, stringToWrite string, messageChannel chan <- string ) {
	backpack.WriteString(stringToWrite)
	messageChannel <- stringToWrite
	log.Printf("Sent message: %s", stringToWrite)        
}

//
// Example of a countdown timer
//
func showCountdown(backpack *ledBackpack7Segment.LedBackpack7Segment, counter int, messageChannel chan <- string)() {
	// Countdown from 10 to 0
	ticker := time.NewTicker(time.Second)
	countDownOver := false 
	writeString(backpack,"  10",messageChannel)
	counter --
	for {	
		if countDownOver == true {	
			break 
		}
		<- ticker.C
		countString := fmt.Sprintf("%4d",counter)
		writeString(backpack,countString,messageChannel)
		counter --
		if counter == -1 {
			ticker.Stop()
			countDownOver = true
		}
	}
	// Flash the final 0 for a couple of seconds
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_2hz)
	time.Sleep(2*time.Second)
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_off)
}

//
// An Example of a clock showing seconds
// Clock runs from 9:50 to 10:10
//
func showClock (backpack *ledBackpack7Segment.LedBackpack7Segment, messageChannel chan <- string) () {
	counter := 590
	for {	
		hours := counter/60
		minutes := counter%60
		counter ++
		clockString := fmt.Sprintf("%2d:%02d",hours,minutes)
		writeString(backpack,clockString,messageChannel)
		time.Sleep(time.Second)
		if counter > 610 {
			return 
		}
	}
}

//
// Draw ---- on the 7 Segment Display
// And send the '----' as a Server Side Event
// Note: The 7 Segment Display blinks, but not the corresponding web page
//
func demoSeparator (backpack *ledBackpack7Segment.LedBackpack7Segment, messageChannel chan <- string) () {
	// backpack.WriteError()
	writeString(backpack,"----",messageChannel)
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_2hz)
	time.Sleep(2*time.Second)
	backpack.BlinkRate(ledBackpack7Segment.Ht16k33_blink_off)

}

//
// Show the Current Minutes and Seconds from the system clock on the 7 Segment Display
// Send a Server Side Event with the same information
//
func showCurrentMinuteAndSeconds (backpack *ledBackpack7Segment.LedBackpack7Segment, messageChannel chan <- string, doneSignal chan bool) () {
	for {
		timeNow := time.Now()
		minutes := timeNow.Minute()
		seconds := timeNow.Second()
		minSecString := fmt.Sprintf("%2d:%02d",minutes,seconds)
		writeString(backpack,minSecString,messageChannel)
		time.Sleep(time.Second)
		// Check to see if we have had an OS interruption
		select {
			case <- doneSignal:
				return
			default:
		}
	}

}


func main() {
	fmt.Println("Entered Main") 
	// Channel for incoming signals (SIGTERM, SIGINT)
	osSignal := make(chan os.Signal, 1)
	// doneSignal indicates that an incoming Signal has been received
	doneSignal := make (chan bool, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	// Start monitoring for program interruption from OS
	// This isn't checked until the demo gets to displaying the current minute and seconds
	go func() {
		osSig := <- osSignal
		fmt.Println()
		fmt.Println(osSig)
		doneSignal <- true
	}()
	

	// ================== Setup Server Side Events ==================

	// Make a new Broker instance
	// This is for sending Server Side Events (SSE) to web pages
	broker := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	// Start processing events
	broker.Start()

	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	http.Handle("/events/", broker)
	// Files get served out of the ./templates directory
	http.Handle("/", http.FileServer(http.Dir("./templates")))
        // When we get a request at "/test", call `MainPageHandler`
	// in a new goroutine.	
	// The URL path should match the one in MainPageHandler
	http.HandleFunc("/test", http.HandlerFunc(MainPageHandler))

	fmt.Println("Starting Server") ;
	// Start the server and listen forever on port 8000.
	go func() { 
		http.ListenAndServe(":8000", nil)
	}()

	// =================== Seven Segment Display Demo =================

	// 1, 0x70 on Jetson TK1 I2C Pins 18&20
	backpack,err := ledBackpack7Segment.NewLedBackpack7Segment(1, 0x70) 
        fmt.Println("backpack: ",backpack)
        if err != nil {
		panic(err)
	}
	defer backpack.Close() 
	// Turn on the display
	backpack.Begin()
	// When we exit the program, turn the 7 Segment Display off
	defer backpack.End() 
	writeString(backpack,"    ",broker.messages)
	// Give clients a chance to connect
	time.Sleep(4*time.Second) 
	writeString(backpack,"1.234",broker.messages)
	time.Sleep(3*time.Second) 
	demoSeparator(backpack,broker.messages)
	// Should ignore last invalid character
	writeString(backpack,"9.8.:7.6.",broker.messages)
	time.Sleep(3*time.Second) 
	demoSeparator(backpack,broker.messages) ;


	writeString(backpack,"12:15",broker.messages) ;
	time.Sleep(3*time.Second)
	demoSeparator(backpack,broker.messages) ;
	// Countdown from 10
	showCountdown(backpack,10,broker.messages)
	demoSeparator(backpack,broker.messages) ;
	showClock(backpack,broker.messages)
	demoSeparator(backpack,broker.messages) ;
	showCurrentMinuteAndSeconds(backpack,broker.messages,doneSignal)
	fmt.Println("Cleaning Up, Done!")
}

