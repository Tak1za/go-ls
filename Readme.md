## go-ls

This is a simple implementation of the ever popular ***ls*** command from Unix systems. <br />
***ls*** command lists down the contents of the current working directory by default.

## Need ?

While working on Unix, this might probably be the most used command and if you switch to Windows every now and then like me there would have been countless number of times when instead of writing ***dir*** in the prompt (I don't like Powershell either), you end up writing ***ls*** (like me!). <br />
**Why not remove the hassle?!**

## Installation

1. Install Go
2. Clone the repository and enter the folder: <br />
   ```git clone https://github.com/Tak1za/go-ls.git && cd go-ls```
3. Make sure GOPATH is set. Check by running: ```go env```.
4. Add GOPATH to Windows Environment Variables.
5. Install the binary: <br />
   ```go build -o %GOPATH%/bin/ls```
6. If all goes well, run ```ls``` in any directory.
