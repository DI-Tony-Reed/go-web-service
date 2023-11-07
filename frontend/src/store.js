export default {
  albums: [],

  errors: [],
  waitingOnAjax: false,

  setWaitingOnAjax(boolean) {
    this.waitingOnAjax = boolean
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
}