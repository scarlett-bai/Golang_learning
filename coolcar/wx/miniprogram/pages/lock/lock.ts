import { IAppOption } from "../../appoption"
import { CarService } from "../../service/car"
import { car } from "../../service/proto_gen/car/car_pb"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

const shareLocationKey = "share_location"
// pages/lock/lock.ts
Page({

  carID: '',
  carRefresher: 0,
  data: {
    shareLocation: false,
    avatarURL: '',
    userinfo: {},
    hasUserInfo: false,
    CanIUseGetUserProfile: false,
  },

  async onLoad(opt: Record<'car_id', string>) {
    const o: routing.LockOpts = opt
    this.carID = o.car_id
    // console.log('unlocking car', o.car_id)
    const userInfo = await getApp<IAppOption>().globalData.userInfo
    this.setData({
      avatarURL: userInfo.avatarUrl,
      shareLocation: wx.getStorageSync(shareLocationKey) || false,
    })
  },

  onShareLocation(e: any){
    this.data.shareLocation = e.detail.value
    wx.setStorageSync(shareLocationKey, this.data.shareLocation)
  },
  onGetUserInfo(e: any) {
    const userInfo: WechatMiniprogram.UserInfo = e.detail.userInfo
    if (userInfo) {
      getApp<IAppOption>().resolveUserInfo(userInfo)
      this.setData({
        shareLocation: true,
      })
      wx.setStorageSync(shareLocationKey, true)
    }
  },

  onUnlockTap(){
    wx.getLocation({
      type: 'gcj02',
      success: async loc => {
        if (!this.carID) {
          console.error('no carID specified')
          return
        }
        let trip: rental.v1.ITripEntity
        try {
            trip = await TripService.ceateTrip({
            start: loc,
            carId: this.carID,
            avatarUrl: this.data.shareLocation
              ? this.data.avatarURL: ''
          })
        
          if (!trip.id) {
            console.error('no tripID in response', trip)
            return
          }
        } catch (err) {
          wx.showToast({
            title: '创建行程失败',
            icon: 'none',
          })
          return
        }
        
        // return
        // const tripID = 'trip456'
         wx.showLoading({
            title: '开锁中',
            mask: true,
          })

          // 
          this.carRefresher = setInterval(async () => {
            const c = await CarService.getCar(this.carID)
            if (c.status === car.v1.CarStatus.UNLOCKED) {
                this.clearCarRefresher() 
                wx.redirectTo({
                // url: `/pages/driving/driving?trip_id=${{tripID}}`,
                url: routing.driving({
                  trip_id: trip.id!,
                }),
                complete: () => {
                  wx.hideLoading()
                }
              })
            }
          }, 2000);
         
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title:'请前往设置页面位置授权',
        })
      }
    })
   
  },

  clearCarRefresher() {
    if (this.carRefresher) {
      clearInterval(this.carRefresher)
      this.carRefresher = 0
    }
  },
  /**
   * 页面的初始数据
   */
  

  /**
   * 生命周期函数--监听页面加载
   */
  
  // getUserProfile(){
  //   wx.getUserInfo({
  //     // desc: '用于展示头像',
  //     success:res => {
  //       this.setData({
  //         userinfo: res.userInfo,
  //         hasUserInfo: true,
  //       })
  //     }
  //   })
  // },
  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow() {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide() {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload() {
    this.clearCarRefresher()
    wx.hideLoading()
  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh() {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom() {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage() {

  }
})