declare module 'vue' {
    interface ComponentCustomProperties {
        $http: AxiosStatic
        $translate: (key: string) => string
    }
}