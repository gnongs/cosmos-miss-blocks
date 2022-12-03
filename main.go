package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Block struct {
		LastCommit struct {
			Height     string `json: "height"`
			Signatures []struct {
				ValidatorAddress string `json:"validator_address"`
			} `json:"signatures"`
		} `json:"last_commit"`
	} `json:"block"`
}

const (
	rpcAddr    = "https://agoric-mainnet-rpc.allthatnode.com:1317"
	targetAddr = "9A4CCB10366FDBAE60E5CE58495684FBEF959FDE"
)

var (
	missBlocks []string
)

func main() {
	for {
		err := getMissBlock(rpcAddr)
		if err != nil {
			fmt.Sprintln(err)
		}
		log.Info().Str("module", "miss").Msgf("Current miss blocks: %v", missBlocks)
		time.Sleep(7 * time.Second)
	}
}

func getMissBlock(rpc string) error {
	var returnData Response

	isMissed := true
	queryURL := rpc + "/blocks/latest"

	client := &http.Client{}

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &returnData)

	for _, s := range returnData.Block.LastCommit.Signatures {
		if targetAddr == s.ValidatorAddress {
			isMissed = false
			break
		}
	}

	if isMissed {
		missBlocks = append(missBlocks, returnData.Block.LastCommit.Height)
	}

	return nil
}
