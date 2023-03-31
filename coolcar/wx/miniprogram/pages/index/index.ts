// index.ts

import { routing } from "../../utils/routing"

// 获取应用实例
const defaultAvatarUrl = "/resources/account.png"
Page({
  isPageShowing: false,
  data: {
    avatarUrl: defaultAvatarUrl,
    setting: {
    skew: 0,
    rotate:0,
    showLocation: true,
    showScale:true,
    subKey: '',
    layerStyle:-1,
    enableZoom: true,
    enableScroll:true,
    enableRotate:false,
    showCompass:false,
    enable3D:false,
    enableOverlooking:false,
    enableSatellite:false,
    enableTraffic:false,
  },
  location: {
    latitude: 31,
    longitude: 120,
  },
  scale: 10,
  markers: [
    {
      iconPath : "/resources/car.png",
      id: 0,
      latitude: 23.099994,
      longitude: 113.234520,
      width:50,
      height:50
    },
    {
      iconPath : "/resources/car.png",
      id: 1,
      latitude: 23.099994,
      longitude: 113.234520,
      width:50,
      height:50
    }
  ],
},
onChooseAvatar(e:any) {
  console.log("avatarUrl",e.detail.avatarUrl)
  // const { avatarUrl } = e.detail.avatarUrl
  // wx.uploadFile({
  //   url: 
  // })
  this.setData({
    avatarUrl: e.detail.avatarUrl,
  })

  

},
onMyTripsTap(){
  wx.navigateTo({
    url: routing.mytrips(),
  })
},
onMyLocationTap(){
  wx.getLocation({
    type: 'gcj02',
    success: res => {
      this.setData({
        location: {
          latitude: res.latitude,
          longitude: res.longitude,
        },
      })
    },
    fail: () =>{
      wx.showToast({
        icon: 'none',
        title: '请前往设置页授权'
      })
    }
  })
},



// async onLoad(){
//   const userInfo = await getApp<IAppOption>().globalData.userInfo
//   this.setData({
//     avatarURL: userInfo.avatarUrl,
//   })
// },
  onScanTap(){
    wx.scanCode({
      success: async () => {
        await this.selectComponent('#licModel').showModel()
        const carID = 'car123'
        const redirectURL = routing.lock({
          car_id: carID,
        })
        wx.navigateTo({
                // url: `/pages/register/register?redirect=${encodeURIComponent(redirectURL)}`,
                url: routing.register({
                  redirectURL: redirectURL
                })
              })
       
      },
      fail: console.error
    })
  },
  onShow(){
    this.isPageShowing = true;
  },
  onHide(){
    this.isPageShowing = false;
  },

  onModelOK(){
    console.log('ok clicked')
  },
  moveCars(){
    const map = wx.createMapContext("map")
    const dest = {
      latitude: 23.099994,
      longitude: 113.234520,
    }

    const moveCar = () => {
      dest.latitude += 0.1
      dest.longitude += 0.1
      map.translateMarker({
      destination: {
         latitude: dest.latitude,
         longitude: dest.longitude,
      },
      markerId: 0,
      autoRotate: false,
      rotate: 0,
      duration: 5000,
      animationEnd: () => {
        if (this.isPageShowing) {
          moveCar()
        }
      },
    })
    }
    moveCar()
    
  }


})
