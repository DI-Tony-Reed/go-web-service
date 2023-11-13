import {reactive} from 'vue'

export const store = reactive({
  albums: [],

  messages: [],
  errors: [],
  waitingOnAjax: false,

  setWaitingOnAjax(boolean) {
    this.waitingOnAjax = boolean
  },

  setMessages: function (messages) {
    this.messages = [messages]
  },

  clearMessages() {
    this.messages = []
  },

  setErrors: function (errors) {
    if (typeof errors === 'string') {
      errors = [errors]
    }

    this.errors = errors
  },

  clearErrors() {
    this.errors = []
  },
})