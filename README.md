# goi2c
Jetson I2C interface in Go (Golang)
This repository contains some sketches which interface the Jetson Development Kits
with an I2C interface and web servers.

These examples interface with an Adafruit 0.56" 4-Digit 7-Segment Display w/I2C Backpack.
The configuration of the interface will depend on how the I2C Backpack is wired to the the Jetson.
This can be changed in the 'main' function of each example.

The folders:
devices          - holds code for interfacing with the I2C backback 7 segment display
i2c              - I2C interface library, derived from https://github.com/SpaceLeap/go-embedded/blob/master/i2c/i2c.go
i2cExample       - Interface with the Jetson to the 7 Segment Display
i2cExampleServer - A more involved example. Shows digits on the Segment Display, and broadcasts over the Web using Server Side Events (SSE)
                   HTML 5 Browser required.

Additional Licenses:

-----------------------------------------------------

Go I2C library
Originally from:
https://github.com/SpaceLeap/go-embedded
Specifically:
https://github.com/SpaceLeap/go-embedded/blob/master/i2c/i2c.go

License Terms:
The MIT License (MIT)

Copyright (c) 2013 Erik Unger

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
Contact GitHub API Training Shop Blog About

-----------------------------------------------------------------------------
Server Code mostly derived from:

https://github.com/kljensen/golang-html5-sse-example

License Terms:

License (the Unlicense)

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or distribute this software, either in source code form or as a compiled binary, for any purpose, commercial or non-commercial, and by any means.

In jurisdictions that recognize copyright laws, the author or authors of this software dedicate any and all copyright interest in the software to the public domain. We make this dedication for the benefit of the public at large and to the detriment of our heirs and successors. We intend this dedication to be an overt act of relinquishment in perpetuity of all present and future rights to this software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to http://unlicense.org/

------------------------------------------

 * segment-display.js
 *
 * Copyright 2012, Rüdiger Appel
 * http://www.3quarks.com
 * Published under Creative Commons 3.0 License.
 *
 * Date: 2012-02-14
 * Version: 1.0.0

