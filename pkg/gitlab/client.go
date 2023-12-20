package gitlab

import (
	"github.com/fwiedmann/dispatch"
	"github.com/xanzy/go-gitlab"
)

func NewClient(url, token, projectId string, mergeRequestId int) (Client, error) {
	c, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return Client{}, err
	}
	return Client{
		gitLabApiClient: c,
	}, nil
}

type Client struct {
	gitLabApiClient *gitlab.Client
	projectId       string
	mergeRequestId  int
}

func (c *Client) Dispatch(notes []dispatch.Note) error {
	userIds, err := c.getUserIds(notes)
	if err != nil {
		return err
	}

	return c.updateMergeRequest(userIds, notes)
}

func (c *Client) getUserIds(notes []dispatch.Note) ([]int, error) {
	userIdsMap := map[int]any{}

	list, _, err := c.gitLabApiClient.ProjectMembers.ListAllProjectMembers(c.projectId, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		for _, member := range note.Members {
			for _, projectMember := range list {
				if member == projectMember.Email {
					userIdsMap[projectMember.ID] = struct{}{}
				}
			}
		}
	}

	userIds := make([]int, 0)
	for k := range userIdsMap {
		userIds = append(userIds, k)
	}
	return userIds, nil
}

func (c *Client) updateMergeRequest(memberIds []int, notes []dispatch.Note) error {
	mr, _, err := c.gitLabApiClient.MergeRequests.GetMergeRequest(c.projectId, c.mergeRequestId, nil)
	if err != nil {
		return err
	}

	updateOptions := &gitlab.UpdateMergeRequestOptions{}
	updateOptions.AssigneeIDs = c.addMissingAssignees(memberIds, mr.Assignees)

	_, _, err = c.gitLabApiClient.MergeRequests.UpdateMergeRequest(c.projectId, c.mergeRequestId, updateOptions)
	return err
}

func (c *Client) addMissingAssignees(memberIds []int, existingAssignees []*gitlab.BasicUser) *[]int {
	updatedAssignees := make([]int, 0)
	for _, assignee := range existingAssignees {
		updatedAssignees = append(updatedAssignees, assignee.ID)
	}
	for _, memberId := range memberIds {
		found := false
		for _, assigneeId := range updatedAssignees {
			if memberId == assigneeId {
				found = true
			}
		}
		if !found {
			updatedAssignees = append(updatedAssignees, memberId)
		}
	}
	return &updatedAssignees
}
