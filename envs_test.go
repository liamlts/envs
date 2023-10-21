package envs

import (
	"log"
	"testing"
)

func TestEnvs(t *testing.T) {
	hmap, err := parseEnvSFile("./envs/.envs")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hmap["S3_BUCKET"], hmap["MALFORM"])

}
