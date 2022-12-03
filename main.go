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

func getLatestBlock(rpc string) int {
	var r Response
	queryURL := rpc + "/blocks/latest"
	client := &http.Client{}

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
	json.Unmarshal(body, &r)

	latestHeight, _ := strconv.Atoi(r.Block.LastCommit.Height)

	return latestHeight
}

func checkMissed(rpc string, startBlock int) error {
	var returnData Response

	isMissed := true
	queryURL := rpc + "/blocks/" + strconv.Itoa(startBlock)

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
