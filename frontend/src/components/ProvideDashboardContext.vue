<template>
  <slot />
</template>

<script lang="ts">
import {
  useEnvironmentStore,
  useMemberStore,
  useRoleStore,
  usePrincipalStore,
  useSettingStore,
  useUIStateStore,
  useProjectStore,
  useDebugStore,
  useUserStore,
} from "@/store";
import { defineComponent } from "vue";
import { DEFAULT_PROJECT_ID } from "../types";

export default defineComponent({
  name: "ProvideDashboardContext",
  async setup() {
    await Promise.all([
      useSettingStore().fetchSetting(),
      // Fetch so MemberSelect can have the data.
      useMemberStore().fetchMemberList(),
      useRoleStore().fetchRoleList(),
      useUserStore().fetchUserList(),
      // Though fetchMemberList also return the principal info, it's possible that a principal is no longer a member.
      // since all record types have creator, updater which are associated with principal, so we need to fetch
      // the principal list as well.
      // We also need this to render the proper inbox and activity entry.
      usePrincipalStore().fetchPrincipalList(),
      useEnvironmentStore().fetchEnvironmentList(),
      // The default project hosts databases not explicitly assigned to other users project.
      useProjectStore().fetchProjectById(DEFAULT_PROJECT_ID),
      useUIStateStore().restoreState(),
      useDebugStore().fetchDebug(),
    ]);
  },
});
</script>
