package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	height := getLatestBlock(rpcAddr)
	log.Info().Str("module", "miss").Msgf("Start height: %v", height)

	for {
		err := checkMissed(rpcAddr, height)
		if err != nil {
			fmt.Sprintln(err)
		}
		log.Info().Str("module", "miss").Msgf("This checked height: %v, Total miss blocks: %v", height, missBlocks)
		time.Sleep(7 * time.Second)

		height++
	}
}

func query(route string) Response {
	var data Response

	queryURL := rpcAddr + route
	client := http.Client{}

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		fmt.Sprintln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Sprintln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Sprintln(err)
	}

	json.Unmarshal(body, &data)

	return data
}

func getLatestBlock(rpc string) int {
	route := "/blocks/latest"
	respData := query(route)

	latestHeight, _ := strconv.Atoi(respData.Block.LastCommit.Height)

	return latestHeight
}

func checkMissed(rpc string, startBlock int) error {
	isMissed := true
	route := "/blocks/" + strconv.Itoa(startBlock)

	respData := query(route)

	for _, s := range respData.Block.LastCommit.Signatures {
		if targetAddr == s.ValidatorAddress {
			isMissed = false
			break
		}
	}

	if isMissed {
		missBlocks = append(missBlocks, respData.Block.LastCommit.Height)
	}

	return nil
}
