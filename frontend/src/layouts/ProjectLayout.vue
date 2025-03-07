<template>
  <div class="py-2">
    <ArchiveBanner v-if="project.rowStatus == 'ARCHIVED'" />
  </div>
  <h1 class="px-6 pb-4 text-xl font-bold leading-6 text-main truncate">
    <template v-if="project.id === DEFAULT_PROJECT_ID">
      {{ $t("database.unassigned-databases") }}
    </template>
    <template v-else>
      {{ project.name }}
    </template>
    <span
      v-if="project.tenantMode === 'TENANT'"
      class="text-sm font-normal px-2 ml-2 rounded whitespace-nowrap inline-flex items-center bg-gray-200"
    >
      {{ $t("project.mode.tenant") }}
    </span>
  </h1>
  <BBAttention
    v-if="project.id === DEFAULT_PROJECT_ID"
    class="mx-6 mb-4"
    :style="'INFO'"
    :title="$t('project.overview.info-slot-content')"
  />
  <BBTabFilter
    class="px-3 pb-2 border-b border-block-border"
    :responsive="false"
    :tab-item-list="tabItemList"
    :selected-index="state.selectedIndex"
    @select-index="selectTab"
  />

  <div class="py-6 px-6">
    <router-view
      :project-slug="projectSlug"
      :project-webhook-slug="projectWebhookSlug"
      :allow-edit="allowEdit"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, reactive, watch } from "vue";
import { useRouter } from "vue-router";
import { startCase } from "lodash-es";

import {
  idFromSlug,
  hasWorkspacePermission,
  hasPermissionInProject,
} from "../utils";
import ArchiveBanner from "../components/ArchiveBanner.vue";
import { BBTabFilterItem } from "../bbkit/types";
import { useI18n } from "vue-i18n";
import { Project, DEFAULT_PROJECT_ID } from "../types";
import { useCurrentUser, useProjectStore } from "@/store";

type ProjectTabItem = {
  name: string;
  hash: string;
};

interface LocalState {
  selectedIndex: number;
}

export default defineComponent({
  name: "ProjectLayout",
  components: {
    ArchiveBanner,
  },
  props: {
    projectSlug: {
      required: true,
      type: String,
    },
    projectWebhookSlug: {
      type: String,
      default: undefined,
    },
  },
  setup(props) {
    const router = useRouter();
    const { t } = useI18n();

    const currentUser = useCurrentUser();
    const projectStore = useProjectStore();

    const project = computed((): Project => {
      return projectStore.getProjectById(idFromSlug(props.projectSlug));
    });

    const isDefaultProject = computed((): boolean => {
      return project.value.id === DEFAULT_PROJECT_ID;
    });

    const isTenantProject = computed((): boolean => {
      return project.value.tenantMode === "TENANT";
    });

    const projectTabItemList = computed((): ProjectTabItem[] => {
      const list: (ProjectTabItem | null)[] = [
        { name: t("common.overview"), hash: "overview" },
        { name: t("common.databases"), hash: "databases" },

        isTenantProject.value
          ? null // Hide "Change History" tab for tenant projects
          : { name: t("common.change-history"), hash: "change-history" },

        { name: startCase(t("slow-query.slow-queries")), hash: "slow-query" },

        { name: t("common.activities"), hash: "activity" },
        { name: t("common.gitops"), hash: "gitops" },
        { name: t("common.webhooks"), hash: "webhook" },
        isDefaultProject.value
          ? null
          : { name: t("common.settings"), hash: "setting" },
      ];
      const filteredList = list.filter(
        (item) => item !== null
      ) as ProjectTabItem[];

      return filteredList;
    });

    const findTabIndexByHash = (hash: string) => {
      hash = hash.replace(/^#?/g, "");
      const index = projectTabItemList.value.findIndex(
        (item) => item.hash === hash
      );
      if (index >= 0) {
        return index;
      }
      // otherwise fallback to the first tab
      return 0;
    };

    const state = reactive<LocalState>({
      selectedIndex: findTabIndexByHash(router.currentRoute.value.hash),
    });

    const allowEdit = computed(() => {
      if (project.value.rowStatus == "ARCHIVED") {
        return false;
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
          project.value,
          currentUser.value,
          "bb.permission.project.manage-general"
        )
      ) {
        return true;
      }
      return false;
    });

    const tabItemList = computed((): BBTabFilterItem[] => {
      return projectTabItemList.value.map((item) => {
        return {
          title: item.name,
          alert: false,
        };
      });
    });

    const selectProjectTabOnHash = () => {
      const { name, hash } = router.currentRoute.value;
      if (name == "workspace.project.detail") {
        const index = findTabIndexByHash(hash);
        selectTab(index);
      } else if (
        name == "workspace.project.hook.create" ||
        name == "workspace.project.hook.detail"
      ) {
        state.selectedIndex = findTabIndexByHash("webhook");
      }
    };

    const selectTab = (index: number) => {
      state.selectedIndex = index;
      router.replace({
        name: "workspace.project.detail",
        hash: "#" + projectTabItemList.value[index].hash,
      });
    };

    watch(
      () => router.currentRoute.value.hash,
      () => {
        selectProjectTabOnHash();
      },
      {
        immediate: true,
      }
    );

    return {
      DEFAULT_PROJECT_ID,
      state,
      project,
      allowEdit,
      tabItemList,
      selectTab,
    };
  },
});
</script>
