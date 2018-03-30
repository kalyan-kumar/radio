# YouTube Radio

An application to share music. All the users
connected will be in sync. Anyone can request
a song by searching for it's name on the
webpage, which will be added to the queue.

Backend is implemented in `golang`. The server
uses YouTube Data API to query for the search
string to fetch the top match.

Frontend is simple HTML and JavaScript, with
the player rendered by through IFrame API.

### Usage
Setup [Go](https://golang.org/) environment.
To run the server Google API key is required.
Generate one at [this](https://console.developers.google.com/apis/dashboard)
console and add it in `src/Jockey.go` file in
this package.
```
go get github.com/kalyan-kumar/radio
cd $GOPATH/src/github.com/kalyan-kumar/radio
go build
./radio
```
To connect to the server, hit the url:

`<servers-ip>:8042/radio`

### Development
The code is in alpha quality, developed
occasionally. Areas of improvement include
the displayed webpage and some convenience
features at the server's disposal.

A multi-station version could be implemented
so that different genres of songs could be
played on different stations.