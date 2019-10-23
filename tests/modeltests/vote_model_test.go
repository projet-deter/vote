package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"github.com/m2fof/vote/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllVotes(t *testing.T) {

	err := refreshUserAndVoteTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = seedUsersAndVotes()
	if err != nil {
		log.Fatalf("Error seeding user and post  table %v\n", err)
	}
	votes, err := voteInstance.FindAllVotes(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*votes), 2)
}

func TestSavePost(t *testing.T) {

	err := refreshUserAndVoteTable()
	if err != nil {
		log.Fatalf("Error user and post refreshing table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	newVote := models.Vote{
		ID:       1,
		Title:    "This is the title",
		Desc:     "This is the content",
		AuthorID: user.ID,
	}
	savedVote, err := newVote.SaveVote(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	assert.Equal(t, newVote.ID, savedVote.ID)
	assert.Equal(t, newVote.Title, savedVote.Title)
	assert.Equal(t, newVote.Desc, savedVote.Desc)
	assert.Equal(t, newVote.AuthorID, savedVote.AuthorID)

}

func TestGetPostByID(t *testing.T) {

	err := refreshUserAndVoteTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	vote, err := seedOneUserAndOneVote()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundVote, err := voteInstance.FindVoteByID(server.DB, vote.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundVote.ID, vote.ID)
	assert.Equal(t, foundVote.Title, vote.Title)
	assert.Equal(t, foundVote.Desc, vote.Desc)
}

func TestUpdateAPost(t *testing.T) {

	err := refreshUserAndVoteTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	vote, err := seedOneUserAndOneVote()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	voteUpdate := models.Vote{
		ID:       1,
		Title:    "modiUpdate",
		Desc:     "modiupdate@gmail.com",
		AuthorID: vote.AuthorID,
	}
	updatedVote, err := voteUpdate.UpdateVote(server.DB, vote.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedVote.ID, voteUpdate.ID)
	assert.Equal(t, updatedVote.Title, voteUpdate.Title)
	assert.Equal(t, updatedVote.Desc, voteUpdate.Desc)
	assert.Equal(t, updatedVote.AuthorID, voteUpdate.AuthorID)
}

func TestDeleteAPost(t *testing.T) {

	err := refreshUserAndVoteTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table: %v\n", err)
	}
	post, err := seedOneUserAndOneVote()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	isDeleted, err := voteInstance.DeleteVote(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	//one shows that the record has been deleted or:
	// assert.Equal(t, int(isDeleted), 1)

	//Can be done this way too
	assert.Equal(t, isDeleted, int64(1))
}
