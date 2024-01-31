package models

import "time"

// API Proposal
type ProposalJson struct {
	Proposal struct {
		ProposalID string `json:"proposal_id"`
		Content    struct {
			Type        string `json:"@type"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Changes     []struct {
				Subspace string `json:"subspace"`
				Key      string `json:"key"`
				Value    string `json:"value"`
			} `json:"changes"`
		} `json:"content"`
		Status           string `json:"status"`
		FinalTallyResult struct {
			Yes        string `json:"yes"`
			Abstain    string `json:"abstain"`
			No         string `json:"no"`
			NoWithVeto string `json:"no_with_veto"`
		} `json:"final_tally_result"`
		SubmitTime     time.Time `json:"submit_time"`
		DepositEndTime time.Time `json:"deposit_end_time"`
		TotalDeposit   []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total_deposit"`
		VotingStartTime time.Time `json:"voting_start_time"`
		VotingEndTime   time.Time `json:"voting_end_time"`
	} `json:"proposal"`
}

type V1proposalsJSON struct {
	Proposals []struct {
		ID       string `json:"id"`
		Messages []struct {
			Type    string `json:"@type"`
			Content struct {
				Type               string `json:"@type"`
				Title              string `json:"title"`
				Description        string `json:"description"`
				SubjectClientID    string `json:"subject_client_id"`
				SubstituteClientID string `json:"substitute_client_id"`
				Changes            []struct {
					Subspace string `json:"subspace"`
					Key      string `json:"key"`
					Value    string `json:"value"`
				} `json:"changes"`
			} `json:"content"`
			Authority string `json:"authority"`
			Plan      struct {
				Name                string      `json:"name"`
				Time                time.Time   `json:"time"`
				Height              string      `json:"height"`
				Info                string      `json:"info"`
				UpgradedClientState interface{} `json:"upgraded_client_state"`
			} `json:"plan"`
			Params struct {
				MintDenom            string `json:"mint_denom"`
				InflationRateChange  string `json:"inflation_rate_change"`
				InflationMax         string `json:"inflation_max"`
				InflationMin         string `json:"inflation_min"`
				GoalBonded           string `json:"goal_bonded"`
				BlocksPerYear        string `json:"blocks_per_year"`
				MaxBundleSize        int    `json:"max_bundle_size"`
				EscrowAccountAddress string `json:"escrow_account_address"`
				ReserveFee           struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"reserve_fee"`
				MinBidIncrement struct {
					Denom  string `json:"denom"`
					Amount string `json:"amount"`
				} `json:"min_bid_increment"`
				FrontRunningProtection bool   `json:"front_running_protection"`
				ProposerFee            string `json:"proposer_fee"`
			} `json:"params"`
		} `json:"messages"`
		Status           string `json:"status"`
		FinalTallyResult struct {
			YesCount        string `json:"yes_count"`
			AbstainCount    string `json:"abstain_count"`
			NoCount         string `json:"no_count"`
			NoWithVetoCount string `json:"no_with_veto_count"`
		} `json:"final_tally_result"`
		SubmitTime     time.Time `json:"submit_time"`
		DepositEndTime time.Time `json:"deposit_end_time"`
		TotalDeposit   []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total_deposit"`
		VotingStartTime time.Time `json:"voting_start_time"`
		VotingEndTime   time.Time `json:"voting_end_time"`
		Metadata        string    `json:"metadata"`
		Title           string    `json:"title"`
		Summary         string    `json:"summary"`
		Proposer        string    `json:"proposer"`
	} `json:"proposals"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}

// API Proposals
type ProposalsJson struct {
	Proposals []struct {
		ProposalID string `json:"proposal_id"`
		Content    struct {
			Type        string `json:"@type"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Recipient   string `json:"recipient"`
			Amount      []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
			Changes []struct {
				Subspace string `json:"subspace"`
				Key      string `json:"key"`
				Value    string `json:"value"`
			} `json:"changes"`
			SubjectClientID    string `json:"subject_client_id"`
			SubstituteClientID string `json:"substitute_client_id"`
			Plan               struct {
				Name                string      `json:"name"`
				Time                time.Time   `json:"time"`
				Height              string      `json:"height"`
				Info                string      `json:"info"`
				UpgradedClientState interface{} `json:"upgraded_client_state"`
			} `json:"plan"`
		} `json:"content"`
		Status           string `json:"status"`
		FinalTallyResult struct {
			Yes        string `json:"yes"`
			Abstain    string `json:"abstain"`
			No         string `json:"no"`
			NoWithVeto string `json:"no_with_veto"`
		} `json:"final_tally_result"`
		SubmitTime     time.Time `json:"submit_time"`
		DepositEndTime time.Time `json:"deposit_end_time"`
		TotalDeposit   []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"total_deposit"`
		VotingStartTime time.Time `json:"voting_start_time"`
		VotingEndTime   time.Time `json:"voting_end_time"`
	} `json:"proposals"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}
