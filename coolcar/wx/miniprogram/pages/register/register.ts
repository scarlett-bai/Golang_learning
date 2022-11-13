import { routing } from "../../utils/routing"

// pages/register/register.ts
Page({
  redirectURL: '',
  onUploadLic(){
    wx.chooseMedia({
      success: res => {
        if (res.tempFiles[0].tempFilePath.length > 0){
            this.setData({
            licImgURL: res.tempFiles[0].tempFilePath
          })
          //TODO: upload image
          setTimeout(()=>{
              this.setData({
                licNo:'1234566789',
                name:'落意',
                genderIndex:2,
                birthDate: '1992-06-27',
              })
          }, 1000)
        }
      }
    })
  },
  onGenderChange(e: any){
    this.setData({
      genderIndex: e.detail.value,
    })
  },

  onBirthDateChange(e: any){
    this.setData({
      birthDate: e.detail.value,
    })
  },
  onSubmit(){
    // TODO: submit the form to server
    this.setData({
      state: 'PENDING',
    })
    setTimeout(() => {
       this.onLicVerified()
    }, 3000);
  },

  onResubmit(){
    this.setData({
      state: 'UNSUBMITTED',
      licImgURL: '', 
    })
  },
  onLicVerified(){
    this.setData({
      state:'VERIFIED'
    })
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
    genders: ['未知', '男', '女', '其他'],
    birthDate: '1990-01-01',
    licImgURL: '',
    state:'UNSUBMITTED' as 'UNSUBMITTED' | 'PENDING' | 'VERIFIED',
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(opt: Record<'redirect', string>) {
      const o: routing.RegisterOpts = opt
      if (o.redirect) {
      this.redirectURL = decodeURIComponent(o.redirect)
   }
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