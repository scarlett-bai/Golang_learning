import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { formatDuration, formatFee } from "../../utils/format"
import { routing } from "../../utils/routing"

// pages/driving/driving.ts
// const centPerSec = 0.7
const updateIntervalSec = 5  // 每5s向服务器上报数据

// function formatDuration(sec: number) {
//   const padString = (n: number) => 
//     n < 10 ? '0'+n.toFixed(0) : n.toFixed(0)
//   const h = Math.floor(sec/3600)
//   sec -= 3600 * h
//   const m = Math.floor(sec/60)
//   sec -= 60 * m
//   const s = Math.floor(sec)

//   return `${padString(h)}:${padString(m)}:${padString(s)}`
// }

// function formatFee(cents: number) {
//   return (cents / 100).toFixed(2)
// }

function durationStr(sec: number){
  const dur = formatDuration(sec)
  return `${dur.hh}:${dur.mm}:${dur.ss}`
}

Page({
  timer: undefined as number|undefined,
  tripID: '',
  data: {
    location: {
      latitude: 32.92,
      longitude: 118.46,
    },
    scale:14,
    elapsed: '00:00:00',
    fee: '0.00',
  },

  setupLocationUpdator(){
    wx.startLocationUpdate({
      fail: console.error,
    })
    wx.onLocationChange(loc => {
      this.setData({
        location: {
          latitude: loc.latitude,
          longitude: loc.longitude,
        },
      })
    })
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(opt: Record<'trip_id', string>) {
    const o: routing.DrivingOpts = opt
    console.log('current trip', o.trip_id)
    // o.trip_id = '641163c4d98effc263012d70'
    this.tripID = o.trip_id
    TripService.getTrip(o.trip_id).then(console.log)
    this.setupLocationUpdator()
    this.setupTimer(o.trip_id)
  },

  onEndTripTap(){
    TripService.finishTrip(this.tripID).then(() => {
        wx.redirectTo({
        url: routing.mytrips(),
    })
    }).catch(err => {
      console.log(err)
      wx.showToast({
        title: '结束行程失败',
        icon: 'none',
      })
    })
   
  },

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
    wx.stopLocationUpdate()
    if (this.timer) {
      clearInterval(this.timer)
    }
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

  },

   async setupTimer(tripID: string){
    const trip = await TripService.updateTripPos(tripID)
    if (trip.status !== rental.v1.TripStatus.IN_PROGRESS) {
        console.error('Trip not in progress')
        return
    }
    let secSinceLastUpdate = 0   // 自上次更新已经过去了多少秒
    let lastUpdateDurationSec = trip.current!.timestampSec! = trip.start!.timestampSec!
    this.setData({
      elapsed: durationStr(lastUpdateDurationSec),
      fee: formatFee(trip.current!.feeCent!)
    }) 
    this.timer = setInterval(() => {
      secSinceLastUpdate++
      if (updateIntervalSec % 5 === 0) {
        TripService.getTrip(tripID).then(trip => {
          lastUpdateDurationSec = trip.current!.timestampSec! = trip.start!.timestampSec!
          secSinceLastUpdate = 0
          this.setData({
            fee: formatFee(trip.current!.feeCent!)
          })
        }).catch(console.error)
      }
        this.setData({
          elapsed: durationStr(lastUpdateDurationSec + secSinceLastUpdate),
        })
    }, 1000)
  },
})