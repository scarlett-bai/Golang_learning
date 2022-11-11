// index.ts
// 获取应用实例
const app = getApp<IAppOption>()

Page({
  data: {
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
  }
})
