package modules

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/vote_bot/log"
	"github.com/vote_bot/models"
	"github.com/vote_bot/util"
)

func Proposals() {
	start := time.Now() // Start time

	// Import proposal list from DB
	models.GetDBProposals()

	// Find and change the DEPOSIT proposal that has been destroyed over time
	cancelProposals()

	// Import proposal list from API
	proposalAPI, err := getProposalsFromAPI()
	if err != nil {
		if err != nil {
			SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("GetProposalsFromAPI Failed : %v", err))
			log.Logger.Error.Fatal("GetProposalsFromAPI Failed : ", err)
		}
	}

	// Compare DB to API proposals
	compareChainProposals(proposalAPI)

	// Forward to Proposition Telegram in VOTING state
	votingAlert()

	elapsed := time.Since(start)
	log.Logger.Trace.Println("Elapsed Time (getNodeInfo) : ", elapsed)

}

// Take the proposition from LCD and put it in proposal API
func getProposalsFromAPI() ([]models.Proposals, error) {
	var proposalAPI []models.Proposals
	var wg1 sync.WaitGroup
	var mu sync.Mutex // Adding a Mutex

	for _, mainnetLCD := range models.Config.MainnetLCD {
		wg1.Add(1)
		go getProposals(&wg1, &mu, &proposalAPI, mainnetLCD.URL, mainnetLCD.ChainName)
	}
	wg1.Wait()

	return proposalAPI, nil
}

func getProposals(wg *sync.WaitGroup, mu *sync.Mutex, proposalAPI *[]models.Proposals, URL, chainName string) {
	defer wg.Done()
	var proposalsJSON models.ProposalsJson
	var v1proposalsJSON models.V1proposalsJSON
	proposalURL := URL + "/cosmos/gov/v1beta1/proposals?pagination.limit=20&pagination.reverse=true"

	res, err := util.CallUrl(proposalURL, 10) // timeout 10 sec
	if err != nil {
		log.Logger.Error.Println("GetProposalJson error:", err)
		if strings.Contains(err.Error(), "Client.Timeout") {
			res2, err := util.CallUrl(proposalURL, 10)
			res = res2
			if err != nil {
				SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("GetProposalJson error : %v", err))
				log.Logger.Error.Println("GetProposalJson error:", err)
				return
			}
		}
	}

	json.Unmarshal([]byte(res), &proposalsJSON)
	if len(proposalsJSON.Proposals) == 0 {
		proposalURL := URL + "/cosmos/gov/v1/proposals?pagination.limit=20&pagination.reverse=true"

		res, err := util.CallUrl(proposalURL, 10) // timeout 10 sec
		if err != nil {
			log.Logger.Error.Println("GetProposalJson error:", err)
			if strings.Contains(err.Error(), "Client.Timeout") {
				res2, err := util.CallUrl(proposalURL, 10)
				res = res2
				if err != nil {
					SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("GetV1ProposalJson error : %v", err))
					log.Logger.Error.Println("GetV1ProposalJson error:", err)
					return
				}
			}
		}
		json.Unmarshal([]byte(res), &v1proposalsJSON)
		log.Logger.Trace.Println(chainName+" v1", len(v1proposalsJSON.Proposals), proposalURL)
	} else {
		log.Logger.Trace.Println(chainName, len(proposalsJSON.Proposals), proposalURL)
	}
	if len(v1proposalsJSON.Proposals) == 0 && len(proposalsJSON.Proposals) == 0 {
		SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("Pasing error : %v LCD에 문제가 있습니다. 확인 해주세요.", chainName))
	}

	// Use Mutex to avoid modifying slices at the same time
	mu.Lock()
	defer mu.Unlock()

	// Each time, a new proposal is added to the proposition API
	for _, proposal := range proposalsJSON.Proposals {
		var votingStartTime, votingEndTime time.Time
		if proposal.VotingStartTime.Format(time.RFC3339) == "0001-01-01T00:00:00Z" {
			votingStartTime = time.Now().UTC()
			votingEndTime = time.Now().UTC()
		} else {
			votingStartTime = proposal.VotingStartTime
			votingEndTime = proposal.VotingEndTime
		}

		*proposalAPI = append(*proposalAPI,
			models.Proposals{
				ChainName:       chainName,
				ProposalID:      proposal.ProposalID,
				Type:            proposal.Content.Type,
				Title:           proposal.Content.Title,
				Status:          proposal.Status,
				SubmitTime:      proposal.SubmitTime,
				DepositEndTime:  proposal.DepositEndTime,
				VotingStartTime: votingStartTime,
				VotingEndTime:   votingEndTime,
				UpdateDate:      time.Now().UTC(),
				Description:     convertEnterToEscape(proposal.Content.Description),
			})
	}
	for _, proposal := range v1proposalsJSON.Proposals {
		var votingStartTime, votingEndTime time.Time
		if proposal.VotingStartTime.Format(time.RFC3339) == "0001-01-01T00:00:00Z" {
			votingStartTime = time.Now().UTC()
			votingEndTime = time.Now().UTC()
		} else {
			votingStartTime = proposal.VotingStartTime
			votingEndTime = proposal.VotingEndTime
		}

		var types string
		if len(proposal.Messages) == 0 {
			types = ""
		} else {
			types = proposal.Messages[0].Type
		}
		*proposalAPI = append(*proposalAPI,
			models.Proposals{
				ChainName:       chainName,
				ProposalID:      proposal.ID,
				Type:            types,
				Title:           proposal.Title,
				Status:          proposal.Status,
				SubmitTime:      proposal.SubmitTime,
				DepositEndTime:  proposal.DepositEndTime,
				VotingStartTime: votingStartTime,
				VotingEndTime:   votingEndTime,
				UpdateDate:      time.Now().UTC(),
				Description:     convertEnterToEscape(proposal.Summary),
			})
	}
}

// A function that compares the list of proposals for two Structures
func compareChainProposals(newProposals []models.Proposals) {
	oldProposals := models.ProposalsDB
	// Convert previous proposals to maps for quick search
	oldProposalsMap := make(map[string]map[int]models.Proposals)
	for _, oldProposal := range oldProposals {
		// Create a map for each chain independently
		chainMap, exists := oldProposalsMap[oldProposal.ChainName]
		if !exists {
			chainMap = make(map[int]models.Proposals)
			oldProposalsMap[oldProposal.ChainName] = chainMap
		}

		oldID, _ := strconv.Atoi(oldProposal.ProposalID)

		// When adding a value to a map in each chain, create a new value and add it to the map
		if _, exists := chainMap[oldID]; !exists {
			chainMap[oldID] = models.Proposals{
				ChainName:       oldProposal.ChainName,
				ProposalID:      oldProposal.ProposalID,
				Type:            oldProposal.Type,
				Title:           oldProposal.Title,
				Status:          oldProposal.Status,
				SubmitTime:      oldProposal.SubmitTime,
				DepositEndTime:  oldProposal.DepositEndTime,
				VotingStartTime: oldProposal.VotingStartTime,
				VotingEndTime:   oldProposal.VotingEndTime,
				// VoteOption:      oldProposal.VoteOption,
				// VoteTx:          oldProposal.VoteTx,
				UpdateDate:  oldProposal.UpdateDate,
				Description: oldProposal.Description,
			}
		}
	}
	// Current time output
	log.Logger.Trace.Println(time.Now().UTC())

	// Check the new proposal list to see what changes have been made
	for _, newProposal := range newProposals {
		// Add the current chain name to the previous proposal
		oldChainMap, chainNameExists := oldProposalsMap[newProposal.ChainName]
		if !chainNameExists {
			oldChainMap = make(map[int]models.Proposals)
			// If there is no chain name, add the chain name to the map
			oldProposalsMap[newProposal.ChainName] = oldChainMap
		}
		newID, _ := strconv.Atoi(newProposal.ProposalID)
		oldProposal, exists := oldChainMap[newID]
		oldChainMap[newID] = newProposal

		if !exists {
			// If a new proposal is added
			log.Logger.Trace.Printf("[%s] %s 새로운 프로포절이 추가되었습니다.\n", newProposal.ChainName, newProposal.ProposalID)
			// Update to DB
			err := insertProposal(newProposal)
			if err != nil {
				SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("Insert Proposal failed: %v\n", err))
				log.Logger.Error.Fatal("Insert Proposal failed: ", err)
			}
			// Add added proposals to oldProducts
			oldProposals = append(oldProposals, newProposal)
			models.ProposalsDB = append(models.ProposalsDB, newProposal)

		} else if oldProposal.ProposalID == newProposal.ProposalID && oldProposal.ChainName == newProposal.ChainName {

			if diff := compareProposals(oldProposal, newProposal); diff != "" {
				// If the proposal has been modified
				log.Logger.Trace.Printf("[%s] %s 프로포절이 \n%s 상태에서 \n%s 상태로 변경되었습니다. \n", newProposal.ChainName, newProposal.ProposalID, oldProposal.Status, newProposal.Status)
				// Update to DB
				err := updateProposal(newProposal)
				if err != nil {
					SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("Insert Proposal failed: %v\n", err))
					log.Logger.Error.Fatal("Update Proposal failed: ", err)
				}
				// Update modified product information to oldProducts
				models.ProposalsDB[models.DBFind(newProposal.ChainName, newProposal.ProposalID)] = newProposal
			}
		}
	}
}

// Voting alarm module. Alarm when voting starts. Check by giving another alarm when less than 24 hours
func votingAlert() {
	oldProposals := models.ProposalsDB
	for _, proposal := range oldProposals {
		if proposal.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
			if proposal.MsgID == 0 {
				var explore string
				for _, LCD := range models.Config.MainnetLCD {
					if LCD.ChainName == proposal.ChainName {
						lowerChainName := strings.ToLower(proposal.ChainName)
						if LCD.EX == "https://www.mintscan.io/" || LCD.EX == "https://bigdipper.live/" {
							explore = LCD.EX + lowerChainName + "/proposals/" + proposal.ProposalID
						} else {
							explore = LCD.EX + "/proposal/" + proposal.ProposalID
						}
					}
				}
				log.Logger.Trace.Printf("[%s] %s The proposal voting has started. Current Time %v, End Time %v \n%s", proposal.ChainName, proposal.ProposalID, time.Now().UTC().Format(time.RFC3339), proposal.VotingEndTime.Format(time.RFC3339), explore)
				SendTelegramProposalMessage(proposal, explore)

			}
		}
	}

	// Categorized to output a statement that says 24 hours are left later
	for _, proposal := range oldProposals {
		if proposal.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
			const remainTime = 24
			if proposal.VotingEndTime.Sub(time.Now().UTC()) < remainTime*time.Hour && proposal.IsAlert == false {
				log.Logger.Trace.Printf("[%s] %s Proposal voting period is less than 24 hours. Current Time %v, End Time %v \n", proposal.ChainName, proposal.ProposalID, time.Now().UTC().Format(time.RFC3339), proposal.VotingEndTime.Format(time.RFC3339))
				ReplyTelegramButtonMessage(proposal.MsgID, fmt.Sprintf("[%s] %s Proposal voting period is less than 24 hours. Current Time : %v\nEnd Time : %v \n", proposal.ChainName, proposal.ProposalID, time.Now().UTC().Format(time.RFC3339), proposal.VotingEndTime.Format(time.RFC3339)))
				updateAlertProposal(proposal.ChainName, proposal.ProposalID)
			}
		}
	}
}

func compareProposals(oldProposal, newProposal models.Proposals) string {
	opt := cmp.FilterPath(func(p cmp.Path) bool {
		// UpdateDate 필드를 무시
		return p.Last().String() == "UpdateDate"
	}, cmp.Ignore())

	if oldProposal.Status != newProposal.Status {
		if diff := cmp.Diff(oldProposal, newProposal, opt); diff != "" {
			return diff
		}
	}
	return ""
}

/*
Check the DEPOSIT_END_TIME time in DEPOSIT state, and check PROPOSAL_STATUS_CANCELED if it exceeds the current time

The process of changing the canceled PROPOSALs to cancel because the deposit was not made in the DEPOSIT state
*/
func cancelProposals() {
	oldProposals := models.ProposalsDB
	for i, proposal := range oldProposals {
		if proposal.Status == "PROPOSAL_STATUS_DEPOSIT_PERIOD" && proposal.DepositEndTime.Sub(time.Now()) < time.Hour {
			// If the proposal has been modified
			log.Logger.Trace.Printf("[%s] %s 프로포절이 DEPOSIT 에서 캔슬되었습니다. \n", proposal.ChainName, proposal.ProposalID)
			// Update to DB
			err := updateCanceledProposal(proposal.ChainName, proposal.ProposalID)
			if err != nil {
				SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("update CanceledProposal failed: %v\n", err))
				log.Logger.Error.Fatal("update CanceledProposal failed: ", err)
			}
			// Update modified proposal information to ProposalsDB
			models.ProposalsDB[i] = proposal
		}
	}
}

// Functions which is insert to DB
func insertProposal(proposal models.Proposals) error {
	err := models.InsertProposal(proposal)
	if err != nil {
		log.Logger.Error.Println("Insert Proposal failed: ", err)
		return err
	}
	return nil
}

// Functions which is updateProposal to DB
func updateProposal(proposal models.Proposals) error {
	err := models.UpdateDBProposals(proposal)
	if err != nil {
		log.Logger.Error.Println("update Proposal failed: ", err)
		return err
	}
	return nil
}

// Functions which is update calceled Proposal to DB
func updateCanceledProposal(chainName, proposalID string) error {
	proposal := models.ProposalsDB[models.DBFind(chainName, proposalID)]
	proposal.Status = "PROPOSAL_STATUS_CANCELED"
	err := models.UpdateDBProposals(proposal)
	if err != nil {
		log.Logger.Error.Println("update CanceledProposal failed: ", err)
		return err
	}
	return nil
}

// Functions which is update alerted Proposal to DB
func updateAlertProposal(chainName, proposalID string) error {
	proposal := models.ProposalsDB[models.DBFind(chainName, proposalID)]
	proposal.IsAlert = true
	err := models.UpdateDBProposals(proposal)
	if err != nil {
		log.Logger.Error.Println("update AlertProposal failed: ", err)
		return err
	}
	return nil
}

func convertEnterToEscape(input string) string {
	// Change Enter to <br>
	result := strings.ReplaceAll(input, "\n", "<br />")
	return result
}
