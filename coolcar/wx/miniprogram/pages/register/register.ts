import { UploadFileOpts } from './../../service/request';
import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { Coolcar } from "../../service/request"
import { padString } from "../../utils/format"
import { routing } from "../../utils/routing"


function formatDate(millis: number) {
  const dt = new Date(millis)
  const y = dt.getFullYear()
  const m = dt.getMonth() + 1
  const d = dt.getDate()
  return `${padString(y)}-${padString(m)}-${padString(d)}`
}

// pages/register/register.ts
Page({
  redirectURL: '',
  profileRefresher: 0,
  onUploadLic(){
    wx.chooseImage({
      success: async res => {
        if (res.tempFilePaths.length === 0){
            return
        }
        this.setData({
            licImgURL: res.tempFilePaths[0]
          })
          //TODO: upload image
          const photoRes = await ProfileService.createProfilePhoto()
          if (!photoRes.uploadUrl) {
            return
          }
          await Coolcar.uploadfile({
            localPath: res.tempFilePaths[0],
            url: photoRes.uploadUrl,
          })

          const identity = await ProfileService.completeProfilePhoto()
          this.renderIdentity(identity)
        }
      })
    
  },
  onGenderChange(e: any){
    this.setData({
      genderIndex: parseInt(e.detail.value),
    })
  },

  onBirthDateChange(e: any){
    this.setData({
      birthDate: e.detail.value,
    })
  },
  onSubmit(){
   ProfileService.submitProfile({
     licNumber: this.data.licNo,
     name: this.data.name,
     gender: this.data.genderIndex,
     birthDateMillis: Date.parse(this.data.birthDate),
   }).then(p => {
      this.renderProfile(p)
      this.scheduleProfileRefresher()
    })
  },

  onUnload(){
    this.clearProfileRefresher()
  },

  scheduleProfileRefresher() {
      this.profileRefresher = setInterval(() => {
        console.log('aaaaaaaa')
        ProfileService.getProfile().then(p => {
          this.renderProfile(p)
          if (p.identityStatus !== rental.v1.IdentityStatus.PENDING) {
            this.clearProfileRefresher() 
          }
          if (p.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
            this.onLicVerified()
          }
        })
      }, 1000)
    },

  clearProfileRefresher(){
    if (this.profileRefresher) {
      clearInterval(this.profileRefresher)
      this.profileRefresher = 0
    }
  },

  onResubmit(){
    ProfileService.clearProfile().then(p => this.renderProfile(p))
    ProfileService.clearProfilePhoto().then(() => {
      this.setData({
        licImgURL: '',
      })
    })
  },
  onLicVerified(){
    if (this.redirectURL) {
      wx.redirectTo({
      url:this.redirectURL,
    })
    }
    
  },
  /**
   * 页面的初始数据
   */
  data: {
    licNo: '',
    name: '',
    genderIndex: 0,
    genders: ['未知', '男', '女'],
    birthDate: '1990-01-01',
    licImgURL: '',
    state:rental.v1.IdentityStatus[rental.v1.IdentityStatus.UNSUBMITTED],
  },

  /**
   * 生命周期函数--监听页面加载
   */



  renderProfile(p: rental.v1.IProfile) {
    this.renderIdentity(p.identity!)
    this.setData({
      state: rental.v1.IdentityStatus[p.identityStatus||0]
    })
  },

  renderIdentity(i?: rental.v1.IIdentity) {
    this.setData({
      licNo: i?.licNumber || '',
      name: i?.name||'',
      genderIndex: i?.gender||0,
      birthDate: formatDate(i?.birthDateMillis||0)
    })
  },
  
  onLoad(opt: Record<'redirect', string>) {
      const o: routing.RegisterOpts = opt
      if (o.redirect) {
      this.redirectURL = decodeURIComponent(o.redirect)
   }
      ProfileService.getProfile().then(p => this.renderProfile(p))
      ProfileService.getProfilePhoto().then(p => {
        this.setData({
          licImgURL: p.url || '',
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