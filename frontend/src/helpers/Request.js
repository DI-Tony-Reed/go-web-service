import Store from "../store"

export default class Request {
    constructor(url) {
        this.base = 'http://localhost:8081'
        this.url = url
    }

    // Return 'thenable' promise
    async call(parameters = {}, method = 'GET') {
        Store.clearErrors()

        let options = {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: method,
        }

        if (parameters && method !== "GET") {
            options.body = JSON.stringify(parameters)
        }

        Store.setWaitingOnAjax(true)

        return await fetch(`${ this.base }/${ this.url }`, options)
            .then(response => {
                Store.setWaitingOnAjax(false)

                return response.json()
            }).then(response => {
                if (response.errors) {
                    Store.setErrors(response.errors)
                }

                return response
            })
    }

    async post(parameters) {
        return await this.call(parameters, 'POST')
    }

    async get(parameters) {
        return await this.call(parameters, 'GET')
    }
}