package config

import "os"

var Url = os.Getenv("URL")
var RiskID = os.Getenv("RISK_ID")
var CompanyID = os.Getenv("COMPANY_ID")
var SecuriwiserApi = os.Getenv("SECURIWISER_API")
