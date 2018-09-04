package eos

import (
	"io/ioutil"
	"log"
)

type logger struct {
	Decoder    *log.Logger
	Encoder    *log.Logger
	ABIDecoder *log.Logger
	ABIEncoder *log.Logger
}

var Logger = &logger{
	Decoder: log.New(ioutil.Discard, "[Decoder]		", 0),
	Encoder: log.New(ioutil.Discard, "[Encoder]		", 0),
	ABIDecoder: log.New(ioutil.Discard, "[ABIDecoder]	", 0),
	ABIEncoder: log.New(ioutil.Discard, "[ABIEncoder]	", 0),
}
