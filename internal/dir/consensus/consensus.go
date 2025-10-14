package consensus

import (
	"compress/zlib"
	"errors"
	"io"
	"net/http"
)

func GetValidConsensus() ([]byte, error) {

	consesusList := []string{
		Tor26Consensus, Moria1Consensus, DizumConsensus,
		GabelmooConsensus, DannenbergConsensus, MaatuskaConsensus,
		LongclawConsensus, BastetConsensus, FaravaharConsensus,
	}

	for _, v := range consesusList {

		resp, err := http.Get(v)
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
			defer read.Close()
			reader = read
		} else {
			reader = resp.Body
		}

		conensusBlob, err := io.ReadAll(reader)
		if err != nil {
			continue
		}
		if !validConsensus(conensusBlob) {
			continue
		}

		return conensusBlob, nil

	}

	return nil, errors.New("failed to fetch all consensus")
}
