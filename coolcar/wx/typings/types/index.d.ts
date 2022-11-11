/// <reference path="./wx/index.d.ts" />

interface IAppOption {
    globalData: {
        userInfo: Promise<WechatMiniprogram.UserInfo>,
    }
}