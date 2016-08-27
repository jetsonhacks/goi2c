/*
 	Package ledBackpack7Segment implements a simple library for the
	Adafruit 0.56" 4-Digit 7-Segment Display w/I2C Backpack
*/
package ledBackpack7Segment

import (
	_ "fmt"
	"github.com/jetsonhacks/goi2c/i2c"
)

// Ht16k33 is the name of the I2C display driver chip 
const (
	Ht16k33_blink_cmd = 0x80
	Ht16k33_blink_displayon = 0x01
	Ht16k33_blink_off = 0
	Ht16k33_blink_2hz = 1
	Ht16k33_blink_1hz = 2
	Ht16k33_blink_halfhz  = 3
	Ht16k33_cmd_brightness = 0xE0
	seven_digits = 5
)

// Translation table for characters
var numbertable = [...] uint8 {
    0x3F, /* 0 */
    0x06, /* 1 */
    0x5B, /* 2 */
    0x4F, /* 3 */
    0x66, /* 4 */
    0x6D, /* 5 */
    0x7D, /* 6 */
    0x07, /* 7 */
    0x7F, /* 8 */
    0x6F, /* 9 */
    0x77, /* a */
    0x7C, /* b */
    0x39, /* C */
    0x5E, /* d */
    0x79, /* E */
    0x71, /* F */
}

// Translate from a character to the display glyph
var translateTable = map[rune]uint8 {
    '0' : 0x3F, 
    '1' : 0x06, 
    '2' : 0x5B, 
    '3' : 0x4F, 
    '4' : 0x66, 
    '5' : 0x6D, 
    '6' : 0x7D, 
    '7' : 0x07, 
    '8' : 0x7F, 
    '9' : 0x6F,
    'a' : 0x77, 
    'b' : 0x7C, 
    'c' : 0x39, 
    'd' : 0x5E, 
    'e' : 0x79, 
    'f' : 0x71, 
    'A' : 0x77, 
    'B' : 0x7C, 
    'C' : 0x39, 
    'D' : 0x5E, 
    'E' : 0x79, 
    'F' : 0x71, 
    ' ' : 0x00,
    '-' : 0x40,  
}


type LedBackpack7Segment struct {
        I2CDevice *i2c.I2C
        Bus int
        FileDescriptor int 
        I2CAddress int 
        Error int
	DisplayBuffer [16] uint8 
        // Position int
}

// The default address is 0x70 
func NewLedBackpack7Segment ( bus, address int ) (*LedBackpack7Segment, error) {
	i2cdevice,err := i2c.NewI2C(bus, address) 
	if err != nil {
		return nil, err
	}
	// Setup the structure to return
        backpack := &LedBackpack7Segment{Bus: bus, I2CAddress: address, I2CDevice: i2cdevice}
	return backpack, nil 
}

func (backpack *LedBackpack7Segment) Close() error {
	return backpack.I2CDevice.Close() 
}

func (backpack *LedBackpack7Segment) i2cwrite ( writeValue uint8 ) (error) {
	err := backpack.I2CDevice.WriteUint8(writeValue)
	if err != nil {
		return err ;
	}
	return nil ;
}

func (backpack *LedBackpack7Segment) Begin () {
	backpack.i2cwrite(0x21)  // turn on oscillator
	backpack.BlinkRate(Ht16k33_blink_off) 
	backpack.SetBrightness(15)
}

func (ht16K33 *LedBackpack7Segment) End () {
	ht16K33.i2cwrite(0x20)  // turn off oscillator
}

func (backpack *LedBackpack7Segment) BlinkRate ( blinkRate uint8 ) {
	if blinkRate > 3 {
		blinkRate = 0	// turn off if not sure
	}
	backpack.i2cwrite(Ht16k33_blink_cmd | Ht16k33_blink_displayon | (blinkRate << 1))
}

func (backpack *LedBackpack7Segment) SetBrightness ( brightness uint8 ) {
	if brightness > 15 {
		brightness = 15
	}
	backpack.i2cwrite(brightness | Ht16k33_cmd_brightness)
}


func (backpack *LedBackpack7Segment) Clear () {
	for i := 0; i < len(backpack.DisplayBuffer); i++ {
		backpack.DisplayBuffer[i] = 0 
	}
}



func (backpack *LedBackpack7Segment) WriteDigitNum (place uint8, number uint8, dot bool) {

	if place > 4 || place < 0 {
		return 
	}
	var mask uint8 = 0 
	if dot == true {
		mask = 1 << 7
	}

	backpack.WriteDigitRaw(place,numbertable[number] | mask) 
}


func (backpack *LedBackpack7Segment) WriteDigitRaw (place uint8, bitmask uint8) {
	if place > 4 || place < 0 {
		return 
	}
	backpack.DisplayBuffer[place*2] = bitmask;
}

func (backpack *LedBackpack7Segment) WriteDisplay () (error) {
	// fmt.Println(backpack.DisplayBuffer[:]) 
	err := backpack.I2CDevice.WriteI2CBlock(0x00, backpack.DisplayBuffer[:])
	return err 
}

func (backpack *LedBackpack7Segment) WriteError () {
	backpack.WriteString("----") 
}

func (backpack *LedBackpack7Segment) WriteString (writeString string) {
	backpack.Clear() 
	var position uint8
	position = 0 
	var mask uint8 
	for _, char := range writeString {
		// On the 7 segment display, position 2 is the ':'
		// Only the ':' can be displayed at that position
		if position == 2 {
			if char != '.' {
				position ++ 
				if char == ':' {
					// Write the colon glyph
					backpack.WriteDigitRaw(2,2)
					continue
				}
			}
		} 

		if char == '.' {
			// fmt.Println("Char .")
			// fmt.Println("DisplayBuffer: ",backpack.DisplayBuffer)
			place := position 
			place = place - 1 
			// fmt.Println("Place: ",place)
			if place == 255 {	
				// string started with '.'
				place = 0
				position++ 
			}
			backpack.DisplayBuffer[place*2] = backpack.DisplayBuffer[place*2] | (1 << 7)
			// fmt.Println("DisplayBuffer: ",backpack.DisplayBuffer)
			continue
		}

		mask = 0
		number := translateTable[char] ;
		if position > 4 {
			// string to write is over maximum size
			break 
		}
		backpack.WriteDigitRaw(position,number | mask) 
		position++
	}
	backpack.WriteDisplay() 
}
