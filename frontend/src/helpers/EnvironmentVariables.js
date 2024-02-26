export default class EnvironmentVariables {
  toArray() {
    return import.meta.env
  }

  applicationPort() {
    return this.toArray()['VITE_APPLICATION_PORT']
  }

  applicationUrl() {
    return this.toArray()['VITE_APPLICATION_URL']
  }

  applicationProtocol() {
    return this.toArray()['VITE_APPLICATION_PROTOCOL']
  }
}