<template>
  <div class="py-4 space-y-4">
    <ArchiveBanner v-if="instance.rowStatus == 'ARCHIVED'" />
    <BBAttention
      v-else-if="state.migrationSetupStatus != 'OK'"
      :style="'WARN'"
      :title="attentionTitle"
      :description="attentionText"
      :action-text="attentionActionText"
      @click-action="state.showCreateMigrationSchemaModal = true"
    />
    <div class="px-6 space-y-6">
      <InstanceForm :instance="instance" />
      <div
        v-if="hasDataSourceFeature"
        class="py-6 space-y-4 border-t divide-control-border"
      >
        <DataSourceTable :instance="instance" />
      </div>
      <div v-else>
        <div class="mb-4 flex items-center justify-between">
          <BBTabFilter
            :tab-item-list="tabItemList"
            :selected-index="state.selectedIndex"
            @select-index="
              (index: number) => {
                state.selectedIndex = index;
              }
            "
          />
          <div class="flex items-center space-x-4">
            <div>
              <BBSpin
                v-if="state.syncingSchema"
                :title="$t('instance.syncing')"
              />
            </div>
            <button
              v-if="allowEdit"
              :disabled="state.syncingSchema"
              type="button"
              class="btn-normal"
              @click.prevent="syncSchema"
            >
              {{ $t("common.sync-now") }}
            </button>
            <button
              v-if="
                instance.rowStatus == 'NORMAL' &&
                instanceHasCreateDatabase(instance)
              "
              type="button"
              class="btn-primary"
              @click.prevent="createDatabase"
            >
              {{ $t("instance.new-database") }}
            </button>
          </div>
        </div>
        <div v-if="state.selectedIndex == DATABASE_TAB">
          <div class="mb-2 pl-4">
            <BBCheckbox
              :title="$t('database.show-missing-databases')"
              :value="state.showMissingDatabases"
              @toggle="state.showMissingDatabases = $event"
            />
          </div>
          <DatabaseTable
            mode="INSTANCE"
            :scroll-on-page-change="false"
            :show-missing-databases="state.showMissingDatabases"
            :database-list="databaseList"
          />
        </div>
        <InstanceUserTable
          v-else-if="state.selectedIndex == USER_TAB"
          :instance-user-list="instanceUserList"
        />
      </div>
      <template v-if="allowArchiveOrRestore">
        <template v-if="instance.rowStatus == 'NORMAL'">
          <BBButtonConfirm
            :style="'ARCHIVE'"
            :button-text="$t('instance.archive-this-instance')"
            :ok-text="$t('common.archive')"
            :require-confirm="true"
            :confirm-title="
              $t('instance.archive-instance-instance-name', [instance.name])
            "
            :confirm-description="
              $t(
                'instance.archived-instances-will-not-be-shown-on-the-normal-interface-you-can-still-restore-later-from-the-archive-page'
              )
            "
            @confirm="doArchive"
          />
        </template>
        <template v-else-if="instance.rowStatus == 'ARCHIVED'">
          <BBButtonConfirm
            :style="'RESTORE'"
            :button-text="$t('instance.restore-this-instance')"
            :ok-text="$t('instance.restore')"
            :require-confirm="true"
            :confirm-title="
              $t('instance.restore-instance-instance-name-to-normal-state', [
                instance.name,
              ])
            "
            :confirm-description="''"
            @confirm="doRestore"
          />
        </template>
      </template>
    </div>
  </div>

  <BBAlert
    v-if="state.showCreateMigrationSchemaModal"
    :style="'INFO'"
    :ok-text="$t('common.create')"
    :title="$t('instance.create-migration-schema') + '?'"
    :description="
      $t(
        'instance.bytebase-relies-on-migration-schema-to-manage-gitops-based-schema-migration-for-databases-belonged-to-this-instance'
      )
    "
    :in-progress="state.creatingMigrationSchema"
    @ok="
      () => {
        doCreateMigrationSchema();
      }
    "
    @cancel="state.showCreateMigrationSchemaModal = false"
  ></BBAlert>

  <BBModal
    v-if="state.showCreateDatabaseModal"
    :title="$t('quick-action.create-db')"
    @close="state.showCreateDatabaseModal = false"
  >
    <CreateDatabasePrepForm
      :environment-id="instance.environment.id"
      :instance-id="instance.id"
      @dismiss="state.showCreateDatabaseModal = false"
    />
  </BBModal>
  <FeatureModal
    v-if="state.showFeatureModal"
    feature="bb.feature.instance-count"
    @cancel="state.showFeatureModal = false"
  />
</template>

<script lang="ts" setup>
import { computed, reactive, watchEffect } from "vue";
import {
  idFromSlug,
  hasWorkspacePermission,
  instanceHasCreateDatabase,
  isMemberOfProject,
} from "../utils";
import ArchiveBanner from "../components/ArchiveBanner.vue";
import DatabaseTable from "../components/DatabaseTable.vue";
import DataSourceTable from "../components/DataSourceTable.vue";
import InstanceUserTable from "../components/InstanceUserTable.vue";
import InstanceForm from "../components/InstanceForm.vue";
import CreateDatabasePrepForm from "../components/CreateDatabasePrepForm.vue";
import {
  Database,
  Instance,
  InstanceMigration,
  MigrationSchemaStatus,
  SQLResultSet,
} from "../types";
import { BBTabFilterItem } from "../bbkit/types";
import { useI18n } from "vue-i18n";
import {
  featureToRef,
  pushNotification,
  useCurrentUser,
  useDatabaseStore,
  useInstanceStore,
  useSubscriptionStore,
  useSQLStore,
  useDBSchemaStore,
} from "@/store";

const DATABASE_TAB = 0;
const USER_TAB = 1;

interface LocalState {
  selectedIndex: number;
  migrationSetupStatus: MigrationSchemaStatus;
  showCreateMigrationSchemaModal: boolean;
  creatingMigrationSchema: boolean;
  showCreateDatabaseModal: boolean;
  syncingSchema: boolean;
  showFeatureModal: boolean;
  showMissingDatabases: boolean;
}

const props = defineProps({
  instanceSlug: {
    required: true,
    type: String,
  },
});

const instanceStore = useInstanceStore();
const subscriptionStore = useSubscriptionStore();
const { t } = useI18n();

const currentUser = useCurrentUser();
const sqlStore = useSQLStore();

const state = reactive<LocalState>({
  selectedIndex: DATABASE_TAB,
  migrationSetupStatus: "OK",
  showCreateMigrationSchemaModal: false,
  creatingMigrationSchema: false,
  showCreateDatabaseModal: false,
  syncingSchema: false,
  showFeatureModal: false,
  showMissingDatabases: false,
});

const instance = computed((): Instance => {
  return instanceStore.getInstanceById(idFromSlug(props.instanceSlug));
});

const checkMigrationSetup = () => {
  instanceStore
    .checkMigrationSetup(instance.value.id)
    .then((migration: InstanceMigration) => {
      state.migrationSetupStatus = migration.status;
    });
};

const prepareMigrationSchemaStatus = () => {
  checkMigrationSetup();
};
watchEffect(prepareMigrationSchemaStatus);

const attentionTitle = computed((): string => {
  if (state.migrationSetupStatus == "NOT_EXIST") {
    return t("instance.missing-migration-schema");
  } else if (state.migrationSetupStatus == "UNKNOWN") {
    return t("instance.unable-to-connect-instance-to-check-migration-schema");
  }
  return "";
});

const attentionText = computed((): string => {
  if (state.migrationSetupStatus == "NOT_EXIST") {
    return (
      t(
        "instance.bytebase-relies-on-migration-schema-to-manage-gitops-based-schema-migration-for-databases-belonged-to-this-instance"
      ) +
      (hasWorkspacePermission(
        "bb.permission.workspace.manage-instance",
        currentUser.value.role
      )
        ? ""
        : " " + t("instance.please-contact-your-dba-to-configure-it"))
    );
  } else if (state.migrationSetupStatus == "UNKNOWN") {
    return (
      t(
        "instance.bytebase-relies-on-migration-schema-to-manage-gitops-based-schema-migration-for-databases-belonged-to-this-instance"
      ) +
      (hasWorkspacePermission(
        "bb.permission.workspace.manage-instance",
        currentUser.value.role
      )
        ? " " +
          t("instance.please-check-the-instance-connection-info-is-correct")
        : " " + t("instance.please-contact-your-dba-to-configure-it"))
    );
  }
  return "";
});

const attentionActionText = computed((): string => {
  if (
    hasWorkspacePermission(
      "bb.permission.workspace.manage-instance",
      currentUser.value.role
    )
  ) {
    if (state.migrationSetupStatus == "NOT_EXIST") {
      return t("instance.create-migration-schema");
    } else if (state.migrationSetupStatus == "UNKNOWN") {
      return "";
    }
  }
  return "";
});

const hasDataSourceFeature = featureToRef("bb.feature.data-source");

const databaseList = computed(() => {
  const list = useDatabaseStore().getDatabaseListByInstanceId(
    instance.value.id
  );

  if (
    hasWorkspacePermission(
      "bb.permission.workspace.manage-instance",
      currentUser.value.role
    )
  ) {
    return list;
  }

  // In edge case when the user is no longer an Owner or DBA, we only want to display the database
  // belonging to the project which the user is a member of. The returned list above may contain
  // databases not meeting this criteria and we need to filter out them.
  const filteredList: Database[] = [];
  for (const database of list) {
    if (isMemberOfProject(database.project, currentUser.value)) {
      filteredList.push(database);
    }
  }

  return filteredList;
});

const instanceUserList = computed(() => {
  return instanceStore.getInstanceUserListById(instance.value.id);
});

const allowEdit = computed(() => {
  return (
    instance.value.rowStatus == "NORMAL" &&
    hasWorkspacePermission(
      "bb.permission.workspace.manage-instance",
      currentUser.value.role
    )
  );
});

const allowArchiveOrRestore = computed(() => {
  return hasWorkspacePermission(
    "bb.permission.workspace.manage-instance",
    currentUser.value.role
  );
});

const tabItemList = computed((): BBTabFilterItem[] => {
  return [
    {
      title: t("common.databases"),
      alert: false,
    },
    {
      title: t("instance.users"),
      alert: false,
    },
  ];
});

const doArchive = () => {
  instanceStore
    .patchInstance({
      instanceId: instance.value.id,
      instancePatch: {
        rowStatus: "ARCHIVED",
      },
    })
    .then((updatedInstance) => {
      pushNotification({
        module: "bytebase",
        style: "SUCCESS",
        title: t(
          "instance.successfully-archived-instance-updatedinstance-name",
          [updatedInstance.name]
        ),
      });
    });
};

const doRestore = () => {
  const instanceList = instanceStore.getInstanceList(["NORMAL"]);
  if (subscriptionStore.instanceCount <= instanceList.length) {
    state.showFeatureModal = true;
    return;
  }
  instanceStore
    .patchInstance({
      instanceId: instance.value.id,
      instancePatch: {
        rowStatus: "NORMAL",
      },
    })
    .then((updatedInstance) => {
      pushNotification({
        module: "bytebase",
        style: "SUCCESS",
        title: t(
          "instance.successfully-restored-instance-updatedinstance-name",
          [updatedInstance.name]
        ),
      });
    });
};

const doCreateMigrationSchema = () => {
  state.creatingMigrationSchema = true;
  instanceStore
    .createMigrationSetup(instance.value.id)
    .then((resultSet: SQLResultSet) => {
      state.creatingMigrationSchema = false;
      if (resultSet.error) {
        pushNotification({
          module: "bytebase",
          style: "CRITICAL",
          title: t(
            "instance.failed-to-create-migration-schema-for-instance-instance-value-name",
            [instance.value.name]
          ),
          description: resultSet.error,
        });
      } else {
        checkMigrationSetup();
        pushNotification({
          module: "bytebase",
          style: "SUCCESS",
          title: t(
            "instance.successfully-created-migration-schema-for-instance-value-name",
            [instance.value.name]
          ),
        });
      }
      state.showCreateMigrationSchemaModal = false;
    });
};

const syncSchema = () => {
  state.syncingSchema = true;
  sqlStore
    .syncSchema(instance.value.id)
    .then((resultSet: SQLResultSet) => {
      state.syncingSchema = false;
      if (resultSet.error) {
        pushNotification({
          module: "bytebase",
          style: "CRITICAL",
          title: t(
            "instance.failed-to-sync-schema-for-instance-instance-value-name",
            [instance.value.name]
          ),
          description: resultSet.error,
        });
      } else {
        pushNotification({
          module: "bytebase",
          style: "SUCCESS",
          title: t(
            "instance.successfully-synced-schema-for-instance-instance-value-name",
            [instance.value.name]
          ),
          description: resultSet.error,
        });
      }

      // Clear the db schema metadata cache entities.
      // So we will re-fetch new values when needed.
      const dbSchemaStore = useDBSchemaStore();
      const databaseList = useDatabaseStore().getDatabaseListByInstanceId(
        instance.value.id
      );
      databaseList.forEach((database) =>
        dbSchemaStore.removeCacheByDatabaseId(database.id)
      );
    })
    .catch(() => {
      state.syncingSchema = false;
    });
};

const createDatabase = () => {
  state.showCreateDatabaseModal = true;
};
</script>
