package config

import "os"

var Url = os.Getenv("URL")
var RiskID = os.Getenv("RISK_ID")
