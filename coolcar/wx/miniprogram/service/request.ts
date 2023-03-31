import camelcaseKeys = require("camelcase-keys")
import { auth } from "./proto_gen/auth/auth_pb"

export namespace Coolcar {
    const serverAddr = 'http://localhost:8080'
    const AUTH_ERR = 'AUTH_ERR'

    const authData = {
        token: '',
        expiryMs: 0,
    }
    export interface RequestOption<REQ, RES> {
        method: 'GET' | 'PUT' | 'POST' | 'DELETE'
        path: string
        data?: REQ 
        respMarshaller: (r: object)=>RES
    }

    export interface AuthOption {
        attachAuthHeader: boolean
        retryOnAuthError: boolean   // 打个标记
    }
    
    export async function sendRquestWithAuthRetry<REQ,RES>(o: RequestOption<REQ, RES>, a?: AuthOption): Promise<RES> {
        const authOpt = a || {  // 如果给了参数就信你。若没有给就初始化一个默认值
            attachAuthHeader: true,
            retryOnAuthError: true,
        }
        try {
            await login()
            return sendRequest(o, authOpt)
        } catch(err) {
            if (err === AUTH_ERR && authOpt.retryOnAuthError) {
                // 清除状态
                authData.token = ''
                authData.expiryMs = 0
                return sendRquestWithAuthRetry(o, {
                    attachAuthHeader: authOpt.attachAuthHeader,
                    retryOnAuthError:false,
                })
            } else {
                throw err
            }
        }
        
    }


    export async function login() {
        if (authData.token && authData.expiryMs >= Date.now()) {  // 若有效 就不用login
            return
        }
        const wxResp = await wxLogin()
        const reqTimeMs = Date.now()
        const resp = await sendRequest<auth.v1.ILoginRequest, auth.v1.ILoginResponse>({
            method: 'POST',
            path: '/v1/auth/login',
            data: {
                code: wxResp.code,

            },
            respMarshaller: auth.v1.LoginResponse.fromObject,
        }, {
            attachAuthHeader: false,
            retryOnAuthError: false,
        })
        authData.token = resp.accessToken!
        authData.expiryMs = reqTimeMs + resp.expireIn! * 1000
    }

    function sendRequest<REQ, RES>(o: RequestOption<REQ, RES>, a: AuthOption): Promise<RES>{
        return new Promise((reslove, reject) => {
            const header: Record<string, any> = {}  // 表示键值，键是string  值是any
            if (a.attachAuthHeader) {
                if (authData.token && authData.expiryMs >= Date.now()) {
                     header.authorization = 'Bearer ' + authData.token
                } else {
                    reject(AUTH_ERR)  // token过期的情况
                    return
                }
            } 
            wx.request({
                url: serverAddr + o.path,
                method: o.method,
                data: o.data as any,
                header,
                success: res => {
                    if (res.statusCode === 401) {  // Unauthorized
                        reject(AUTH_ERR)
                    } else if (res.statusCode >= 400) {  // 200-300多 都算成功 
                        reject(res)
                    } else {
                        reslove(o.respMarshaller(
                            camelcaseKeys(res.data as object, {
                                deep: true,
                            })
                        ))
                    }
                    
                },
                fail: reject,
            })
        })
    }

    function wxLogin(): Promise<WechatMiniprogram.LoginSuccessCallbackResult> {
        return new Promise((resolve, reject) => {
             wx.login({
                success: resolve,
                fail: reject
            })
        }) 
    }
}