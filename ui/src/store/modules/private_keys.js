import Vue from 'vue';

export default {
  namespaced: true,

  state: {
    privateKeys: JSON.parse(localStorage.getItem('privateKeys')) || [],
    numberPrivateKeys: (JSON.parse(localStorage.getItem('privateKeys')) || []).length,
  },

  getters: {
    list: (state) => state.privateKeys,
    getNumberPrivateKeys: (state) => state.numberPrivateKeys,
  },

  mutations: {
    setPrivateKey: (state, privateKey) => {
      let { numberPrivateKeys } = state;

      state.privateKeys.push(privateKey);
      Vue.set(state, 'numberPrivateKeys', numberPrivateKeys += 1);
    },

    editPrivateKey: (state, privateKey) => {
      Vue.set(state, 'privateKeys', state.privateKeys.map((i) => (i.data === privateKey.data ? { ...i, name: privateKey.name } : i)));
    },

    removePrivateKey: (state, data) => {
      state.privateKeys.splice(state.privateKeys.findIndex((d) => d.data === data), 1);
    },

    clearListPrivateKeys: (state) => {
      Vue.set(state, 'privateKeys', []);
      Vue.set(state, 'numberPrivateKeys', 0);
    },
  },

  actions: {
    set: async (context, privateKey) => {
      const privateKeys = JSON.parse(localStorage.getItem('privateKeys')) || [];
      privateKeys.push(privateKey);

      localStorage.setItem('privateKeys', JSON.stringify(privateKeys));
      context.commit('setPrivateKey', privateKey);
    },

    edit: async (context, privateKey) => {
      context.commit('editPrivateKey', privateKey);
    },

    remove: async (context, data) => {
      const privateKeys = JSON.parse(localStorage.getItem('privateKeys')) || [];

      if (privateKeys !== null) {
        privateKeys.splice(privateKeys.findIndex((d) => d.data === data), 1);
      }

      localStorage.setItem('privateKeys', JSON.stringify(privateKeys));
      context.commit('removePrivateKey', data);
    },
  },
};
