package consensus

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetValidConsensus() ([]byte, error) {

	consesusList := []string{
		Tor26Consensus, Moria1Consensus, DizumConsensus,
		GabelmooConsensus, DannenbergConsensus, MaatuskaConsensus,
		LongclawConsensus, BastetConsensus, FaravaharConsensus,
	}

	for _, v := range consesusList {

		req, err := http.NewRequest("GET", v, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Accept", "text/plain")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			continue
		}

		var reader io.Reader
		if resp.Header.Get("Content-Encoding") == "deflate" {
			read, err := zlib.NewReader(resp.Body)
			if err != nil {
				continue
			}
			reader = read
		} else {
			reader = resp.Body
		}

		conensusBlob, err := io.ReadAll(reader)
		if err != nil {
			continue
		}
		if !ValidConsensus(conensusBlob) {
			continue
		}

		ParseConsensus(conensusBlob)
		fmt.Println(v)
		return conensusBlob, nil

	}

	return nil, errors.New("failed to fetch all consensus")
}
