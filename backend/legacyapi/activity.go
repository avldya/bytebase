package api

import (
	"encoding/json"

	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/vcs"
)

// ActivityType is the type for an activity.
type ActivityType string

const (
	// Issue related.

	// ActivityIssueCreate is the type for creating issues.
	ActivityIssueCreate ActivityType = "bb.issue.create"
	// ActivityIssueCommentCreate is the type for creating issue comments.
	ActivityIssueCommentCreate ActivityType = "bb.issue.comment.create"
	// ActivityIssueFieldUpdate is the type for updating issue fields.
	ActivityIssueFieldUpdate ActivityType = "bb.issue.field.update"
	// ActivityIssueStatusUpdate is the type for updating issue status.
	ActivityIssueStatusUpdate ActivityType = "bb.issue.status.update"
	// ActivityPipelineStageStatusUpdate is the type for stage begins or ends.
	ActivityPipelineStageStatusUpdate ActivityType = "bb.pipeline.stage.status.update"
	// ActivityPipelineTaskStatusUpdate is the type for updating pipeline task status.
	ActivityPipelineTaskStatusUpdate ActivityType = "bb.pipeline.task.status.update"
	// ActivityPipelineTaskFileCommit is the type for committing pipeline task file.
	ActivityPipelineTaskFileCommit ActivityType = "bb.pipeline.task.file.commit"
	// ActivityPipelineTaskStatementUpdate is the type for updating pipeline task SQL statement.
	ActivityPipelineTaskStatementUpdate ActivityType = "bb.pipeline.task.statement.update"
	// ActivityPipelineTaskEarliestAllowedTimeUpdate is the type for updating pipeline task the earliest allowed time.
	ActivityPipelineTaskEarliestAllowedTimeUpdate ActivityType = "bb.pipeline.task.general.earliest-allowed-time.update"

	// Member related.

	// ActivityMemberCreate is the type for creating members.
	ActivityMemberCreate ActivityType = "bb.member.create"
	// ActivityMemberRoleUpdate is the type for updating member roles.
	ActivityMemberRoleUpdate ActivityType = "bb.member.role.update"
	// ActivityMemberActivate is the type for activating members.
	ActivityMemberActivate ActivityType = "bb.member.activate"
	// ActivityMemberDeactivate is the type for deactivating members.
	ActivityMemberDeactivate ActivityType = "bb.member.deactivate"

	// Project related.

	// ActivityProjectRepositoryPush is the type for pushing repositories.
	ActivityProjectRepositoryPush ActivityType = "bb.project.repository.push"
	// ActivityProjectDatabaseTransfer is the type for transferring databases.
	ActivityProjectDatabaseTransfer ActivityType = "bb.project.database.transfer"
	// ActivityProjectMemberCreate is the type for creating project members.
	ActivityProjectMemberCreate ActivityType = "bb.project.member.create"
	// ActivityProjectMemberDelete is the type for deleting project members.
	ActivityProjectMemberDelete ActivityType = "bb.project.member.delete"
	// ActivityProjectMemberRoleUpdate is the type for updating project member roles.
	ActivityProjectMemberRoleUpdate ActivityType = "bb.project.member.role.update"

	// SQL Editor related.

	// ActivitySQLEditorQuery is the type for executing query.
	ActivitySQLEditorQuery ActivityType = "bb.sql-editor.query"

	// Database related.

	// ActivityDatabaseRecoveryPITRDone is the type for performing PITR on the database successfully.
	ActivityDatabaseRecoveryPITRDone ActivityType = "bb.database.recovery.pitr.done"
)

// ActivityLevel is the level of activities.
type ActivityLevel string

const (
	// ActivityInfo is the INFO level of activities.
	ActivityInfo ActivityLevel = "INFO"
	// ActivityWarn is the WARN level of activities.
	ActivityWarn ActivityLevel = "WARN"
	// ActivityError is the ERROR level of activities.
	ActivityError ActivityLevel = "ERROR"
)

// ActivityIssueCreatePayload is the API message payloads for creating issues.
// These payload types are only used when marshalling to the json format for saving into the database.
// So we annotate with json tag using camelCase naming which is consistent with normal
// json naming convention. More importantly, frontend code can simply use JSON.parse to
// convert to the expected struct there.
type ActivityIssueCreatePayload struct {
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

// TaskRollbackBy records an issue rollback activity.
// The task with taskID in IssueID is rollbacked by the task with RollbackByTaskID in RollbackByIssueID.
type TaskRollbackBy struct {
	IssueID           int `json:"issueId"`
	TaskID            int `json:"taskId"`
	RollbackByIssueID int `json:"rollbackByIssueId"`
	RollbackByTaskID  int `json:"rollbackByTaskId"`
}

// ApprovalEvent duplicates store/approval.proto.
type ApprovalEvent struct {
	Status string `json:"status"`
}

// ActivityIssueCommentCreatePayload is the API message payloads for creating issue comments.
type ActivityIssueCommentCreatePayload struct {
	ExternalApprovalEvent *ExternalApprovalEvent `json:"externalApprovalEvent,omitempty"`

	TaskRollbackBy *TaskRollbackBy `json:"taskRollbackBy,omitempty"`

	ApprovalEvent *ApprovalEvent `json:"approvalEvent,omitempty"`

	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

// ActivityIssueFieldUpdatePayload is the API message payloads for updating issue fields.
type ActivityIssueFieldUpdatePayload struct {
	FieldID  IssueFieldID `json:"fieldId"`
	OldValue string       `json:"oldValue,omitempty"`
	NewValue string       `json:"newValue,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

// ActivityIssueStatusUpdatePayload is the API message payloads for updating issue status.
type ActivityIssueStatusUpdatePayload struct {
	OldStatus IssueStatus `json:"oldStatus,omitempty"`
	NewStatus IssueStatus `json:"newStatus,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

// ActivityPipelineStageStatusUpdatePayload is the API message payloads for stage status updates.
type ActivityPipelineStageStatusUpdatePayload struct {
	StageID               int                   `json:"stageId"`
	StageStatusUpdateType StageStatusUpdateType `json:"stageStatusUpdateType"`

	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
	StageName string `json:"stageName"`
}

// ActivityPipelineTaskStatusUpdatePayload is the API message payloads for updating pipeline task status.
type ActivityPipelineTaskStatusUpdatePayload struct {
	TaskID    int        `json:"taskId"`
	OldStatus TaskStatus `json:"oldStatus,omitempty"`
	NewStatus TaskStatus `json:"newStatus,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
	TaskName  string `json:"taskName"`
}

// ActivityPipelineTaskFileCommitPayload is the API message payloads for committing pipeline task files.
type ActivityPipelineTaskFileCommitPayload struct {
	TaskID             int    `json:"taskId"`
	VCSInstanceURL     string `json:"vcsInstanceUrl,omitempty"`
	RepositoryFullPath string `json:"repositoryFullPath,omitempty"`
	Branch             string `json:"branch,omitempty"`
	FilePath           string `json:"filePath,omitempty"`
	CommitID           string `json:"commitId,omitempty"`
}

// ActivityPipelineTaskStatementUpdatePayload is the API message payloads for pipeline task statement updates.
type ActivityPipelineTaskStatementUpdatePayload struct {
	TaskID       int    `json:"taskId"`
	OldStatement string `json:"oldStatement,omitempty"`
	NewStatement string `json:"newStatement,omitempty"`
	OldSheetID   int    `json:"oldSheetId,omitempty"`
	NewSheetID   int    `json:"newSheetId,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
	TaskName  string `json:"taskName"`
}

// ActivityPipelineTaskEarliestAllowedTimeUpdatePayload is the API message payloads for pipeline task the earliest allowed time updates.
type ActivityPipelineTaskEarliestAllowedTimeUpdatePayload struct {
	TaskID               int   `json:"taskId"`
	OldEarliestAllowedTs int64 `json:"oldEarliestAllowedTs,omitempty"`
	NewEarliestAllowedTs int64 `json:"newEarliestAllowedTs,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
	TaskName  string `json:"taskName"`
}

// ActivityMemberCreatePayload is the API message payloads for creating members.
type ActivityMemberCreatePayload struct {
	PrincipalID    int          `json:"principalId"`
	PrincipalName  string       `json:"principalName"`
	PrincipalEmail string       `json:"principalEmail"`
	MemberStatus   MemberStatus `json:"memberStatus"`
	Role           Role         `json:"role"`
}

// ActivityMemberRoleUpdatePayload is the API message payloads for updating member roles.
type ActivityMemberRoleUpdatePayload struct {
	PrincipalID    int    `json:"principalId"`
	PrincipalName  string `json:"principalName"`
	PrincipalEmail string `json:"principalEmail"`
	OldRole        Role   `json:"oldRole"`
	NewRole        Role   `json:"newRole"`
}

// ActivityMemberActivateDeactivatePayload is the API message payloads for activating or deactivating members.
type ActivityMemberActivateDeactivatePayload struct {
	PrincipalID    int    `json:"principalId"`
	PrincipalName  string `json:"principalName"`
	PrincipalEmail string `json:"principalEmail"`
	Role           Role   `json:"role"`
}

// ActivityProjectRepositoryPushPayload is the API message payloads for pushing repositories.
type ActivityProjectRepositoryPushPayload struct {
	VCSPushEvent vcs.PushEvent `json:"pushEvent"`
	// Used by activity table to display info without paying the join cost
	// IssueID/IssueName only exist if the push event leads to the issue creation.
	IssueID   int    `json:"issueId,omitempty"`
	IssueName string `json:"issueName,omitempty"`
}

// ActivityProjectDatabaseTransferPayload is the API message payloads for transferring databases.
type ActivityProjectDatabaseTransferPayload struct {
	DatabaseID int `json:"databaseId,omitempty"`
	// Used by activity table to display info without paying the join cost
	DatabaseName string `json:"databaseName,omitempty"`
}

// ActivitySQLEditorQueryPayload is the API message payloads for the executed query info.
type ActivitySQLEditorQueryPayload struct {
	// Used by activity table to display info without paying the join cost
	Statement  string `json:"statement"`
	DurationNs int64  `json:"durationNs"`
	InstanceID int    `json:"instanceId"`
	// DeprecatedInstanceName is deprecated and should be removed from future version.
	DeprecatedInstanceName string           `json:"instanceName"`
	DatabaseID             int              `json:"databaseId"`
	DatabaseName           string           `json:"databaseName"`
	Error                  string           `json:"error"`
	AdviceList             []advisor.Advice `json:"adviceList"`
}

// Activity is the API message for an activity.
type Activity struct {
	ID int `jsonapi:"primary,activity"`

	// Standard fields
	CreatorID int
	Creator   *Principal `jsonapi:"relation,creator"`
	CreatedTs int64      `jsonapi:"attr,createdTs"`
	UpdaterID int
	Updater   *Principal `jsonapi:"relation,updater"`
	UpdatedTs int64      `jsonapi:"attr,updatedTs"`

	// Related fields
	// The object where this activity belongs
	// e.g if Type is "bb.issue.xxx", then this field refers to the corresponding issue's id.
	ContainerID int `jsonapi:"attr,containerId"`

	// Domain specific fields
	Type    ActivityType  `jsonapi:"attr,type"`
	Level   ActivityLevel `jsonapi:"attr,level"`
	Comment string        `jsonapi:"attr,comment"`
	Payload string        `jsonapi:"attr,payload"`
}

// ActivityResponse is the API message for an activity response.
type ActivityResponse struct {
	ActivityList []*Activity `jsonapi:"relation,activityList"`
	NextToken    string      `jsonapi:"attr,nextToken"`
}

// ActivityCreate is the API message for creating an activity.
type ActivityCreate struct {
	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	CreatorID int

	// Domain specific fields
	ContainerID int          `jsonapi:"attr,containerId"`
	Type        ActivityType `jsonapi:"attr,type"`
	Level       ActivityLevel
	Comment     string `jsonapi:"attr,comment"`
	Payload     string `jsonapi:"attr,payload"`
}

// ActivityFind is the API message for finding activities.
type ActivityFind struct {
	ID *int

	// Domain specific fields
	CreatorID       *int
	TypePrefixList  []string
	LevelList       []ActivityLevel
	ContainerID     *int
	Limit           *int
	CreatedTsAfter  *int64
	CreatedTsBefore *int64
	// If specified, only find activities whose ID is smaller than SinceID.
	SinceID *int
	// If specified, sorts the returned list by created_ts in <<ORDER>>
	// Different use cases want different orders.
	// e.g. Issue activity list wants ASC, while view recent activity list wants DESC.
	Order *SortOrder
}

func (find *ActivityFind) String() string {
	str, err := json.Marshal(*find)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

// ActivityPatch is the API message for patching an activity.
type ActivityPatch struct {
	ID int `jsonapi:"primary,activityPatch"`

	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	UpdaterID int

	// Domain specific fields
	Comment *string `jsonapi:"attr,comment"`

	// TODO(d): remove the payload after the backfill.
	Payload *string
}
