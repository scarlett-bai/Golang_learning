// index.ts

import { CarService } from "../../service/car"
import { ProfileService } from "../../service/profile"
import { car } from "../../service/proto_gen/car/car_pb"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

interface Marker {
  iconPath: string
  id: number
  latitude: number
  longitude: number
  width: number
  height: number
}

// 获取应用实例
const defaultAvatarUrl = "/resources/account.png"
const defaultAvatar = "/resources/car.png"
const initialLat = 30
const initialLng = 120


Page({
  isPageShowing: false,
  socket: undefined as WechatMiniprogram.SocketTask | undefined,
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
    latitude: initialLat,
    longitude: initialLng,
  },
  scale: 10,
  markers: [] as Marker[],
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



async onLoad(){
  

  
  // const userInfo = await getApp<IAppOption>().globalData.userInfo
  // this.setData({
  //   avatarURL: userInfo.avatarUrl,
  // })
},
  async onScanTap(){
    const trips = await TripService.getTrips(rental.v1.TripStatus.IN_PROGRESS)
    if ((trips.trips?.length || 0) > 0) {
      await this.selectComponent('#tripModel').showModel()
      wx.navigateTo({
        url: routing.driving({
          trip_id: trips.trips![0].id!,
        })
      })
      return
    }
    wx.scanCode({
      success: async () => {
        const carID = '643cf69c4c1759ffea405a8d'
        // const carID = 'car123'
        const lockURL = routing.lock({
          car_id: carID,
        })
        const prof = await ProfileService.getProfile()
        if (prof.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
          wx.navigateTo({
            url: lockURL,
          }) 
        } else {
           await this.selectComponent('#licModel').showModel()
            wx.navigateTo({
              // url: `/pages/register/register?redirect=${encodeURIComponent(redirectURL)}`,
              url: routing.register({
                redirectURL: lockURL
              })
            })
        }
        
       
      },
      fail: console.error
    })
  },
  onShow(){
    this.isPageShowing = true;
    if (!this.socket) {
      this.setData({
        markers: []
      }, () => {
        this.setupCarPosUpdater()
      })
    }
  },
  onHide(){
    this.isPageShowing = false;
    if (this.socket) {
      this.socket.close({
        success: () => {
          this.socket = undefined
        }
      })
    }
  },

  onModelOK(){
    console.log('ok clicked')
  },

  setupCarPosUpdater() {
    // move markers
    const map = wx.createMapContext("map")
    const markersBNyCarID = new Map<string, Marker>()
    const endTranslation = () => {
          translationInProgress = false
        }
    let translationInProgress = false
      this.socket = CarService.subscrbe(car => {
        if (!car.id || translationInProgress || !this.isPageShowing) {
          console.log('dropped')
          return
        }
       const marker = markersBNyCarID.get(car.id)
       if (!marker) {
        // Insert new marker
        const newMarker: Marker = {
          id: this.data.markers.length,
          iconPath: car.car?.driver?.avatarUrl || defaultAvatar,
          latitude: car.car?.position?.latitude || initialLat,
          longitude: car.car?.position?.longitude || initialLng,
          height: 20,
          width: 20,
        }
        markersBNyCarID.set(car.id, newMarker)
        this.data.markers.push(newMarker)
        translationInProgress = true
        this.setData({
          markers: this.data.markers,
        }, endTranslation)
        return
       }
       const nweAvatar = car.car?.driver?.avatarUrl || defaultAvatar
       const newLat = car.car?.position?.latitude || initialLat
       const newLng = car.car?.position?.longitude || initialLng
       if (marker.iconPath !== nweAvatar) {
          //  Change iconPath and possibly position
          marker.iconPath = nweAvatar
          marker.latitude =  newLat
          marker.longitude = newLng 
          translationInProgress = true
          this.setData({
            markers: this.data.markers,
          }, endTranslation)
          return
       }
       if (marker.latitude != newLat || marker.longitude !== newLng) {
        translationInProgress = true
         map.translateMarker({
          markerId: marker.id,
          destination: {
            latitude: newLat,
            longitude: newLng,
          },
          autoRotate: false,
          rotate: 0,
          duration: 900,
          animationEnd: endTranslation,
         })
       }
      })
  },


})
