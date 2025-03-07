<template>
  <div
    class="flex-1 overflow-auto focus:outline-none"
    tabindex="0"
    v-bind="$attrs"
  >
    <main class="flex-1 relative overflow-y-auto">
      <!-- Highlight Panel -->
      <div
        class="px-4 pb-4 space-y-2 lg:space-y-0 lg:flex lg:items-center lg:justify-between"
      >
        <div class="flex-1 min-w-0 shrink-0">
          <!-- Summary -->
          <div class="flex items-center">
            <div>
              <div class="flex items-center">
                <h1
                  class="pt-2 pb-2.5 text-xl font-bold leading-6 text-main truncate flex items-center gap-x-3"
                >
                  {{ database.name }}

                  <ProductionEnvironmentIcon
                    :environment="database.instance.environment"
                    tooltip
                    class="w-5 h-5"
                  />

                  <BBBadge
                    v-if="isPITRDatabase(database)"
                    text="PITR"
                    :can-remove="false"
                    class="text-xs"
                  />
                </h1>
              </div>
            </div>
          </div>
          <dl
            class="flex flex-col space-y-1 md:space-y-0 md:flex-row md:flex-wrap"
            data-label="bb-database-detail-info-block"
          >
            <dt class="sr-only">{{ $t("common.environment") }}</dt>
            <dd class="flex items-center text-sm md:mr-4">
              <span class="textlabel"
                >{{ $t("common.environment") }}&nbsp;-&nbsp;</span
              >
              <router-link
                :to="`/environment/${environmentSlug(
                  database.instance.environment
                )}`"
                class="normal-link"
              >
                {{ environmentName(database.instance.environment) }}
              </router-link>
            </dd>
            <dt class="sr-only">{{ $t("common.instance") }}</dt>
            <dd class="flex items-center text-sm md:mr-4">
              <InstanceEngineIcon :instance="database.instance" />
              <span class="ml-1 textlabel"
                >{{ $t("common.instance") }}&nbsp;-&nbsp;</span
              >
              <router-link
                :to="`/instance/${instanceSlug(database.instance)}`"
                class="normal-link"
                >{{ instanceName(database.instance) }}</router-link
              >
            </dd>
            <dt class="sr-only">{{ $t("common.project") }}</dt>
            <dd class="flex items-center text-sm md:mr-4">
              <span class="textlabel"
                >{{ $t("common.project") }}&nbsp;-&nbsp;</span
              >
              <router-link
                :to="`/project/${projectSlug(database.project)}#databases`"
                class="normal-link"
                >{{ projectName(database.project) }}</router-link
              >
            </dd>
            <SQLEditorButton
              class="text-sm md:mr-4"
              :database="database"
              :label="true"
              :disabled="!allowQuery"
              @failed="handleGotoSQLEditorFailed"
            />
            <dd
              v-if="hasSchemaDiagramFeature"
              class="flex items-center text-sm md:mr-4 textlabel cursor-pointer hover:text-accent"
              @click.prevent="state.showSchemaDiagram = true"
            >
              <span class="mr-1">{{ $t("schema-diagram.self") }}</span>
              <SchemaDiagramIcon />
            </dd>
            <DatabaseLabelProps
              :label-list="database.labels"
              :database="database"
              :allow-edit="allowEditDatabaseLabels"
              @update:label-list="updateLabels"
            >
              <template #label="{ label }">
                <span class="textlabel capitalize">
                  {{ hidePrefix(label) }}&nbsp;-&nbsp;
                </span>
              </template>
            </DatabaseLabelProps>
          </dl>
        </div>
        <div
          class="flex flex-row justify-end items-center flex-wrap shrink gap-x-2 gap-y-2"
          data-label="bb-database-detail-action-buttons-container"
        >
          <BBSpin v-if="state.syncingSchema" :title="$t('instance.syncing')" />
          <button
            type="button"
            class="btn-normal"
            :disabled="state.syncingSchema"
            @click.prevent="syncDatabaseSchema"
          >
            {{ $t("common.sync-now") }}
          </button>
          <button
            v-if="allowTransferProject"
            type="button"
            class="btn-normal"
            @click.prevent="tryTransferProject"
          >
            <span>{{ $t("database.transfer-project") }}</span>
            <heroicons-outline:switch-horizontal
              class="-mr-1 ml-2 h-5 w-5 text-control-light"
            />
          </button>
          <button
            v-if="allowAlterSchemaOrChangeData"
            type="button"
            class="btn-normal"
            @click="createMigration('bb.issue.database.data.update')"
          >
            <span>{{ $t("database.change-data") }}</span>
          </button>
          <button
            v-if="allowAlterSchema"
            type="button"
            class="btn-normal"
            @click="createMigration('bb.issue.database.schema.update')"
          >
            <span>{{ $t("database.alter-schema") }}</span>
          </button>
        </div>
      </div>
    </main>

    <BBTabFilter
      class="px-3 pb-2 border-b border-block-border"
      :responsive="false"
      :tab-item-list="tabItemList"
      :selected-index="state.selectedIndex"
      data-label="bb-database-detail-tab"
      @select-index="
        (index: number) => {
          selectTab(index);
        }
      "
    />
    <div class="py-6 px-6">
      <template v-if="selectedTabItem?.hash === 'overview'">
        <DatabaseOverviewPanel :database="database" />
      </template>
      <template v-if="selectedTabItem?.hash === 'change-history'">
        <DatabaseMigrationHistoryPanel
          :database="database"
          :allow-edit="allowEdit"
        />
      </template>
      <template v-if="selectedTabItem?.hash === 'backup-and-restore'">
        <DatabaseBackupPanel
          :database="database"
          :allow-admin="allowAdmin"
          :allow-edit="allowEdit"
        />
      </template>
      <template v-if="selectedTabItem?.hash === 'slow-query'">
        <DatabaseSlowQueryPanel :database="database" />
      </template>
      <template v-if="selectedTabItem?.hash === 'settings'">
        <DatabaseSettingsPanel :database="database" />
      </template>
    </div>

    <BBModal
      v-if="state.showTransferDatabaseModal"
      :title="$t('database.transfer-project')"
      @close="state.showTransferDatabaseModal = false"
    >
      <div class="w-112 flex flex-col items-center">
        <div class="col-span-1 w-64">
          <label for="user" class="textlabel">{{ $t("common.project") }}</label>
          <!-- Only allow to transfer database to the project having OWNER role -->
          <ProjectSelect
            id="project"
            class="mt-1"
            name="project"
            :allowed-role-list="['OWNER']"
            :include-default-project="allowTransferToDefaultProject"
            :selected-id="state.editingProjectId as number"
            @select-project-id="
              (projectId) => {
                state.editingProjectId = projectId;
              }
            "
          />
        </div>
        <SelectDatabaseLabel
          :database="database"
          :target-project-id="state.editingProjectId"
          class="mt-4"
          @next="doTransfer"
        >
          <template #buttons="{ next }">
            <div
              class="w-full pt-4 mt-6 flex justify-end border-t border-block-border"
            >
              <button
                type="button"
                class="btn-normal py-2 px-4"
                @click.prevent="state.showTransferDatabaseModal = false"
              >
                {{ $t("common.cancel") }}
              </button>
              <!--
                We are not allowed to transfer a db either its labels are not valid
                or transferring into its project itself.
              -->
              <button
                type="button"
                class="btn-primary ml-3 inline-flex justify-center py-2 px-4"
                :disabled="state.editingProjectId == database.project.id"
                @click.prevent="next"
              >
                {{ $t("common.transfer") }}
              </button>
            </div>
          </template>
        </SelectDatabaseLabel>
      </div>
    </BBModal>
    <BBModal
      v-if="state.showIncorrectProjectModal"
      :title="$t('common.warning')"
      @close="state.showIncorrectProjectModal = false"
    >
      <div class="col-span-1 w-96">
        {{ $t("database.incorrect-project-warning") }}
      </div>
      <div class="pt-6 flex justify-end">
        <button
          type="button"
          class="btn-normal py-2 px-4"
          @click.prevent="state.showIncorrectProjectModal = false"
        >
          {{ $t("common.cancel") }}
        </button>
        <button
          type="button"
          class="btn-primary ml-3 inline-flex justify-center py-2 px-4"
          @click.prevent="
            state.showIncorrectProjectModal = false;
            state.showTransferDatabaseModal = true;
          "
        >
          {{ $t("database.go-to-transfer") }}
        </button>
      </div>
    </BBModal>
  </div>

  <GhostDialog ref="ghostDialog" />

  <BBModal
    v-if="state.showSchemaDiagram"
    :title="$t('schema-diagram.self')"
    class="h-[calc(100vh-40px)] !max-h-[calc(100vh-40px)]"
    header-class="!border-0"
    container-class="flex-1 !pt-0"
    @close="state.showSchemaDiagram = false"
  >
    <div class="w-[80vw] h-full">
      <SchemaDiagram
        :database="database"
        :database-metadata="
          dbSchemaStore.getDatabaseMetadataByDatabaseId(database.id)
        "
      />
    </div>
  </BBModal>

  <SchemaEditorModal
    v-if="state.showSchemaEditorModal"
    :database-id-list="[database.id]"
    alter-type="SINGLE_DB"
    @close="state.showSchemaEditorModal = false"
  />
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, watch, ref } from "vue";
import { useRouter } from "vue-router";
import dayjs from "dayjs";
import { useI18n } from "vue-i18n";
import { startCase } from "lodash-es";

import ProjectSelect from "@/components/ProjectSelect.vue";
import DatabaseBackupPanel from "@/components/DatabaseBackupPanel.vue";
import DatabaseMigrationHistoryPanel from "@/components/DatabaseMigrationHistoryPanel.vue";
import DatabaseOverviewPanel from "@/components/DatabaseOverviewPanel.vue";
import DatabaseSlowQueryPanel from "@/components/DatabaseSlowQueryPanel.vue";
import { DatabaseSettingsPanel } from "@/components/DatabaseDetail";
import InstanceEngineIcon from "@/components/InstanceEngineIcon.vue";
import { DatabaseLabelProps } from "@/components/DatabaseLabels";
import { SelectDatabaseLabel } from "@/components/TransferDatabaseForm";
import {
  idFromSlug,
  hasWorkspacePermission,
  hidePrefix,
  allowGhostMigration,
  isPITRDatabase,
  isDatabaseAccessible,
  allowUsingSchemaEditor,
  isArchivedDatabase,
  instanceHasBackupRestore,
  instanceHasAlterSchema,
  instanceSupportSlowQuery,
  hasPermissionInProject,
} from "@/utils";
import {
  ProjectId,
  UNKNOWN_ID,
  DEFAULT_PROJECT_ID,
  Database,
  DatabaseLabel,
  SQLResultSet,
} from "@/types";
import { BBTabFilterItem } from "@/bbkit/types";
import { GhostDialog } from "@/components/AlterSchemaPrepForm";
import { SchemaDiagram, SchemaDiagramIcon } from "@/components/SchemaDiagram";
import { SQLEditorButton } from "@/components/DatabaseDetail";
import {
  pushNotification,
  useCurrentUser,
  useDatabaseStore,
  useDBSchemaStore,
  usePolicyByDatabaseAndType,
  useSQLStore,
} from "@/store";

type DatabaseTabItem = {
  name: string;
  hash: string;
};

interface LocalState {
  showTransferDatabaseModal: boolean;
  showIncorrectProjectModal: boolean;
  showSchemaEditorModal: boolean;
  editingProjectId: ProjectId;
  selectedIndex: number;
  syncingSchema: boolean;
  showSchemaDiagram: boolean;
}

const props = defineProps({
  databaseSlug: {
    required: true,
    type: String,
  },
});

const { t } = useI18n();
const router = useRouter();
const databaseStore = useDatabaseStore();
const dbSchemaStore = useDBSchemaStore();
const sqlStore = useSQLStore();
const ghostDialog = ref<InstanceType<typeof GhostDialog>>();

const databaseTabItemList = computed((): DatabaseTabItem[] => {
  return [
    { name: t("common.overview"), hash: "overview" },
    { name: t("change-history.self"), hash: "change-history" },
    { name: t("common.backup-and-restore"), hash: "backup-and-restore" },
    { name: startCase(t("slow-query.slow-queries")), hash: "slow-query" },
    { name: t("common.settings"), hash: "settings" },
  ];
});

const state = reactive<LocalState>({
  showTransferDatabaseModal: false,
  showIncorrectProjectModal: false,
  showSchemaEditorModal: false,
  editingProjectId: UNKNOWN_ID,
  selectedIndex: 0,
  syncingSchema: false,
  showSchemaDiagram: false,
});

const currentUser = useCurrentUser();

const database = computed((): Database => {
  return databaseStore.getDatabaseById(idFromSlug(props.databaseSlug));
});

const hasSchemaDiagramFeature = computed((): boolean => {
  return instanceHasAlterSchema(database.value.instance);
});

const accessControlPolicy = usePolicyByDatabaseAndType(
  computed(() => ({
    databaseId: database.value.id,
    type: "bb.policy.access-control",
  }))
);
const allowQuery = computed(() => {
  const policy = accessControlPolicy.value;
  const list = policy ? [policy] : [];
  return isDatabaseAccessible(database.value, list, currentUser.value);
});

// Project can be transferred if meets either of the condition below:
// - Database is in default project
// - Workspace role can manage instance
// - Project role can transfer database
const allowTransferProject = computed(() => {
  if (isArchivedDatabase(database.value)) {
    return false;
  }

  if (database.value.project.id == DEFAULT_PROJECT_ID) {
    return true;
  }

  if (
    hasWorkspacePermission(
      "bb.permission.workspace.manage-project",
      currentUser.value.role
    )
  ) {
    return true;
  }

  if (
    hasPermissionInProject(
      database.value.project,
      currentUser.value,
      "bb.permission.project.transfer-database"
    )
  ) {
    return true;
  }

  return false;
});

const allowTransferToDefaultProject = computed(() => {
  if (database.value.project.id === DEFAULT_PROJECT_ID) {
    return true;
  }

  // Allow to transfer a database to DEFAULT project only if the current user
  // can manage all projects.
  // AKA DBA or workspace owner.
  return hasWorkspacePermission(
    "bb.permission.workspace.manage-project",
    currentUser.value.role
  );
});

// Database can be admined if meets either of the condition below:
// - Workspace role can manage instance
// - Project role can admin database
//
// The admin operation includes
// - Edit database label
// - Enable/disable backup
const allowAdmin = computed(() => {
  if (isArchivedDatabase(database.value)) {
    return false;
  }

  if (
    hasWorkspacePermission(
      "bb.permission.workspace.manage-instance",
      currentUser.value.role
    )
  ) {
    return true;
  }

  if (
    hasPermissionInProject(
      database.value.project,
      currentUser.value,
      "bb.permission.project.admin-database"
    )
  ) {
    return true;
  }
  return false;
});

// Database can be edited if meets either of the condition below:
// - Workspace role can manage instance
// - Project role can change database
//
// The edit operation includes
// - Take manual backup
const allowEdit = computed(() => {
  if (isArchivedDatabase(database.value)) {
    return false;
  }

  if (
    hasWorkspacePermission(
      "bb.permission.workspace.manage-instance",
      currentUser.value.role
    )
  ) {
    return true;
  }

  if (
    hasPermissionInProject(
      database.value.project,
      currentUser.value,
      "bb.permission.project.change-database"
    )
  ) {
    return true;
  }
  return false;
});

const allowAlterSchemaOrChangeData = computed(() => {
  if (database.value.project.id === DEFAULT_PROJECT_ID) {
    return false;
  }
  return allowEdit.value;
});

const allowAlterSchema = computed(() => {
  return (
    allowAlterSchemaOrChangeData.value &&
    instanceHasAlterSchema(database.value.instance)
  );
});

const allowEditDatabaseLabels = computed((): boolean => {
  // only allowed to edit database labels when allowAdmin
  return allowAdmin.value;
});

const availableDatabaseTabItemList = computed(() => {
  const db = database.value;
  return databaseTabItemList.value.filter((item) => {
    if (item.hash === "backup-and-restore") {
      return instanceHasBackupRestore(db.instance);
    }
    if (item.hash === "slow-query") {
      return instanceSupportSlowQuery(db.instance);
    }
    return true;
  });
});

const tabItemList = computed((): BBTabFilterItem[] => {
  return availableDatabaseTabItemList.value.map((item) => {
    return { title: item.name, alert: false };
  });
});

const tryTransferProject = () => {
  state.editingProjectId = database.value.project.id;
  state.showTransferDatabaseModal = true;
};

// 'normal' -> normal migration
// 'online' -> online migration
// false -> user clicked cancel button
const isUsingGhostMigration = async (databaseList: Database[]) => {
  if (database.value.project.tenantMode === "TENANT") {
    // Not available for tenant mode now.
    return "normal";
  }

  // check if all selected databases supports gh-ost
  if (allowGhostMigration(databaseList)) {
    // open the dialog to ask the user
    const { result, mode } = await ghostDialog.value!.open();
    if (!result) {
      return false; // return false when user clicked the cancel button
    }
    return mode;
  }

  // fallback to normal
  return "normal";
};

const createMigration = async (
  type: "bb.issue.database.schema.update" | "bb.issue.database.data.update"
) => {
  type AlterMode = "online" | "normal" | false;
  let mode: AlterMode = "normal";
  if (type === "bb.issue.database.schema.update") {
    if (
      database.value.syncStatus === "OK" &&
      allowUsingSchemaEditor([database.value])
    ) {
      state.showSchemaEditorModal = true;
      return;
    }

    // Check and show a normal/online selection modal dialog if needed.
    mode = await isUsingGhostMigration([database.value]);
  }
  if (mode === false) return;

  // Create a user friendly default issue name
  const issueNameParts: string[] = [];
  issueNameParts.push(`[${database.value.name}]`);
  if (mode === "online") {
    issueNameParts.push("Online schema change");
  } else {
    issueNameParts.push(
      type === "bb.issue.database.schema.update"
        ? `Alter schema`
        : `Change data`
    );
  }
  const datetime = dayjs().format("@MM-DD HH:mm");
  const tz = "UTC" + dayjs().format("ZZ");
  issueNameParts.push(`${datetime} ${tz}`);

  const query: Record<string, any> = {
    template: type,
    name: issueNameParts.join(" "),
    project: database.value.project.id,
    databaseList: database.value.id,
  };
  if (mode === "online") {
    query.ghost = "1";
  }

  router.push({
    name: "workspace.issue.detail",
    params: {
      issueSlug: "new",
    },
    query,
  });
};

const updateProject = (newProjectId: ProjectId, labels?: DatabaseLabel[]) => {
  databaseStore
    .transferProject({
      databaseId: database.value.id,
      projectId: newProjectId,
      labels,
    })
    .then((updatedDatabase) => {
      pushNotification({
        module: "bytebase",
        style: "SUCCESS",
        title: t(
          "database.successfully-transferred-updateddatabase-name-to-project-updateddatabase-project-name",
          [updatedDatabase.name, updatedDatabase.project.name]
        ),
      });
    });
};

const updateLabels = (labels: DatabaseLabel[]) => {
  databaseStore.patchDatabaseLabels({
    databaseId: database.value.id,
    labels,
  });
};

const selectedTabItem = computed(() => {
  return availableDatabaseTabItemList.value[state.selectedIndex];
});

const selectTab = (index: number) => {
  const item = availableDatabaseTabItemList.value[index];
  state.selectedIndex = index;
  router.replace({
    name: "workspace.database.detail",
    hash: "#" + item.hash,
  });
};

const selectDatabaseTabOnHash = () => {
  if (router.currentRoute.value.hash) {
    for (let i = 0; i < availableDatabaseTabItemList.value.length; i++) {
      if (
        availableDatabaseTabItemList.value[i].hash ==
        router.currentRoute.value.hash.slice(1)
      ) {
        selectTab(i);
        break;
      }
    }
  } else {
    selectTab(0);
  }
};

const handleGotoSQLEditorFailed = () => {
  state.editingProjectId = database.value.project.id;
  state.showIncorrectProjectModal = true;
};

onMounted(() => {
  selectDatabaseTabOnHash();
});

watch(
  () => router.currentRoute.value.hash,
  () => {
    if (router.currentRoute.value.name == "workspace.database.detail") {
      selectDatabaseTabOnHash();
    }
  }
);

const doTransfer = (labels: DatabaseLabel[]) => {
  updateProject(state.editingProjectId, labels);
  state.showTransferDatabaseModal = false;
};

const syncDatabaseSchema = () => {
  state.syncingSchema = true;
  sqlStore
    .syncDatabaseSchema(database.value.id)
    .then((resultSet: SQLResultSet) => {
      state.syncingSchema = false;
      if (resultSet.error) {
        pushNotification({
          module: "bytebase",
          style: "CRITICAL",
          title: t(
            "db.failed-to-sync-schema-for-database-database-value-name",
            [database.value.name]
          ),
          description: resultSet.error,
        });
      } else {
        pushNotification({
          module: "bytebase",
          style: "SUCCESS",
          title: t(
            "db.successfully-synced-schema-for-database-database-value-name",
            [database.value.name]
          ),
          description: resultSet.error,
        });
      }
      useDBSchemaStore().getOrFetchDatabaseMetadataById(
        database.value.id,
        true // skip cache
      );
    })
    .catch(() => {
      state.syncingSchema = false;
    });
};
</script>
