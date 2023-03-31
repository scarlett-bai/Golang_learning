import { IAppOption } from "./appoption";
import { Coolcar } from "./service/request";
// import { coolcar } from "./service/proto_gen/trip_pb";
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
    // wx.request({
    //   url: 'http://localhost:8080/trip/trip123',
    //   method:'GET',
    //   success: res => {
    //     const GetTripRes = coolcar.GetTripResponse.fromObject(camelcaseKeys(res.data as object, {
    //       deep: true,
    //     }))
    //     console.log("getTripRes:", GetTripRes)
    //     console.log('status is', coolcar.TripStatus[GetTripRes.trip?.status!])
    //   },
    //   fail: console.error,
    // })
    // 登录
    Coolcar.login()
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