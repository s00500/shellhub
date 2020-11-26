<template>
  <fragment>
    <v-menu
      :close-on-content-click="false"
      offset-y
    >
      <template #activator="{ on }">
        <v-chip v-on="on">
          <v-icon left>
            mdi-server
          </v-icon>
          <span>
            {{ namespace.name }}
          </span>
          <v-icon right>
            mdi-chevron-down
          </v-icon>
        </v-chip>
      </template>

      <v-card>
        <v-subheader>Tenant ID</v-subheader>

        <v-list
          class="pt-0 pb-0"
        >
          <v-list-item>
            <v-list-item-content>
              <v-chip>
                <v-list-item-title>
                  <span data-test="tenantID-text">{{ tenant }}</span>
                </v-list-item-title>
                <v-icon
                  v-clipboard="tenant"
                  v-clipboard:success="() => {
                    this.$store.dispatch('snackbar/showSnackbarCopy', this.$copy.tenantId);
                  }"
                  right
                >
                  mdi-content-copy
                </v-icon>
              </v-chip>
            </v-list-item-content>
          </v-list-item>
        </v-list>

        <v-divider />

        <v-list
          class="pt-0"
        >
          <v-subheader>Namespaces</v-subheader>

          <v-list-item-group>
            <v-virtual-scroll
              :max-height="149"
              item-height="50"
              :items="namespaces"
            >
              <template #default="{ item }">
                <v-list-item
                  :key="item.tenant_id"
                  @click="switchIn(item.tenant_id)"
                >
                  <v-list-item-icon>
                    <v-icon>mdi-login</v-icon>
                  </v-list-item-icon>
                  <v-list-item-content>
                    <v-list-item-title>
                      {{ item.name }}
                    </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-virtual-scroll>
          </v-list-item-group>
        </v-list>

        <v-divider />

        <v-list
          class="pt-0 pb-0"
        >
          <v-list-item
            v-show="show"
            @click="dialog=!dialog"
          >
            <v-list-item-icon>
              <v-icon>mdi-plus-box</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              Create Namespace
            </v-list-item-content>
          </v-list-item>

          <v-divider />

          <v-list-item
            to="/settings/namespace-manager"
          >
            <v-list-item-icon>
              <v-icon>mdi-cog</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              Settings
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-card>

      <NamespaceAdd
        :show.sync="dialog"
      />
    </v-menu>
  </fragment>
</template>

<script>
import NamespaceAdd from '@/components/namespace/NamespaceAdd';

export default {
  name: 'NamespaceMenu',

  components: {
    NamespaceAdd,
  },

  props: {
    inANamespace: {
      type: Boolean,
      required: true,
    },
  },

  data() {
    return {
      model: true,
      dialog: false,
    };
  },

  computed: {
    isOwner() {
      return this.owner === this.$store.getters['auth/id'];
    },

    owner() {
      return this.$store.getters['namespaces/get'].owner;
    },

    namespace() {
      return this.$store.getters['namespaces/get'];
    },

    namespaces() {
      return this.$store.getters['namespaces/list'];
    },

    tenant() {
      return localStorage.getItem('tenant');
    },

    show() {
      return this.$env.isHosted;
    },
  },

  watch: {
    dialog(value) {
      if (!value) {
        this.model = false;
      }
    },
  },

  created() {
    if (this.$props.inANamespace) {
      this.getNamespaces();
    }
  },

  methods: {
    async getNamespaces() {
      try {
        // load namespaces
        await this.$store.dispatch('namespaces/get', this.tenant);
      } catch (e) {
        if (e.response.status === 403) {
          this.$store.dispatch('snackbar/showSnackbarErrorAssociation');
        } else {
          this.$store.dispatch('snackbar/showSnackbarErrorLoading', this.$errors.namespaceList);
        }
      }
    },

    async switchIn(tenant) {
      try {
        await this.$store.dispatch('namespaces/switchNamespace', {
          tenant_id: tenant,
        });
        window.location.reload();
      } catch {
        this.$store.dispatch('snackbar/showSnackbarErrorLoading', this.$errors.namespaceSwitch);
      }
    },
  },
};
</script>
