package main

import (
	"nsrlFerret/api"
	"nsrlFerret/db"
	"nsrlFerret/nsrl"
)

func main() {
	nsrl.GetNSrl()
	api.StartApi(db.ProcessNSRLtxt("rds_modernm/rds_modernm/NSRLFile.txt"))
}
