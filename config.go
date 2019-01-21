package main

import "time"

const countyWorkers     = 3
const detailWorkers     = 20
const httpRetryCount    = 7
const httpRetryWait     = time.Millisecond * 750
const httpNextDelay     = time.Millisecond * 250
const httpClientTimeout = time.Second * 60
const baseUrl           = "http://prog.rarom.ro/servicenou/"

var counties = []string{"AB", "AG", "AR", "BB", "BC", "BH", "BN", "BT", "BV", "BR", "BZ", "CS",
	"CL", "CJ", "CT", "CV", "DB", "DJ", "GL", "GR", "GJ", "HR", "HD", "IL", "IS", "IF", "MM",
	"MH", "MS", "NT", "OT", "PH", "SM", "SJ", "SB", "SV", "TR", "TM", "TL", "VS", "VL", "VN"}