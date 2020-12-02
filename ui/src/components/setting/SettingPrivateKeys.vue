<template>
  <v-form>
    <v-container>
      <v-card class="mt-2 pa-4 pb-0 elevation-1">
        <div class="d-flex pa-0 align-center">
          <h3>Private Keys</h3>
          <v-spacer />
          <v-spacer />
          <PrivateKeyFormDialog
            :create-key="true"
            action="Private"
          />
        </div>

        <v-data-table
          :headers="headers"
          :items="getListPrivateKeys"
          data-test="dataTable-field"
          :server-items-length="getNumberPrivateKeys"
          hide-default-footer
        >
          <template #[`item.name`]="{ item }">
            {{ item.name }}
          </template>

          <template #[`item.privatekey`]="{ item }">
            {{ item.privatekey }}
          </template>

          <template #[`item.actions`]="{ item }">
            <PrivateKeyFormDialog
              :key-object="item"
              :create-key="false"
              action="Private"
            />
            <PrivateKeyDelete
              :fingerprint="item.data"
              action="Private"
            />
          </template>
        </v-data-table>
      </v-card>
    </v-container>
  </v-form>
</template>

<script>

import PrivateKeyFormDialog from '@/components/public_key/KeyFormDialog';
import PrivateKeyDelete from '@/components/public_key/KeyDelete';

export default {
  name: 'SettingPrivateKeys',

  components: {
    PrivateKeyFormDialog,
    PrivateKeyDelete,
  },

  data() {
    return {
      pagination: {},
      headers: [
        {
          text: 'Name',
          value: 'name',
          align: 'center',
        },
        {
          text: 'Private Key',
          value: 'data',
          align: 'center',
        },
        {
          text: 'Actions',
          value: 'actions',
          align: 'center',
        },
      ],
    };
  },

  computed: {
    getListPrivateKeys() {
      return this.$store.getters['privatekeys/list'];
    },

    getNumberPrivateKeys() {
      return this.$store.getters['privatekeys/getNumberPrivateKeys'];
    },
  },
};
</script>
