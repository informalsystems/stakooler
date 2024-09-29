package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type BlockResponse struct {
	BlockID struct {
		Hash          string `json:"hash"`
		PartSetHeader struct {
			Total int    `json:"total"`
			Hash  string `json:"hash"`
		} `json:"part_set_header"`
	} `json:"block_id"`
	Block struct {
		Header struct {
			Version struct {
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"version"`
			ChainID     string    `json:"chain_id"`
			Height      string    `json:"height"`
			Time        time.Time `json:"time"`
			LastBlockID struct {
				Hash          string `json:"hash"`
				PartSetHeader struct {
					Total int    `json:"total"`
					Hash  string `json:"hash"`
				} `json:"part_set_header"`
			} `json:"last_block_id"`
			LastCommitHash     string `json:"last_commit_hash"`
			DataHash           string `json:"data_hash"`
			ValidatorsHash     string `json:"validators_hash"`
			NextValidatorsHash string `json:"next_validators_hash"`
			ConsensusHash      string `json:"consensus_hash"`
			AppHash            string `json:"app_hash"`
			LastResultsHash    string `json:"last_results_hash"`
			EvidenceHash       string `json:"evidence_hash"`
			ProposerAddress    string `json:"proposer_address"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
		Evidence struct {
			Evidence []struct {
				DuplicateVoteEvidence struct {
					VoteA struct {
						Type    string `json:"type"`
						Height  string `json:"height"`
						Round   int    `json:"round"`
						BlockID struct {
							Hash          string `json:"hash"`
							PartSetHeader struct {
								Total int    `json:"total"`
								Hash  string `json:"hash"`
							} `json:"part_set_header"`
						} `json:"block_id"`
						Timestamp        time.Time `json:"timestamp"`
						ValidatorAddress string    `json:"validator_address"`
						ValidatorIndex   int       `json:"validator_index"`
						Signature        string    `json:"signature"`
					} `json:"vote_a"`
					VoteB struct {
						Type    string `json:"type"`
						Height  string `json:"height"`
						Round   int    `json:"round"`
						BlockID struct {
							Hash          string `json:"hash"`
							PartSetHeader struct {
								Total int    `json:"total"`
								Hash  string `json:"hash"`
							} `json:"part_set_header"`
						} `json:"block_id"`
						Timestamp        time.Time `json:"timestamp"`
						ValidatorAddress string    `json:"validator_address"`
						ValidatorIndex   int       `json:"validator_index"`
						Signature        string    `json:"signature"`
					} `json:"vote_b"`
					TotalVotingPower string    `json:"total_voting_power"`
					ValidatorPower   string    `json:"validator_power"`
					Timestamp        time.Time `json:"timestamp"`
				} `json:"duplicate_vote_evidence"`
				LightClientAttackEvidence struct {
					ConflictingBlock struct {
						SignedHeader struct {
							Header struct {
								Version struct {
									Block string `json:"block"`
									App   string `json:"app"`
								} `json:"version"`
								ChainID     string    `json:"chain_id"`
								Height      string    `json:"height"`
								Time        time.Time `json:"time"`
								LastBlockID struct {
									Hash          string `json:"hash"`
									PartSetHeader struct {
										Total int    `json:"total"`
										Hash  string `json:"hash"`
									} `json:"part_set_header"`
								} `json:"last_block_id"`
								LastCommitHash     string `json:"last_commit_hash"`
								DataHash           string `json:"data_hash"`
								ValidatorsHash     string `json:"validators_hash"`
								NextValidatorsHash string `json:"next_validators_hash"`
								ConsensusHash      string `json:"consensus_hash"`
								AppHash            string `json:"app_hash"`
								LastResultsHash    string `json:"last_results_hash"`
								EvidenceHash       string `json:"evidence_hash"`
								ProposerAddress    string `json:"proposer_address"`
							} `json:"header"`
							Commit struct {
								Height  string `json:"height"`
								Round   int    `json:"round"`
								BlockID struct {
									Hash          string `json:"hash"`
									PartSetHeader struct {
										Total int    `json:"total"`
										Hash  string `json:"hash"`
									} `json:"part_set_header"`
								} `json:"block_id"`
								Signatures []struct {
									BlockIDFlag      string    `json:"block_id_flag"`
									ValidatorAddress string    `json:"validator_address"`
									Timestamp        time.Time `json:"timestamp"`
									Signature        string    `json:"signature"`
								} `json:"signatures"`
							} `json:"commit"`
						} `json:"signed_header"`
						ValidatorSet struct {
							Validators []struct {
								Address string `json:"address"`
								PubKey  struct {
									Ed25519   string `json:"ed25519"`
									Secp256K1 string `json:"secp256k1"`
								} `json:"pub_key"`
								VotingPower      string `json:"voting_power"`
								ProposerPriority string `json:"proposer_priority"`
							} `json:"validators"`
							Proposer struct {
								Address string `json:"address"`
								PubKey  struct {
									Ed25519   string `json:"ed25519"`
									Secp256K1 string `json:"secp256k1"`
								} `json:"pub_key"`
								VotingPower      string `json:"voting_power"`
								ProposerPriority string `json:"proposer_priority"`
							} `json:"proposer"`
							TotalVotingPower string `json:"total_voting_power"`
						} `json:"validator_set"`
					} `json:"conflicting_block"`
					CommonHeight        string `json:"common_height"`
					ByzantineValidators []struct {
						Address string `json:"address"`
						PubKey  struct {
							Ed25519   string `json:"ed25519"`
							Secp256K1 string `json:"secp256k1"`
						} `json:"pub_key"`
						VotingPower      string `json:"voting_power"`
						ProposerPriority string `json:"proposer_priority"`
					} `json:"byzantine_validators"`
					TotalVotingPower string    `json:"total_voting_power"`
					Timestamp        time.Time `json:"timestamp"`
				} `json:"light_client_attack_evidence"`
			} `json:"evidence"`
		} `json:"evidence"`
		LastCommit struct {
			Height  string `json:"height"`
			Round   int    `json:"round"`
			BlockID struct {
				Hash          string `json:"hash"`
				PartSetHeader struct {
					Total int    `json:"total"`
					Hash  string `json:"hash"`
				} `json:"part_set_header"`
			} `json:"block_id"`
			Signatures []struct {
				BlockIDFlag      string    `json:"block_id_flag"`
				ValidatorAddress string    `json:"validator_address"`
				Timestamp        time.Time `json:"timestamp"`
				Signature        string    `json:"signature"`
			} `json:"signatures"`
		} `json:"last_commit"`
	} `json:"block"`
}

func (b *BlockResponse) GetLatestBlock(endpoint string, client *http.Client) error {
	var body []byte

	url := endpoint + "/cosmos/base/tendermint/v1beta1/blocks/latest"
	body, err := HttpGet(url, client)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, b); err != nil {
		return err
	}
	return nil
}

func (b *BlockResponse) GetBlock(height string, endpoint string) error {
	url := endpoint + "/cosmos/base/tendermint/v1beta1/blocks/" + height
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, b)
	if err != nil {
		return err
	}
	return nil
}
