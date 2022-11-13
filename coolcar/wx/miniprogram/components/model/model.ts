import { ModelResult } from './../types';
// components/model/model.ts
Component({
  /**
   * 组件的属性列表
   */
  properties: {
    showModel: Boolean,
    showCancel: Boolean,
    title: String,
    contents: String,
  },

  options: {
    addGlobalClass:true,
  },
  /**
   * 组件的初始数据
   */
  data: {
    resolve: undefined as ((r: ModelResult) => void) | undefined,

  },

  /**
   * 组件的方法列表
   */
  methods: {
    onCancel(){
       this.hideModel('cancel')
       
    },
    onOK(){
      this.hideModel('ok')
    },
    hideModel(res: ModelResult){
      this.setData({
        showModel: false,
      })
      this.triggerEvent(res)
      if (this.data.resolve) {
        this.data.resolve(res)
      }
    },

    showModel(): Promise<ModelResult>{
      this.setData({
        showModel:true,
      })
      return new Promise((resolve) => {
        // save resolve function
        this.data.resolve = resolve
      })
    }
  }
})
