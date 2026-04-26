import { defineStore } from 'pinia'

interface State {
    uesr: any
}

interface UserInfo {
    [key: string]: any
}

const useUserStore = defineStore('user', {
    state: (): State => ({
        uesr: {}
    }),
    getters: {
        getUserInfo(): UserInfo {
            return this.uesr
        }
    },
    actions: {
        SET_USERINFO(user: UserInfo) {
            this.uesr = user
        }
    }
})

export default useUserStore