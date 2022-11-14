import { getSetting, getUserInfo } from "./utils/wxapi"

let resolveUserInfo: (value: WechatMiniprogram.UserInfo | PromiseLike<WechatMiniprogram.UserInfo>) => void;
let rejectUserInfo: (reason?: any) => void
// app.ts
App<IAppOption>({
  globalData: {
    userInfo: new Promise((resolve, reject) => {
      resolveUserInfo = resolve
      rejectUserInfo = reject
    })
  },
  async onLaunch() {

    // 登录
    wx.login({
      success(res) {
        wx.request({
          url: 'https://example.com/onLogin',
          data: {
            code: res.code
          }
        })
      console.log(res)
      }
    })
    try  {
        const setting = await getSetting()
        if (setting.authSetting['scope.userInfo']){
          const userInfoRes = await getUserInfo()
          resolveUserInfo(userInfoRes.userInfo)
        } 
      } catch (err){
      rejectUserInfo(err)
    }
  

},
  resolveUserInfo(userInfo: WechatMiniprogram.UserInfo){
    resolveUserInfo(userInfo)
  }
})