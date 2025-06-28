import { defineStore } from 'pinia'
import type { Client } from 'src/services/api-service'

interface State {
  clients: Client[]
  selectedClientId: number | null
}

export const useStore = defineStore('main', {
  state: (): State => ({
    clients: [],
    selectedClientId: null,
  }),

  actions: {
    setClients(clients: Client[]) {
      this.clients = clients
      // Set the first client as selected if available
      if (clients.length > 0) {
        this.selectedClientId = this.clients[0]?.id ?? null
      }
    },

    setSelectedClient(clientId: number) {
      this.selectedClientId = clientId
    },
  },
})
