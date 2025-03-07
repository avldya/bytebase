import { Activity } from "./activity";
import { Anomaly } from "./anomaly";
import { BackupSetting } from "./backup";
import { Bookmark } from "./bookmark";
import { EMPTY_ID, UNKNOWN_ID } from "./const";
import { Database } from "./database";
import { DataSource } from "./dataSource";
import { Environment } from "./environment";
import { CommandId, CommandRegisterId, PrincipalId } from "./id";
import { Inbox } from "./inbox";
import { Instance } from "./instance";
import { Issue } from "./issue";
import { Member } from "./member";
import { Pipeline, Stage, Task, TaskProgress } from "./pipeline";
import { Principal } from "./principal";
import { Project, ProjectMember } from "./project";
import { ProjectWebhook } from "./projectWebhook";
import { Repository } from "./repository";
import { VCS } from "./vcs";
import { DeploymentConfig } from "./deployment";
import { Policy, DefaultApprovalPolicy } from "./policy";
import { Sheet } from "./sheet";
import { SQLReviewPolicy } from "./sqlReview";
import { AuditLog, AuditActivityType, AuditActivityLevel } from "./auditLog";

// System bot id
export const SYSTEM_BOT_ID = 1;

// The project to hold those databases synced from the instance but haven't been assigned an application
// project yet. We can't use UNKNOWN_ID because of referential integrity.
export const DEFAULT_PROJECT_ID = 1;

export const ALL_DATABASE_NAME = "*";

// For text input, we do validation if there is no further keystroke after 1s
export const TEXT_VALIDATION_DELAY = 1000;

// Normally, we poll issue every 10s to fetch any update from the server side.
// If change occurs, then we will start the poll from 0.2s, 0.4s, 0.8s, 1.6s, 3.2s, 5s, 10s, 10s ... with jitter
// We do this because new update is more likely to happen after the initial change (e.g task gets new update after changing its status)
export const NORMAL_POLL_INTERVAL = 10000;
export const MINIMUM_POLL_INTERVAL = 200;
// Add jitter to avoid timer from different clients converging to the same polling frequency.
export const POLL_JITTER = 200;

// It may take a while to perform instance related operations since we are
// connecting the remote instance. And certain operations just take longer for
// a particular database type due to its unique property (e.g. create migration schema
// is a heavier operation in TiDB than traditional RDBMS).
export const INSTANCE_OPERATION_TIMEOUT = 60000;

// RowStatus
export type RowStatus = "NORMAL" | "ARCHIVED";

// Router
export type RouterSlug = {
  principalId?: PrincipalId;
  environmentSlug?: string;
  projectSlug?: string;
  projectWebhookSlug?: string;
  issueSlug?: string;
  instanceSlug?: string;
  databaseSlug?: string;
  tableName?: string;
  dataSourceSlug?: string;
  migrationHistorySlug?: string;
  vcsSlug?: string;
  connectionSlug?: string;
  sheetSlug?: string;
  sqlReviewPolicySlug?: string;
  ssoName?: string;
};

// Quick Action Type
export type Command = {
  id: CommandId;
  registerId: CommandRegisterId;
  run: () => void;
};

export type EnvironmentQuickActionType =
  | "quickaction.bb.environment.create"
  | "quickaction.bb.environment.reorder";
export type ProjectQuickActionType =
  | "quickaction.bb.project.create"
  | "quickaction.bb.project.database.transfer";
export type InstanceQuickActionType = "quickaction.bb.instance.create";
export type UserQuickActionType = "quickaction.bb.user.manage";
export type DatabaseQuickActionType =
  | "quickaction.bb.database.create" // Used by DBA and Owner
  | "quickaction.bb.database.request" // Used by Developer
  | "quickaction.bb.database.schema.update"
  | "quickaction.bb.database.data.update"
  | "quickaction.bb.database.troubleshoot"
  | "quickaction.bb.database.schema.sync";

export type QuickActionType =
  | EnvironmentQuickActionType
  | ProjectQuickActionType
  | InstanceQuickActionType
  | UserQuickActionType
  | DatabaseQuickActionType;

export type ResourceType =
  | "PRINCIPAL"
  | "MEMBER"
  | "ENVIRONMENT"
  | "PROJECT"
  | "PROJECT_HOOK"
  | "PROJECT_MEMBER"
  | "INSTANCE"
  | "DATABASE"
  | "DATA_SOURCE"
  | "BACKUP_SETTING"
  | "ISSUE"
  | "PIPELINE"
  | "POLICY"
  | "STAGE"
  | "TASK_PROGRESS"
  | "TASK"
  | "ACTIVITY"
  | "INBOX"
  | "BOOKMARK"
  | "VCS"
  | "REPOSITORY"
  | "ANOMALY"
  | "DEPLOYMENT_CONFIG"
  | "SHEET"
  | "SQL_REVIEW"
  | "AUDIT_LOG";

interface ResourceMaker {
  (type: "PRINCIPAL"): Principal;
  (type: "MEMBER"): Member;
  (type: "ENVIRONMENT"): Environment;
  (type: "PROJECT"): Project;
  (type: "PROJECT_HOOK"): ProjectWebhook;
  (type: "PROJECT_MEMBER"): ProjectMember;
  (type: "INSTANCE"): Instance;
  (type: "DATABASE"): Database;
  (type: "DATA_SOURCE"): DataSource;
  (type: "BACKUP_SETTING"): BackupSetting;
  (type: "ISSUE"): Issue;
  (type: "PIPELINE"): Pipeline;
  (type: "POLICY"): Policy;
  (type: "STAGE"): Stage;
  (type: "TASK_PROGRESS"): TaskProgress;
  (type: "TASK"): Task;
  (type: "ACTIVITY"): Activity;
  (type: "INBOX"): Inbox;
  (type: "BOOKMARK"): Bookmark;
  (type: "VCS"): VCS;
  (type: "REPOSITORY"): Repository;
  (type: "ANOMALY"): Anomaly;
  (type: "DEPLOYMENT_CONFIG"): DeploymentConfig;
  (type: "SHEET"): Sheet;
  (type: "SQL_REVIEW"): SQLReviewPolicy;
  (type: "AUDIT_LOG"): AuditLog;
}

const makeUnknown = (type: ResourceType) => {
  // Have to omit creator and updater to avoid recursion.
  const UNKNOWN_PRINCIPAL: Principal = {
    id: UNKNOWN_ID,
    type: "END_USER",
    name: "<<Unknown principal>>",
    email: "",
    role: "DEVELOPER",
  } as Principal;

  const UNKNOWN_MEMBER: Member = {
    id: UNKNOWN_ID,
    rowStatus: "NORMAL",
    status: "ACTIVE",
    role: "DEVELOPER",
    principal: UNKNOWN_PRINCIPAL,
  };

  const UNKNOWN_ENVIRONMENT: Environment = {
    id: UNKNOWN_ID,
    resourceId: "",
    rowStatus: "NORMAL",
    name: "<<Unknown environment>>",
    order: 0,
    tier: "UNPROTECTED",
  };

  const UNKNOWN_PROJECT: Project = {
    id: UNKNOWN_ID,
    resourceId: "",
    rowStatus: "NORMAL",
    name: "<<Unknown project>>",
    key: "UNK",
    memberList: [],
    workflowType: "UI",
    visibility: "PUBLIC",
    tenantMode: "DISABLED",
    dbNameTemplate: "",
    schemaChangeType: "DDL",
  };

  const UNKNOWN_PROJECT_HOOK: ProjectWebhook = {
    id: UNKNOWN_ID,
    projectId: UNKNOWN_ID,
    type: "",
    name: "",
    url: "",
    activityList: [],
  };

  const UNKNOWN_PROJECT_MEMBER: ProjectMember = {
    id: `projects/${UNKNOWN_ID}/roles/${UNKNOWN_ID}/principals/${UNKNOWN_ID}`,
    project: UNKNOWN_PROJECT,
    role: "DEVELOPER",
    principal: UNKNOWN_PRINCIPAL,
  };

  const UNKNOWN_INSTANCE: Instance = {
    id: UNKNOWN_ID,
    resourceId: "",
    rowStatus: "NORMAL",
    environment: UNKNOWN_ENVIRONMENT,
    dataSourceList: [],
    name: "<<Unknown instance>>",
    engine: "MYSQL",
    engineVersion: "",
    externalLink: "",
    srv: false,
    authenticationDatabase: "",
  };

  const UNKNOWN_DATABASE: Database = {
    id: UNKNOWN_ID,
    instanceId: UNKNOWN_ID,
    instance: UNKNOWN_INSTANCE,
    projectId: UNKNOWN_ID,
    project: UNKNOWN_PROJECT,
    labels: [],
    dataSourceList: [],
    name: "<<Unknown database>>",
    characterSet: "",
    collation: "",
    syncStatus: "NOT_FOUND",
    lastSuccessfulSyncTs: 0,
    schemaVersion: "",
  };

  const UNKNOWN_DATA_SOURCE: DataSource = {
    id: UNKNOWN_ID,
    instanceId: UNKNOWN_ID,
    databaseId: UNKNOWN_ID,
    name: "<<Unknown data source>>",
    type: "RO",
    username: "",
    password: "",
    sslCa: "",
    sslCert: "",
    sslKey: "",
    host: "",
    port: "",
    database: "",
    options: {
      srv: false,
      authenticationDatabase: "",
      sid: "",
      serviceName: "",
    },
    // UI-only fields
    updateSsl: false,
  };

  const UNKNOWN_BACKUP_SETTING: BackupSetting = {
    id: UNKNOWN_ID,
    databaseId: UNKNOWN_ID,
    enabled: false,
    hour: 0,
    dayOfWeek: 0,
    hookUrl: "",
    retentionPeriodTs: 0,
  };

  const UNKNOWN_PIPELINE: Pipeline = {
    id: UNKNOWN_ID,
    name: "<<Unknown pipeline>>",
    stageList: [],
  };

  const UNKNOWN_POLICY: Policy = {
    id: UNKNOWN_ID,
    rowStatus: "NORMAL",
    resourceType: "",
    resourceId: UNKNOWN_ID,
    environment: UNKNOWN_ENVIRONMENT,
    type: "bb.policy.pipeline-approval",
    inheritFromParent: false,
    payload: {
      value: DefaultApprovalPolicy,
      assigneeGroupList: [],
    },
  };

  const UNKNOWN_ISSUE: Issue = {
    id: UNKNOWN_ID,
    project: UNKNOWN_PROJECT,
    pipeline: UNKNOWN_PIPELINE,
    creator: UNKNOWN_PRINCIPAL,
    createdTs: 0,
    updater: UNKNOWN_PRINCIPAL,
    updatedTs: 0,
    name: "<<Unknown issue>>",
    status: "DONE",
    type: "bb.issue.general",
    description: "",
    assignee: UNKNOWN_PRINCIPAL,
    assigneeNeedAttention: false,
    subscriberList: [],
    payload: {},
  };

  const UNKNOWN_STAGE: Stage = {
    id: UNKNOWN_ID,
    pipeline: UNKNOWN_PIPELINE,
    name: "<<Unknown stage>>",
    environment: UNKNOWN_ENVIRONMENT,
    taskList: [],
  };

  const UNKNOWN_TASK_PROGRESS: TaskProgress = {
    totalUnit: 0,
    completedUnit: 0,
    createdTs: 0,
    updatedTs: 0,
  };

  const UNKNOWN_TASK: Task = {
    id: UNKNOWN_ID,
    pipeline: UNKNOWN_PIPELINE,
    stage: UNKNOWN_STAGE,
    creator: UNKNOWN_PRINCIPAL,
    createdTs: 0,
    updater: UNKNOWN_PRINCIPAL,
    updatedTs: 0,
    name: "<<Unknown task>>",
    type: "bb.task.general",
    status: "DONE",
    instance: UNKNOWN_INSTANCE,
    database: UNKNOWN_DATABASE,
    earliestAllowedTs: 0,
    taskRunList: [],
    taskCheckRunList: [],
    blockedBy: [],
    progress: { ...UNKNOWN_TASK_PROGRESS },
  };

  const UNKNOWN_ACTIVITY: Activity = {
    id: UNKNOWN_ID,
    creator: UNKNOWN_PRINCIPAL,
    createdTs: 0,
    updater: UNKNOWN_PRINCIPAL,
    updatedTs: 0,
    containerId: UNKNOWN_ID,
    type: "bb.issue.create",
    level: "INFO",
    comment: "<<Unknown comment>>",
  };

  const UNKNOWN_INBOX: Inbox = {
    id: UNKNOWN_ID,
    receiver_id: UNKNOWN_ID,
    activity: UNKNOWN_ACTIVITY,
    status: "READ",
  };

  const UNKNOWN_BOOKMARK: Bookmark = {
    id: UNKNOWN_ID,
    creatorID: UNKNOWN_ID,
    name: "",
    link: "",
  };

  const UNKNOWN_VCS: VCS = {
    id: UNKNOWN_ID,
    name: "",
    type: "GITLAB",
    instanceUrl: "",
    apiUrl: "",
    applicationId: "",
    secret: "",
  };

  const UNKNOWN_REPOSITORY: Repository = {
    id: UNKNOWN_ID,
    vcs: UNKNOWN_VCS,
    project: UNKNOWN_PROJECT,
    name: "",
    fullPath: "",
    webUrl: "",
    baseDirectory: "",
    branchFilter: "",
    filePathTemplate: "",
    schemaPathTemplate: "",
    sheetPathTemplate: "",
    enableSQLReviewCI: false,
    sqlReviewCIPullRequestURL: "",
    externalId: UNKNOWN_ID.toString(),
  };

  const UNKNOWN_ANOMALY: Anomaly = {
    id: UNKNOWN_ID,
    creator: UNKNOWN_PRINCIPAL,
    createdTs: 0,
    updater: UNKNOWN_PRINCIPAL,
    updatedTs: 0,
    instanceId: UNKNOWN_ID,
    instance: UNKNOWN_INSTANCE,
    databaseId: UNKNOWN_ID,
    database: UNKNOWN_DATABASE,
    type: "bb.anomaly.database.backup.policy-violation",
    severity: "MEDIUM",
    payload: {
      environmentId: UNKNOWN_ID,
      expectedSchedule: "DAILY",
      actualSchedule: "UNSET",
    },
  };

  const UNKNOWN_DEPLOYMENT_CONFIG: DeploymentConfig = {
    id: UNKNOWN_ID,
    schedule: {
      deployments: [],
    },
  };

  const UNKNOWN_SHEET: Sheet = {
    id: UNKNOWN_ID,
    rowStatus: "NORMAL",
    creator: UNKNOWN_PRINCIPAL,
    creatorId: UNKNOWN_ID,
    createdTs: 0,
    updater: UNKNOWN_PRINCIPAL,
    updatedTs: 0,
    projectId: UNKNOWN_ID,
    project: UNKNOWN_PROJECT,
    databaseId: UNKNOWN_ID,
    database: UNKNOWN_DATABASE,
    name: "<<Unknown sheet>>",
    statement: "",
    visibility: "PRIVATE",
    source: "BYTEBASE",
    type: "SQL",
    starred: false,
    pinned: false,
    payload: {},
    size: 0,
  };

  const UNKNOWN_SQL_REVIEW_POLICY: SQLReviewPolicy = {
    id: UNKNOWN_ID,
    rowStatus: "NORMAL",
    environment: UNKNOWN_ENVIRONMENT,
    name: "",
    ruleList: [],
  };

  switch (type) {
    case "PRINCIPAL":
      return UNKNOWN_PRINCIPAL;
    case "MEMBER":
      return UNKNOWN_MEMBER;
    case "ENVIRONMENT":
      return UNKNOWN_ENVIRONMENT;
    case "PROJECT":
      return UNKNOWN_PROJECT;
    case "PROJECT_HOOK":
      return UNKNOWN_PROJECT_HOOK;
    case "PROJECT_MEMBER":
      return UNKNOWN_PROJECT_MEMBER;
    case "INSTANCE":
      return UNKNOWN_INSTANCE;
    case "DATABASE":
      return UNKNOWN_DATABASE;
    case "DATA_SOURCE":
      return UNKNOWN_DATA_SOURCE;
    case "BACKUP_SETTING":
      return UNKNOWN_BACKUP_SETTING;
    case "ISSUE":
      return UNKNOWN_ISSUE;
    case "PIPELINE":
      return UNKNOWN_PIPELINE;
    case "POLICY":
      return UNKNOWN_POLICY;
    case "STAGE":
      return UNKNOWN_STAGE;
    case "TASK_PROGRESS":
      return UNKNOWN_TASK_PROGRESS;
    case "TASK":
      return UNKNOWN_TASK;
    case "ACTIVITY":
      return UNKNOWN_ACTIVITY;
    case "INBOX":
      return UNKNOWN_INBOX;
    case "BOOKMARK":
      return UNKNOWN_BOOKMARK;
    case "VCS":
      return UNKNOWN_VCS;
    case "REPOSITORY":
      return UNKNOWN_REPOSITORY;
    case "ANOMALY":
      return UNKNOWN_ANOMALY;
    case "DEPLOYMENT_CONFIG":
      return UNKNOWN_DEPLOYMENT_CONFIG;
    case "SHEET":
      return UNKNOWN_SHEET;
    case "SQL_REVIEW":
      return UNKNOWN_SQL_REVIEW_POLICY;
  }
};
export const unknown = makeUnknown as ResourceMaker;

const makeEmpty = (type: ResourceType) => {
  // Have to omit creator and updater to avoid recursion.
  const EMPTY_PRINCIPAL: Principal = {
    id: EMPTY_ID,
    type: "END_USER",
    name: "",
    email: "",
    role: "DEVELOPER",
  } as Principal;

  const EMPTY_MEMBER: Member = {
    id: EMPTY_ID,
    rowStatus: "NORMAL",
    status: "ACTIVE",
    role: "DEVELOPER",
    principal: EMPTY_PRINCIPAL,
  };

  const EMPTY_ENVIRONMENT: Environment = {
    id: EMPTY_ID,
    resourceId: "",
    rowStatus: "NORMAL",
    name: "",
    order: 0,
    tier: "UNPROTECTED",
  };

  const EMPTY_PROJECT: Project = {
    id: EMPTY_ID,
    resourceId: "",
    rowStatus: "NORMAL",
    name: "",
    key: "",
    memberList: [],
    workflowType: "UI",
    visibility: "PUBLIC",
    tenantMode: "DISABLED",
    dbNameTemplate: "",
    schemaChangeType: "DDL",
  };

  const EMPTY_PROJECT_HOOK: ProjectWebhook = {
    id: EMPTY_ID,
    projectId: EMPTY_ID,
    type: "",
    name: "",
    url: "",
    activityList: [],
  };

  const EMPTY_PROJECT_MEMBER: ProjectMember = {
    id: `projects/${EMPTY_ID}/roles/${EMPTY_ID}/principals/${EMPTY_ID}`,
    project: EMPTY_PROJECT,
    role: "DEVELOPER",
    principal: EMPTY_PRINCIPAL,
  };

  const EMPTY_INSTANCE: Instance = {
    id: EMPTY_ID,
    resourceId: "",
    rowStatus: "NORMAL",
    environment: EMPTY_ENVIRONMENT,
    dataSourceList: [],
    name: "",
    engine: "MYSQL",
    engineVersion: "",
    externalLink: "",
    srv: false,
    authenticationDatabase: "",
  };

  const EMPTY_DATABASE: Database = {
    id: EMPTY_ID,
    instanceId: UNKNOWN_ID,
    instance: EMPTY_INSTANCE,
    projectId: UNKNOWN_ID,
    project: EMPTY_PROJECT,
    dataSourceList: [],
    labels: [],
    name: "",
    characterSet: "",
    collation: "",
    syncStatus: "NOT_FOUND",
    lastSuccessfulSyncTs: 0,
    schemaVersion: "",
  };

  const EMPTY_DATA_SOURCE: DataSource = {
    id: EMPTY_ID,
    instanceId: UNKNOWN_ID,
    databaseId: UNKNOWN_ID,
    name: "",
    type: "RO",
    username: "",
    password: "",
    sslCa: "",
    sslCert: "",
    sslKey: "",
    host: "",
    port: "",
    database: "",
    options: { srv: false, authenticationDatabase: "" },
    // UI-only fields
    updateSsl: false,
    sid: "",
    serviceName: "",
  };

  const EMPTY_BACKUP_SETTING: BackupSetting = {
    id: EMPTY_ID,
    databaseId: UNKNOWN_ID,
    enabled: false,
    hour: 0,
    dayOfWeek: 0,
    hookUrl: "",
    retentionPeriodTs: 0,
  };

  const EMPTY_PIPELINE: Pipeline = {
    id: EMPTY_ID,
    name: "",
    stageList: [],
  };

  const EMPTY_POLICY: Policy = {
    id: EMPTY_ID,
    rowStatus: "NORMAL",
    resourceType: "",
    resourceId: EMPTY_ID,
    environment: EMPTY_ENVIRONMENT,
    type: "bb.policy.pipeline-approval",
    inheritFromParent: false,
    payload: {
      value: DefaultApprovalPolicy,
      assigneeGroupList: [],
    },
  };

  const EMPTY_ISSUE: Issue = {
    id: EMPTY_ID,
    pipeline: EMPTY_PIPELINE,
    project: EMPTY_PROJECT,
    creator: EMPTY_PRINCIPAL,
    createdTs: 0,
    updater: EMPTY_PRINCIPAL,
    updatedTs: 0,
    name: "",
    status: "DONE",
    type: "bb.issue.general",
    description: "",
    assignee: EMPTY_PRINCIPAL,
    assigneeNeedAttention: false,
    subscriberList: [],
    payload: {},
  };

  const EMPTY_STAGE: Stage = {
    id: EMPTY_ID,
    pipeline: EMPTY_PIPELINE,
    name: "",
    environment: EMPTY_ENVIRONMENT,
    taskList: [],
  };

  const EMPTY_TASK_PROGRESS: TaskProgress = {
    totalUnit: 0,
    completedUnit: 0,
    createdTs: 0,
    updatedTs: 0,
  };

  const EMPTY_TASK: Task = {
    id: EMPTY_ID,
    pipeline: EMPTY_PIPELINE,
    stage: EMPTY_STAGE,
    creator: EMPTY_PRINCIPAL,
    createdTs: 0,
    updater: EMPTY_PRINCIPAL,
    updatedTs: 0,
    name: "",
    type: "bb.task.general",
    status: "DONE",
    instance: EMPTY_INSTANCE,
    database: EMPTY_DATABASE,
    taskRunList: [],
    taskCheckRunList: [],
    earliestAllowedTs: 0,
    blockedBy: [],
    progress: { ...EMPTY_TASK_PROGRESS },
  };

  const EMPTY_ACTIVITY: Activity = {
    id: EMPTY_ID,
    creator: EMPTY_PRINCIPAL,
    createdTs: 0,
    updater: EMPTY_PRINCIPAL,
    updatedTs: 0,
    containerId: EMPTY_ID,
    type: "bb.issue.create",
    level: "INFO",
    comment: "",
  };

  const EMPTY_INBOX: Inbox = {
    id: EMPTY_ID,
    receiver_id: EMPTY_ID,
    activity: EMPTY_ACTIVITY,
    status: "READ",
  };

  const EMPTY_BOOKMARK: Bookmark = {
    id: EMPTY_ID,
    creatorID: EMPTY_ID,
    name: "",
    link: "",
  };

  const EMPTY_VCS: VCS = {
    id: EMPTY_ID,
    name: "",
    type: "GITLAB",
    instanceUrl: "",
    apiUrl: "",
    applicationId: "",
    secret: "",
  };

  const EMPTY_REPOSITORY: Repository = {
    id: EMPTY_ID,
    vcs: EMPTY_VCS,
    project: EMPTY_PROJECT,
    name: "",
    fullPath: "",
    webUrl: "",
    baseDirectory: "",
    branchFilter: "",
    filePathTemplate: "",
    schemaPathTemplate: "",
    sheetPathTemplate: "",
    enableSQLReviewCI: false,
    sqlReviewCIPullRequestURL: "",
    externalId: EMPTY_ID.toString(),
  };

  const EMPTY_ANOMALY: Anomaly = {
    id: EMPTY_ID,
    creator: EMPTY_PRINCIPAL,
    createdTs: 0,
    updater: EMPTY_PRINCIPAL,
    updatedTs: 0,
    instanceId: EMPTY_ID,
    instance: EMPTY_INSTANCE,
    databaseId: EMPTY_ID,
    database: EMPTY_DATABASE,
    type: "bb.anomaly.database.backup.policy-violation",
    severity: "MEDIUM",
    payload: {
      environmentId: EMPTY_ID,
      expectedSchedule: "DAILY",
      actualSchedule: "UNSET",
    },
  };

  const EMPTY_DEPLOYMENT_CONFIG: DeploymentConfig = {
    id: EMPTY_ID,
    schedule: {
      deployments: [],
    },
  };

  const EMPTY_SHEET: Sheet = {
    id: EMPTY_ID,
    rowStatus: "NORMAL",
    creator: EMPTY_PRINCIPAL,
    creatorId: EMPTY_ID,
    createdTs: 0,
    updater: EMPTY_PRINCIPAL,
    updatedTs: 0,
    projectId: EMPTY_ID,
    project: EMPTY_PROJECT,
    databaseId: EMPTY_ID,
    database: EMPTY_DATABASE,
    name: "<<Empty sheet>>",
    statement: "",
    visibility: "PRIVATE",
    source: "BYTEBASE",
    type: "SQL",
    starred: false,
    pinned: false,
    payload: {},
    size: 0,
  };

  const EMPTY_SQL_REVIEW_POLICY: SQLReviewPolicy = {
    id: EMPTY_ID,
    rowStatus: "NORMAL",
    environment: EMPTY_ENVIRONMENT,
    name: "",
    ruleList: [],
  };

  const EMPTY_AUDIT_LOG: AuditLog = {
    createdTs: 0,
    creator: EMPTY_PRINCIPAL.email,
    type: AuditActivityType.WorkspaceMemberCreate,
    level: AuditActivityLevel.INFO,
    comment: "",
    payload: "",
  };

  switch (type) {
    case "PRINCIPAL":
      return EMPTY_PRINCIPAL;
    case "MEMBER":
      return EMPTY_MEMBER;
    case "ENVIRONMENT":
      return EMPTY_ENVIRONMENT;
    case "PROJECT":
      return EMPTY_PROJECT;
    case "PROJECT_HOOK":
      return EMPTY_PROJECT_HOOK;
    case "PROJECT_MEMBER":
      return EMPTY_PROJECT_MEMBER;
    case "INSTANCE":
      return EMPTY_INSTANCE;
    case "DATABASE":
      return EMPTY_DATABASE;
    case "DATA_SOURCE":
      return EMPTY_DATA_SOURCE;
    case "BACKUP_SETTING":
      return EMPTY_BACKUP_SETTING;
    case "ISSUE":
      return EMPTY_ISSUE;
    case "PIPELINE":
      return EMPTY_PIPELINE;
    case "POLICY":
      return EMPTY_POLICY;
    case "STAGE":
      return EMPTY_STAGE;
    case "TASK_PROGRESS":
      return EMPTY_TASK_PROGRESS;
    case "TASK":
      return EMPTY_TASK;
    case "ACTIVITY":
      return EMPTY_ACTIVITY;
    case "INBOX":
      return EMPTY_INBOX;
    case "BOOKMARK":
      return EMPTY_BOOKMARK;
    case "VCS":
      return EMPTY_VCS;
    case "REPOSITORY":
      return EMPTY_REPOSITORY;
    case "ANOMALY":
      return EMPTY_ANOMALY;
    case "DEPLOYMENT_CONFIG":
      return EMPTY_DEPLOYMENT_CONFIG;
    case "SHEET":
      return EMPTY_SHEET;
    case "SQL_REVIEW":
      return EMPTY_SQL_REVIEW_POLICY;
    case "AUDIT_LOG":
      return EMPTY_AUDIT_LOG;
  }
};
export const empty = makeEmpty as ResourceMaker;
