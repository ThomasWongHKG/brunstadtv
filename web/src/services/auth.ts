import { useAuth0 } from "@/services/auth0"

// type Token = {
//     expiresAt: string
//     token: string
// }

export class Auth {
    public static async signIn() {
        const { loginWithRedirect } = useAuth0()
        return await loginWithRedirect()
    }

    public static async signOut() {
        const { logout } = useAuth0()
        return await logout({
            returnTo: window.location.origin,
        })
    }

    public static loading() {
        const { isLoading } = useAuth0()
        return isLoading;
    }

    // private static _token: Token | null = JSON.parse(
    //     localStorage.getItem("token") ?? "null"
    // )

    public static async getToken() {
        const { getAccessTokenSilently } = useAuth0()
        return getAccessTokenSilently();
        // if (!this.isAuthenticated().value) {
        //     return null
        // }
        // if (this._token) {
        //     const date = new Date(this._token.expiresAt)
        //     if (date.getTime() > new Date().getTime()) return this._token.token
        // }
        // try {
        //     const date = new Date()
        //     date.setHours(date.getHours() + 12)
        //     console.log("what")
        //     this._token = {
        //         expiresAt: date.toISOString(),
        //         token: await getAccessTokenSilently(),
        //     }
        //     localStorage.setItem("token", JSON.stringify(this._token))
        //     return this._token.token
        // } finally {
        // }
    }

    public static isAuthenticated() {
        const { isAuthenticated } = useAuth0()
        return isAuthenticated
    }
}

export default Auth
