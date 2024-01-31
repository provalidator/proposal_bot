package models

import (
	"github.com/vote_bot/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	DbStr string

	ProposalsDB []Proposals
)

func ConnectDatabase() error {
	var err error

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       Config.Database.MysqlUserID + ":" + Config.Database.MysqlUserPW + "@tcp(" + Config.Database.MysqlServerURL + ":" + Config.Database.MysqlServerPort + ")/" + Config.Database.MysqlSelectDBName + "?charset=utf8mb4&parseTime=true",
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		log.Logger.Error.Fatal("Connect Database failed: ", err)
		return err
	}
	return nil
}

// Import all records
func GetDBProposals() {
	DB.Find(&ProposalsDB)
}

// Save if you have a new proposal
func InsertProposal(proposal Proposals) error {

	// Start Transaction
	tx := DB.Begin()

	if tx.Error != nil {
		log.Logger.Error.Println("Insert Proposal Info DB tx Begin failed: ", tx.Error)
		return tx.Error
	}

	// insert
	err := tx.CreateInBatches(&proposal, 200).Error

	// Rollback in case of error during insert
	if err != nil {
		tx.Rollback()
		log.Logger.Error.Println("Insert Proposal Create In Batches failed, Rollback DB: ", err)
		return err
	}

	tx.Commit()
	return nil

}

// Use DB for updates
func UpdateDBProposals(proposal Proposals) error {
	var proposalDB Proposals
	tx := DB.Begin()

	if tx.Error != nil {
		log.Logger.Error.Println("Update Proposals DB tx Begin failed: ", tx.Error)
	}

	err := DB.Model(&proposalDB).Where("chain_name = ? AND proposal_id = ?", proposal.ChainName, proposal.ProposalID).Updates(proposal).Error

	// Rollback in case of error during insert
	if err != nil {
		tx.Rollback()
		log.Logger.Error.Println("Update Send Tx failed, Rollback DB: ", err)
		return err
	}

	tx.Commit()
	return nil

}

// A function that allows you to find the index of the DB's information through chainName, propositionID when updating it to the DB
func DBFind(chainName, proposalID string) int {
	var index int
	for i, oldProposal := range ProposalsDB {
		if oldProposal.ProposalID == proposalID && oldProposal.ChainName == chainName {
			index = i
			break
		}
	}
	return index
}
