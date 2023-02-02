import { IAppOption } from "../../appoption"
import { routing } from "../../utils/routing"

interface Trip {
  id: string
  start: string
  end: string
  duration: string
  fee: string
  distance: string
  status: string
}

interface MainItem {
  id: string
  navId: string
  navScrollId: string
  data: Trip
}

interface NavItem {
  id: string
  mainId: string
  label: string
}

interface MainItemQueryResult {
  id: string
  top: number,
  dataset: {
    navId: string
    navScrollId: string
  }
}
// pages/mytrips/mytrips.ts
Page({
  scrollStates: {
      mainItems: [] as MainItemQueryResult[]
  },

  /**
   * 页面的初始数据
   */
  data: {
      indicatorDots: true,
      autoPlay: false,
      interval: 3000,
      duration: 500,
      circular: true,
      multiItemCount: 1,
      prevMargin:'',
      nextMargin: '',
      vertical: false,
      current:0, 
      promotionItems: [
        {
          img: 'https://wx1.sinaimg.cn/mw2000/005ThgxWly1h82hblqo0yj31w81w7hdt.jpg',
          promotionID: 1,
        },
        {
          img: "https://wx3.sinaimg.cn/mw2000/005ThgxWly1h82hbhvdvoj3236236kjl.jpg",
          promotionID: 2,
        },
        {
          img: "https://wx4.sinaimg.cn/mw2000/df99d00bly1h82j29jmmaj20u00uxgmv.jpg",
          promotionID: 3,
        },
        {
          img: "https://wx2.sinaimg.cn/mw2000/3ff41d7dly1h82kx4itsfj213z0qowig.jpg",
          promotionID:4,
        }

      ],
      avatarURL: '',
      mainItems: [] as MainItem[],
      navItems: [] as NavItem[],
      tripsHeight: 0,
      navCount: 0,
      mainScroll: '',
      navSel:'',
      navScroll: '',
  },

  onSwiperChange(){
    
    // process
  },

  onPromotionItemTap(e: any){
    console.log(e)
    const promotionID = e.currentTarget.dataset.promotionId
    if (promotionID) {

    }
    // promotionID?
  },
  onRegisterTap(){
     wx.navigateTo({
      url: routing.register(),
     })
  },

  /**
   * 生命周期函数--监听页面加载
   */
  async onLoad() {
    this.populateTrips()
    const userInfo = await getApp<IAppOption>().globalData.userInfo
    this.setData({
      avatarURL:userInfo.avatarUrl,
    })
  },


populateTrips(){
    const mainItems: MainItem[] = []
    const navItems: NavItem[] = []
    let navSel =  ''
    let prevNav =  ''
    for (let i = 0; i < 100; i++) {
      const mainId = 'main-' + i
      const navId = 'nav-' + i
      const tripId = (10000+i).toString()
      if (!prevNav) {
        prevNav = navId
      }
      mainItems.push({
        id: mainId,
        navId: navId,
        navScrollId: prevNav,
        data: {
          id: tripId,
          start: '东方明珠',
          end: '迪士尼',
          distance:'27.0公里',
          duration: '0时44分',
         fee: '128.8元',
         status: '已完成',
      },
        
      })
      navItems.push({
        id: navId,
        mainId: mainId,
        label: tripId,
      })
      if (i===0) {
        navSel = navId
      } 
      prevNav = navId
    }
  
    this.setData({
      mainItems,
      navItems,
      navSel,
    }, () => {
        this.prepareScrollStates()
    })
    // console.log(trips)
  },

  prepareScrollStates(){
    wx.createSelectorQuery().selectAll('.main-item').fields({
      id: true,
      dataset: true,
      rect: true,
    }).exec(res => {
        this.scrollStates.mainItems = res[0]
    })
  },

  onNavItemTap(e:any){
      const mainId: string = e.currentTarget?.dataset?.mainId
      const navId: string = e.currentTarget?.id
      if (mainId && navId) {
        this.setData({
          mainScroll: mainId,
          navSel: navId,
        })
      }
  },

  onMainScroll(e: any){
    const top: number = e.currentTarget?.offsetTop + e.detail?.scrollTop
    if (top===undefined) {
      return
    }

    const selItem = this.scrollStates.mainItems.find(v => v.top >= top)
    if (!selItem) {
      return
    }

    this.setData({
      navSel: selItem.dataset.navId,
      navScroll: selItem.dataset.navScrollId
    })
  },
  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady() {
    wx.createSelectorQuery().select('#heading').boundingClientRect(rect => {
      const height = wx.getSystemInfoSync().windowHeight - rect.height
      this.setData({
          tripsHeight: height,
          navCount : Math.round(height/50)
      })
      
    }).exec()

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