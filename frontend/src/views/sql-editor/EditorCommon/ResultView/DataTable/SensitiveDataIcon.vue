<template>
  <NTooltip :animated="false" :delay="250">
    <template #trigger>
      <heroicons-outline:eye-slash
        class="w-[12px] h-[12px] -mb-[2px]"
        :class="[clickable && 'cursor-pointer']"
        v-bind="$attrs"
        @click="handleClick"
      />
    </template>

    <span class="whitespace-nowrap">
      {{ $t("sensitive-data.self") }}
    </span>
  </NTooltip>
</template>

<script lang="ts" setup>
import { NTooltip } from "naive-ui";
import { computed } from "vue";
import { useRouter } from "vue-router";

import { useCurrentUser } from "@/store";
import { hasWorkspacePermission } from "@/utils";

const user = useCurrentUser();
const router = useRouter();

const clickable = computed(() => {
  return hasWorkspacePermission(
    "bb.permission.workspace.manage-sensitive-data",
    user.value.role
  );
});

const handleClick = () => {
  if (!clickable.value) {
    return;
  }
  const url = router.resolve({ name: "setting.workspace.sensitive-data" });
  window.open(url.href, "_BLANK");
};
</script>
