<template>
  <div
    v-if="mode !== 'none'"
    class="space-y-1 mt-2 text-xs"
    :class="[
      mode === 'template-mismatch' && 'text-warning',
      mode === 'value-mismatch' && 'text-warning',
    ]"
  >
    <template v-if="mode === 'normal'">
      <p>{{ $t("database.should-follow-database-name-template") }}</p>
      <p>
        <code>{{ project.dbNameTemplate }}</code>
      </p>
    </template>
    <template v-if="mode === 'template-mismatch'">
      <p>{{ $t("database.doesnt-match-database-name-template") }}</p>
      <p>
        <code>{{ project.dbNameTemplate }}</code>
      </p>
    </template>
    <template v-else-if="mode === 'value-mismatch'">
      <p>{{ $t("database.doesnt-match-tenant-value") }}</p>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";

import type { DatabaseLabel, Project } from "@/types";
import {
  buildDatabaseNameRegExpByTemplate,
  getLabelValueFromLabelList,
} from "@/utils";

type ViewMode = "none" | "normal" | "template-mismatch" | "value-mismatch";

const props = defineProps<{
  name: string;
  project: Project;
  labelList: DatabaseLabel[];
}>();

const mode = computed((): ViewMode => {
  const { project, name, labelList } = props;
  if (!project.dbNameTemplate) return "none";
  if (!name) return "normal";
  const regex = buildDatabaseNameRegExpByTemplate(project.dbNameTemplate);
  const matches = name.match(regex);
  if (!matches) return "template-mismatch";
  const parsedTenant = matches.groups?.["TENANT"];
  const tenant = getLabelValueFromLabelList(labelList, "bb.tenant");
  if (parsedTenant && parsedTenant !== tenant) {
    return "value-mismatch";
  }
  return "none";
});
</script>
