<template>
  <div class="space-y-2">
    <div
      class="text-lg font-medium leading-7 text-main flex items-center justify-between"
    >
      <div class="flex items-center">
        <EnvironmentTabFilter
          :environment="state.environment"
          :include-all="true"
          @update:environment="state.environment = $event ?? UNKNOWN_ID"
        />
      </div>
      <NInputGroup style="width: auto">
        <InstanceSelect
          :instance="state.instance"
          :include-all="true"
          :filter="filterInstance"
          :environment="state.environment"
          @update:instance="state.instance = $event ?? UNKNOWN_ID"
        />
        <SearchBox
          :value="state.keyword"
          :placeholder="$t('database.search-database')"
          @update:value="state.keyword = $event"
        />
      </NInputGroup>
    </div>

    <template v-if="databaseList.length > 0">
      <DatabaseTable
        mode="PROJECT"
        table-class="border"
        :database-list="filteredDatabaseList"
      />
    </template>
    <div v-else class="text-center textinfolabel">
      <i18n-t keypath="project.overview.no-db-prompt" tag="p">
        <template #newDb>
          <span class="text-main">{{ $t("quick-action.new-db") }}</span>
        </template>
        <template #transferInDb>
          <span class="text-main">{{ $t("quick-action.transfer-in-db") }}</span>
        </template>
      </i18n-t>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive, PropType, computed } from "vue";
import { NInputGroup } from "naive-ui";
import { uniqBy } from "lodash-es";

import {
  Database,
  EnvironmentId,
  Instance,
  InstanceId,
  Project,
  UNKNOWN_ID,
} from "../types";
import { filterDatabaseByKeyword } from "@/utils";
import DatabaseTable from "../components/DatabaseTable.vue";
import { EnvironmentTabFilter, InstanceSelect, SearchBox } from "./v2";

interface LocalState {
  environment: EnvironmentId;
  instance: InstanceId;
  keyword: string;
}

const props = defineProps({
  project: {
    required: true,
    type: Object as PropType<Project>,
  },
  databaseList: {
    required: true,
    type: Object as PropType<Database[]>,
  },
});

const state = reactive<LocalState>({
  environment: UNKNOWN_ID,
  instance: UNKNOWN_ID,
  keyword: "",
});

const filteredDatabaseList = computed(() => {
  return props.databaseList
    .filter((db) => {
      return (
        state.environment === UNKNOWN_ID ||
        db.instance.environment.id === state.environment
      );
    })
    .filter((db) => {
      return state.instance === UNKNOWN_ID || db.instance.id === state.instance;
    })
    .filter((db) => {
      return filterDatabaseByKeyword(db, state.keyword, [
        "name",
        "environment",
        "instance",
      ]);
    });
});

const instanceList = computed(() => {
  return uniqBy(
    props.databaseList.map((db) => db.instance),
    (instance) => instance.id
  );
});

const filterInstance = (instance: Instance) => {
  return instanceList.value.findIndex((inst) => inst.id === instance.id) >= 0;
};
</script>
