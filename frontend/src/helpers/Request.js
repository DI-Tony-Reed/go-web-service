import {store} from "../store"

export const GET = 'GET'
export const POST = 'POST'
export const PUT = 'PUT'
export const PATCH = 'PATCH'
export const DELETE = 'DELETE'

export default class Request {

  constructor(url) {
    this.base = 'http://localhost:8081'
    this.url = url
  }

  // Return 'thenable' promise
  async call(parameters = {}, method = GET) {
    store.clearErrors()
    store.clearMessages()

    let options = {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      method: method,
    }

    if (parameters && method !== GET) {
      options.body = JSON.stringify(parameters)
    }

    store.setWaitingOnAjax(true)

    return await fetch(`${this.base}/${this.url}`, options)
      .then(response => {
        store.setWaitingOnAjax(false)

        return response.json()
      }).then(response => {
        if (response.errors) {
          store.setErrors(response.errors)
        }

        if (response.message) {
          store.setMessages(response.message)
        }

        return response
      })
  }

  async post(parameters) {
    return await this.call(parameters, POST)
  }

  async patch(parameters) {
    return await this.call(parameters, PATCH)
  }

  async put(parameters) {
    return await this.call(parameters, PUT)
  }

  async delete(parameters) {
    return await this.call(parameters, DELETE)
  }

  async get(parameters) {
    return await this.call(parameters, GET)
  }
}